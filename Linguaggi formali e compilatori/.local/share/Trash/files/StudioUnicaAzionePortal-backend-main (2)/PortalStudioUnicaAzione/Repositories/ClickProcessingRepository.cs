using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Handlers;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using PortalStudioUnicaAzione.Services.Interfaces;
using System.Security.Cryptography;
using System.Text;

namespace PortalStudioUnicaAzione.Repositories
{
    public class ClickProcessingRepository(
        PortalStudioUnicaAzioneContext appDbContext,
        IDbLogger logger,
        IUserAgentParsingService userAgentService,
        IGeoLocationService geoLocationService,
        IReferrerParsingRepository referrerService
        ) : IClickProcessingRepository
    {
        public async Task ProcessPendingClicksAsync()
        {
            try
            {
                var pendingClicks = await appDbContext.ShortUrlClicks
                    .Where(c => !c.IsProcessed)
                    .OrderBy(c => c.ClickedAt)
                    .Take(100) // Process in batches
                    .ToListAsync();
                foreach (var click in pendingClicks)
                {
                    await ProcessClickAsync(click.Id);
                }
            }
            catch (Exception ex)
            {
                await logger.LogException(ex);
                throw;
            }
        }

        public async Task ProcessClickAsync(Guid clickId)
        {
            try
            {
                var click = await appDbContext.ShortUrlClicks.FindAsync(clickId);
                if (click == null || click.IsProcessed)
                    return;
                // 1. Parse and normalize data
                var deviceInfoId = await userAgentService.GetOrCreateDeviceInfoAsync(click.UserAgent);
                var geoInfoId = await geoLocationService.GetOrCreateGeoInfoAsync(click.IpAddress);
                var referrerInfoId = await referrerService.GetOrCreateReferrerInfoAsync(click.RefererUrl, click.AcceptLanguage);
                // 2. Calculate analytics
                var clickHash = GenerateClickHash(click.IpAddress, click.UserAgent, click.ClickedAt.Date);
                var isUniqueClick = await IsUniqueClickAsync(clickHash, click.ShortUrlId);
                var isUniqueDaily = await IsUniqueDailyClickAsync(clickHash, click.ShortUrlId, click.ClickedAt.Date);
                var isUniqueSession = await IsUniqueSessionClickAsync(click.SessionId, click.ShortUrlId);
                // 3. Create analytics record
                var analytics = new ShortUrlClickAnalytics
                {
                    Id = Uuid7.NewGuid(),
                    ClickId = click.Id,
                    ShortUrlId = click.ShortUrlId,
                    SessionId = click.SessionId,
                    DeviceInfoId = deviceInfoId,
                    GeoInfoId = geoInfoId,
                    ReferrerInfoId = referrerInfoId,
                    ClickHash = clickHash,
                    IsUniqueClick = isUniqueClick,
                    IsUniqueDaily = isUniqueDaily,
                    IsUniqueSession = isUniqueSession,
                    ClickVelocity = await CalculateClickVelocityAsync(click.ShortUrlId, click.ClickedAt),
                    BotScore = await CalculateBotScoreAsync(deviceInfoId, geoInfoId),
                    FraudScore = await CalculateFraudScoreAsync(click),
                    IsWeekend = IsWeekend(click.ClickedAt),
                    IsBusinessHours = IsBusinessHours(click.ClickedAt),
                    LocalHour = await GetLocalHourAsync(geoInfoId, click.ClickedAt),
                    DayOfWeek = (int)click.ClickedAt.DayOfWeek + 1,
                    CreatedAt = DateTime.UtcNow
                };
                await appDbContext.ShortUrlClickAnalytics.AddAsync(analytics);
                // 4. Update session info
                await UpdateSessionInfoAsync(click.SessionId, click.ShortUrlId, click.ClickedAt);
                // 5. Mark click as processed
                click.IsProcessed = true;
                click.ProcessedAt = DateTime.UtcNow;
                await appDbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                await logger.LogException(ex);
                throw;
            }
        }

        private string GenerateClickHash(string ipAddress, string userAgent, DateTime date)
        {
            var input = $"{ipAddress}|{userAgent}|{date:yyyy-MM-dd}";
            using var sha256 = SHA256.Create();
            var hash = sha256.ComputeHash(Encoding.UTF8.GetBytes(input));
            return Convert.ToHexString(hash);
        }

        private async Task<bool> IsUniqueClickAsync(string clickHash, Guid shortUrlId)
        {
            return !await appDbContext.ShortUrlClickAnalytics
                .AnyAsync(a => a.ClickHash == clickHash && a.ShortUrlId == shortUrlId);
        }

        private async Task<bool> IsUniqueDailyClickAsync(string clickHash, Guid shortUrlId, DateTime date)
        {
            return !await appDbContext.ShortUrlClickAnalytics
                .AnyAsync(a => a.ClickHash == clickHash &&
                                a.ShortUrlId == shortUrlId &&
                                a.CreatedAt.Date == date);
        }

        private async Task<bool> IsUniqueSessionClickAsync(string sessionId, Guid shortUrlId)
        {
            return !await appDbContext.ShortUrlClickAnalytics
                .AnyAsync(a => a.SessionId == sessionId && a.ShortUrlId == shortUrlId);
        }

        private async Task<int> CalculateClickVelocityAsync(Guid shortUrlId, DateTime clickTime)
        {
            var oneMinuteAgo = clickTime.AddMinutes(-1);
            var clicksInLastMinute = await appDbContext.ShortUrlClicks
                .CountAsync(c => c.ShortUrlId == shortUrlId &&
                                c.ClickedAt >= oneMinuteAgo &&
                                c.ClickedAt <= clickTime);
            return clicksInLastMinute;
        }

        private async Task<decimal> CalculateBotScoreAsync(Guid deviceInfoId, Guid geoInfoId)
        {
            var deviceInfo = await appDbContext.ShortUrlDeviceInfo.FindAsync(deviceInfoId);
            var geoInfo = await appDbContext.ShortUrlGeoInfo.FindAsync(geoInfoId);
            decimal score = 0.0m;
            if (deviceInfo?.IsBot == true) score += 0.8m;
            if (geoInfo?.IsDataCenter == true) score += 0.3m;
            if (geoInfo?.IsVpn == true) score += 0.2m;
            if (geoInfo?.IsProxy == true) score += 0.2m;
            return Math.Min(score, 1.0m);
        }

        private async Task<decimal> CalculateFraudScoreAsync(ShortUrlClicks click)
        {
            decimal score = 0.0m;
            // Check for rapid clicking from same IP
            var recentClicksFromIP = await appDbContext.ShortUrlClicks
                .CountAsync(c => c.IpAddress == click.IpAddress &&
                                c.ClickedAt >= click.ClickedAt.AddMinutes(-5) &&
                                c.ClickedAt <= click.ClickedAt);
            if (recentClicksFromIP > 10) score += 0.5m;
            if (recentClicksFromIP > 20) score += 0.3m;
            // Check for suspicious user agent
            if (string.IsNullOrEmpty(click.UserAgent) || click.UserAgent.Length < 10)
                score += 0.4m;
            return Math.Min(score, 1.0m);
        }

        private bool IsWeekend(DateTime dateTime)
        {
            return dateTime.DayOfWeek == DayOfWeek.Saturday || dateTime.DayOfWeek == DayOfWeek.Sunday;
        }

        private bool IsBusinessHours(DateTime dateTime)
        {
            var hour = dateTime.Hour;
            return hour >= 9 && hour <= 17 && !IsWeekend(dateTime);
        }

        private async Task<int> GetLocalHourAsync(Guid geoInfoId, DateTime utcTime)
        {
            var geoInfo = await appDbContext.ShortUrlGeoInfo.FindAsync(geoInfoId);
            if (geoInfo?.TimeZone != null)
            {
                try
                {
                    var timeZone = TimeZoneInfo.FindSystemTimeZoneById(geoInfo.TimeZone);
                    var localTime = TimeZoneInfo.ConvertTimeFromUtc(utcTime, timeZone);
                    return localTime.Hour;
                }
                catch
                {
                    // Fall back to UTC if timezone conversion fails
                }
            }
            return utcTime.Hour;
        }

        private async Task UpdateSessionInfoAsync(string sessionId, Guid shortUrlId, DateTime clickTime)
        {
            var sessionInfo = await appDbContext.ShortUrlSessionInfo.FirstOrDefaultAsync(s => s.SessionId == sessionId);
            if (sessionInfo == null)
            {
                sessionInfo = new ShortUrlSessionInfo
                {
                    Id = Uuid7.NewGuid(),
                    SessionId = sessionId,
                    FirstClickAt = clickTime,
                    LastClickAt = clickTime,
                    TotalClicks = 1,
                    UniqueUrlsClicked = 1,
                    SessionDurationSeconds = 0,
                    IsNewSession = true,
                    DeviceConsistency = true,
                    LocationConsistency = true,
                    CreatedAt = DateTime.UtcNow,
                    UpdatedAt = DateTime.UtcNow
                };
                await appDbContext.ShortUrlSessionInfo.AddAsync(sessionInfo);
            }
            else
            {
                sessionInfo.LastClickAt = clickTime;
                sessionInfo.TotalClicks++;
                sessionInfo.SessionDurationSeconds = (int)(clickTime - sessionInfo.FirstClickAt).TotalSeconds;
                sessionInfo.UpdatedAt = DateTime.UtcNow;
                // Check if this is a new unique URL for this session
                var uniqueUrlsInSession = await appDbContext.ShortUrlClicks
                    .Where(c => c.SessionId == sessionId)
                    .Select(c => c.ShortUrlId)
                    .Distinct()
                    .CountAsync();
                sessionInfo.UniqueUrlsClicked = uniqueUrlsInSession;
            }
        }
    }
}
