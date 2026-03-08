using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Repositories.Interfaces;

namespace PortalStudioUnicaAzione.Repositories
{
    public class AnalyticsRepository(PortalStudioUnicaAzioneContext appDbContext) : IAnalyticsRepository
    {
        public async Task<AnalyticsResponseDto> GetAnalyticsByShortCodeAsync(string shortCode, bool includeDetails = false, DateTime? fromDate = null, DateTime? toDate = null)
        {
            var shortUrl = await appDbContext.ShortUrls
                .FirstOrDefaultAsync(u => u.ShortCode == shortCode && !u.Deleted);

            if (shortUrl == null)
                throw new ArgumentException("Short URL not found", nameof(shortCode));

            if (includeDetails)
                return await GetAnalyticsDetailedAsync(shortUrl.Id, fromDate, toDate);
            else
                return await GetAnalyticsSummaryAsync(shortUrl.Id, fromDate, toDate);
        }

        public async Task<AnalyticsResponseDto> GetAnalyticsSummaryAsync(Guid shortUrlId, DateTime? fromDate = null, DateTime? toDate = null)
        {
            var shortUrl = await appDbContext.ShortUrls.FindAsync(shortUrlId);
            if (shortUrl == null)
                throw new ArgumentException("Short URL not found", nameof(shortUrlId));

            var summary = await BuildSummaryAsync(shortUrlId, fromDate, toDate);

            return new AnalyticsResponseDto
            {
                ShortCode = shortUrl.ShortCode,
                OriginalUrl = shortUrl.OriginalUrl,
                CreatedAt = shortUrl.InsertTimestamp ?? DateTime.UtcNow,
                FromDate = fromDate,
                ToDate = toDate,
                Summary = summary
            };
        }

        public async Task<AnalyticsResponseDto> GetAnalyticsDetailedAsync(Guid shortUrlId, DateTime? fromDate = null, DateTime? toDate = null)
        {
            var shortUrl = await appDbContext.ShortUrls.FindAsync(shortUrlId);
            if (shortUrl == null)
                throw new ArgumentException("Short URL not found", nameof(shortUrlId));

            // Get all data in parallel for better performance
            var summaryTask = BuildSummaryAsync(shortUrlId, fromDate, toDate);
            var temporalTask = BuildTemporalAnalyticsAsync(shortUrlId, fromDate, toDate);
            var geographicTask = BuildGeographicAnalyticsAsync(shortUrlId, fromDate, toDate);
            var deviceTask = BuildDeviceAnalyticsAsync(shortUrlId, fromDate, toDate);
            var referrerTask = BuildReferrerAnalyticsAsync(shortUrlId, fromDate, toDate);
            var sessionTask = BuildSessionAnalyticsAsync(shortUrlId, fromDate, toDate);

            await Task.WhenAll(summaryTask, temporalTask, geographicTask, deviceTask, referrerTask, sessionTask);

            var details = new AnalyticsDetailedDto
            {
                Summary = await summaryTask,
                Temporal = await temporalTask,
                Geographic = await geographicTask,
                Devices = await deviceTask,
                Referrers = await referrerTask,
                Sessions = await sessionTask
            };

            return new AnalyticsResponseDto
            {
                ShortCode = shortUrl.ShortCode,
                OriginalUrl = shortUrl.OriginalUrl,
                CreatedAt = shortUrl.InsertTimestamp ?? DateTime.UtcNow,
                FromDate = fromDate,
                ToDate = toDate,
                Summary = details.Summary,
                Details = details
            };
        }

        private async Task<AnalyticsSummaryDto> BuildSummaryAsync(Guid shortUrlId, DateTime? fromDate, DateTime? toDate)
        {
            // Single query to get all summary data
            var query = from analytics in appDbContext.ShortUrlClickAnalytics
                        join clicks in appDbContext.ShortUrlClicks on analytics.ClickId equals clicks.Id
                        where analytics.ShortUrlId == shortUrlId
                        select new { analytics, clicks };

            if (fromDate.HasValue)
                query = query.Where(x => x.analytics.CreatedAt >= fromDate.Value);
            if (toDate.HasValue)
                query = query.Where(x => x.analytics.CreatedAt <= toDate.Value);

            var data = await query.ToListAsync();

            if (!data.Any())
            {
                return new AnalyticsSummaryDto();
            }

            var totalClicks = data.Count;
            var uniqueClicks = data.Count(x => x.analytics.IsUniqueClick);
            var botClicks = data.Count(x => x.analytics.BotScore > 0.5m);
            var avgResponseTime = data.Any() ? (int)data.Average(x => x.clicks.ResponseTimeMs) : 0;
            var lastClick = data.Max(x => x.clicks.ClickedAt);

            // Get session data
            var sessionIds = data.Select(x => x.analytics.SessionId).Distinct().ToList();
            var sessionData = await appDbContext.ShortUrlSessionInfo
                .Where(s => sessionIds.Contains(s.SessionId))
                .ToListAsync();

            return new AnalyticsSummaryDto
            {
                TotalClicks = totalClicks,
                UniqueClicks = uniqueClicks,
                AverageResponseTime = avgResponseTime,
                LastClickAt = lastClick,
                BotClicks = botClicks,
                BotPercentage = totalClicks > 0 ? Math.Round((decimal)botClicks / totalClicks * 100, 2) : 0,
                TotalSessions = sessionData.Count,
                AverageSessionDuration = sessionData.Any() ? (int)sessionData.Average(s => s.SessionDurationSeconds) : 0
            };
        }

        private async Task<TemporalAnalyticsDto> BuildTemporalAnalyticsAsync(Guid shortUrlId, DateTime? fromDate, DateTime? toDate)
        {
            var query = appDbContext.ShortUrlClickAnalytics
                .Where(a => a.ShortUrlId == shortUrlId);

            if (fromDate.HasValue)
                query = query.Where(a => a.CreatedAt >= fromDate.Value);
            if (toDate.HasValue)
                query = query.Where(a => a.CreatedAt <= toDate.Value);

            var data = await query
                .Select(a => new
                {
                    Date = a.CreatedAt.Date,
                    Hour = a.LocalHour,
                    DayOfWeek = a.DayOfWeek,
                    IsUnique = a.IsUniqueClick,
                    IsUniqueDaily = a.IsUniqueDaily,
                    BotScore = a.BotScore
                })
                .ToListAsync();

            var totalClicks = data.Count;

            // Daily breakdown
            var daily = data
                .GroupBy(d => d.Date)
                .Select(g => new DailyClickDto
                {
                    Date = g.Key,
                    TotalClicks = g.Count(),
                    UniqueClicks = g.Count(x => x.IsUniqueDaily),
                    BotClicks = g.Count(x => x.BotScore > 0.5m),
                    AverageResponseTime = 0 // We'll calculate this separately if needed
                })
                .OrderBy(d => d.Date)
                .ToList();

            // Hourly breakdown (0-23)
            var hourly = data
                .GroupBy(d => d.Hour)
                .Select(g => new HourlyClickDto
                {
                    Hour = g.Key,
                    TotalClicks = g.Count(),
                    UniqueClicks = g.Count(x => x.IsUnique),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderBy(h => h.Hour)
                .ToList();

            // Day of week breakdown
            var dayNames = new[] { "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday" };
            var dayOfWeek = data
                .GroupBy(d => d.DayOfWeek)
                .Select(g => new DayOfWeekDto
                {
                    DayNumber = g.Key,
                    DayName = g.Key >= 1 && g.Key <= 7 ? dayNames[g.Key - 1] : "Unknown",
                    TotalClicks = g.Count(),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderBy(d => d.DayNumber)
                .ToList();

            return new TemporalAnalyticsDto
            {
                Daily = daily,
                Hourly = hourly,
                DayOfWeek = dayOfWeek
            };
        }

        private async Task<GeographicAnalyticsDto> BuildGeographicAnalyticsAsync(Guid shortUrlId, DateTime? fromDate, DateTime? toDate)
        {
            var query = from analytics in appDbContext.ShortUrlClickAnalytics
                        join geo in appDbContext.ShortUrlGeoInfo on analytics.GeoInfoId equals geo.Id
                        where analytics.ShortUrlId == shortUrlId
                        select new { analytics, geo };

            if (fromDate.HasValue)
                query = query.Where(x => x.analytics.CreatedAt >= fromDate.Value);
            if (toDate.HasValue)
                query = query.Where(x => x.analytics.CreatedAt <= toDate.Value);

            var data = await query.ToListAsync();
            var totalClicks = data.Count;

            // Countries
            var countries = data
                .Where(x => !string.IsNullOrEmpty(x.geo.Country))
                .GroupBy(x => new { x.geo.Country, x.geo.CountryCode })
                .Select(g => new CountryClickDto
                {
                    Country = g.Key.Country ?? "Unknown",
                    CountryCode = g.Key.CountryCode ?? "XX",
                    TotalClicks = g.Count(),
                    UniqueClicks = g.Count(x => x.analytics.IsUniqueClick),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(c => c.TotalClicks)
                .Take(20)
                .ToList();

            // Cities
            var cities = data
                .Where(x => !string.IsNullOrEmpty(x.geo.City))
                .GroupBy(x => new { x.geo.City, x.geo.Region, x.geo.Country, x.geo.Latitude, x.geo.Longitude })
                .Select(g => new CityClickDto
                {
                    City = g.Key.City ?? "Unknown",
                    Region = g.Key.Region ?? "",
                    Country = g.Key.Country ?? "",
                    TotalClicks = g.Count(),
                    Latitude = g.Key.Latitude,
                    Longitude = g.Key.Longitude
                })
                .OrderByDescending(c => c.TotalClicks)
                .Take(50)
                .ToList();

            // Heat map
            var heatMap = data
                .Where(x => x.geo.Latitude.HasValue && x.geo.Longitude.HasValue)
                .GroupBy(x => new { x.geo.Latitude, x.geo.Longitude, x.geo.City, x.geo.Country })
                .Select(g => new HeatMapPointDto
                {
                    Latitude = g.Key.Latitude!.Value,
                    Longitude = g.Key.Longitude!.Value,
                    ClickCount = g.Count(),
                    Location = $"{g.Key.City}, {g.Key.Country}"
                })
                .OrderByDescending(h => h.ClickCount)
                .Take(100)
                .ToList();

            return new GeographicAnalyticsDto
            {
                Countries = countries,
                Cities = cities,
                HeatMap = heatMap
            };
        }

        private async Task<DeviceAnalyticsDto> BuildDeviceAnalyticsAsync(Guid shortUrlId, DateTime? fromDate, DateTime? toDate)
        {
            var query = from analytics in appDbContext.ShortUrlClickAnalytics
                        join device in appDbContext.ShortUrlDeviceInfo on analytics.DeviceInfoId equals device.Id
                        where analytics.ShortUrlId == shortUrlId
                        select new { analytics, device };

            if (fromDate.HasValue)
                query = query.Where(x => x.analytics.CreatedAt >= fromDate.Value);
            if (toDate.HasValue)
                query = query.Where(x => x.analytics.CreatedAt <= toDate.Value);

            var data = await query.ToListAsync();
            var totalClicks = data.Count;

            // Device types
            var types = data
                .GroupBy(x => x.device.DeviceType)
                .Select(g => new DeviceTypeDto
                {
                    DeviceType = g.Key,
                    TotalClicks = g.Count(),
                    UniqueClicks = g.Count(x => x.analytics.IsUniqueClick),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(d => d.TotalClicks)
                .ToList();

            // Browsers
            var browsers = data
                .Where(x => !string.IsNullOrEmpty(x.device.Browser))
                .GroupBy(x => new { x.device.Browser, x.device.BrowserVersion })
                .Select(g => new BrowserDto
                {
                    Browser = g.Key.Browser!,
                    Version = g.Key.BrowserVersion,
                    TotalClicks = g.Count(),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(b => b.TotalClicks)
                .Take(15)
                .ToList();

            // Operating systems
            var operatingSystems = data
                .Where(x => !string.IsNullOrEmpty(x.device.Os))
                .GroupBy(x => new { x.device.Os, x.device.Osversion })
                .Select(g => new OperatingSystemDto
                {
                    OS = g.Key.Os!,
                    Version = g.Key.Osversion,
                    TotalClicks = g.Count(),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(os => os.TotalClicks)
                .Take(10)
                .ToList();

            // Device models (only top ones to avoid clutter)
            var models = data
                .Where(x => !string.IsNullOrEmpty(x.device.DeviceModel))
                .GroupBy(x => new { x.device.DeviceBrand, x.device.DeviceModel, x.device.DeviceType })
                .Select(g => new DeviceModelDto
                {
                    Brand = g.Key.DeviceBrand,
                    Model = g.Key.DeviceModel!,
                    DeviceType = g.Key.DeviceType,
                    TotalClicks = g.Count(),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(m => m.TotalClicks)
                .Take(15)
                .ToList();

            return new DeviceAnalyticsDto
            {
                Types = types,
                Browsers = browsers,
                OperatingSystems = operatingSystems,
                Models = models
            };
        }

        private async Task<ReferrerAnalyticsDto> BuildReferrerAnalyticsAsync(Guid shortUrlId, DateTime? fromDate, DateTime? toDate)
        {
            var query = from analytics in appDbContext.ShortUrlClickAnalytics
                        join referrer in appDbContext.ShortUrlReferrerInfo on analytics.ReferrerInfoId equals referrer.Id
                        where analytics.ShortUrlId == shortUrlId
                        select new { analytics, referrer };

            if (fromDate.HasValue)
                query = query.Where(x => x.analytics.CreatedAt >= fromDate.Value);
            if (toDate.HasValue)
                query = query.Where(x => x.analytics.CreatedAt <= toDate.Value);

            var data = await query.ToListAsync();
            var totalClicks = data.Count;

            // Referrer types
            var types = data
                .GroupBy(x => x.referrer.ReferrerType)
                .Select(g => new ReferrerTypeDto
                {
                    ReferrerType = g.Key,
                    TotalClicks = g.Count(),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(r => r.TotalClicks)
                .ToList();

            // Domains
            var domains = data
                .Where(x => !string.IsNullOrEmpty(x.referrer.ReferrerDomain))
                .GroupBy(x => new { x.referrer.ReferrerDomain, x.referrer.ReferrerType })
                .Select(g => new ReferrerDomainDto
                {
                    Domain = g.Key.ReferrerDomain!,
                    ReferrerType = g.Key.ReferrerType,
                    TotalClicks = g.Count(),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(d => d.TotalClicks)
                .Take(20)
                .ToList();

            // Social platforms
            var socialPlatforms = data
                .Where(x => !string.IsNullOrEmpty(x.referrer.SocialPlatform))
                .GroupBy(x => x.referrer.SocialPlatform)
                .Select(g => new SocialPlatformDto
                {
                    Platform = g.Key!,
                    TotalClicks = g.Count(),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(s => s.TotalClicks)
                .ToList();

            // UTM Campaigns
            var campaigns = data
                .Where(x => !string.IsNullOrEmpty(x.referrer.Utmsource) ||
                           !string.IsNullOrEmpty(x.referrer.Utmmedium) ||
                           !string.IsNullOrEmpty(x.referrer.Utmcampaign))
                .GroupBy(x => new { x.referrer.Utmsource, x.referrer.Utmmedium, x.referrer.Utmcampaign })
                .Select(g => new CampaignDto
                {
                    UTMSource = g.Key.Utmsource,
                    UTMMedium = g.Key.Utmmedium,
                    UTMCampaign = g.Key.Utmcampaign,
                    TotalClicks = g.Count(),
                    UniqueClicks = g.Count(x => x.analytics.IsUniqueClick)
                })
                .OrderByDescending(c => c.TotalClicks)
                .Take(15)
                .ToList();

            // Search engines
            var searchEngines = data
                .Where(x => !string.IsNullOrEmpty(x.referrer.SearchEngine))
                .GroupBy(x => new { x.referrer.SearchEngine, x.referrer.SearchKeywords })
                .Select(g => new SearchEngineDto
                {
                    SearchEngine = g.Key.SearchEngine!,
                    Keywords = g.Key.SearchKeywords,
                    TotalClicks = g.Count(),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderByDescending(s => s.TotalClicks)
                .Take(10)
                .ToList();

            return new ReferrerAnalyticsDto
            {
                Types = types,
                Domains = domains,
                SocialPlatforms = socialPlatforms,
                Campaigns = campaigns,
                SearchEngines = searchEngines
            };
        }

        private async Task<SessionAnalyticsDto> BuildSessionAnalyticsAsync(Guid shortUrlId, DateTime? fromDate, DateTime? toDate)
        {
            // Get session IDs for this short URL within date range
            var sessionQuery = appDbContext.ShortUrlClickAnalytics
                .Where(a => a.ShortUrlId == shortUrlId);

            if (fromDate.HasValue)
                sessionQuery = sessionQuery.Where(a => a.CreatedAt >= fromDate.Value);
            if (toDate.HasValue)
                sessionQuery = sessionQuery.Where(a => a.CreatedAt <= toDate.Value);

            var sessionIds = await sessionQuery
                .Select(a => a.SessionId)
                .Distinct()
                .ToListAsync();

            if (!sessionIds.Any())
            {
                return new SessionAnalyticsDto();
            }

            // Get session details
            var sessions = await appDbContext.ShortUrlSessionInfo
                .Where(s => sessionIds.Contains(s.SessionId))
                .ToListAsync();

            var totalSessions = sessions.Count;
            var avgSessionDuration = sessions.Any() ? (int)sessions.Average(s => s.SessionDurationSeconds) : 0;
            var bounceRate = totalSessions > 0 ? Math.Round((decimal)sessions.Count(s => s.TotalClicks == 1) / totalSessions, 2) : 0;
            var avgClicksPerSession = sessions.Any() ? (int)sessions.Average(s => s.TotalClicks) : 0;

            // Top sessions
            var topSessions = sessions
                .OrderByDescending(s => s.TotalClicks)
                .ThenByDescending(s => s.SessionDurationSeconds)
                .Take(10)
                .Select(s => new SessionDetailDto
                {
                    SessionId = s.SessionId,
                    FirstClickAt = s.FirstClickAt,
                    LastClickAt = s.LastClickAt,
                    TotalClicks = s.TotalClicks,
                    SessionDuration = s.SessionDurationSeconds,
                    SessionSource = s.SessionSource
                })
                .ToList();

            return new SessionAnalyticsDto
            {
                TotalSessions = totalSessions,
                AverageSessionDuration = avgSessionDuration,
                BounceRate = bounceRate,
                AverageClicksPerSession = avgClicksPerSession,
                TopSessions = topSessions
            };
        }

        public async Task<List<DailyClickDto>> GetDailyStatsAsync(Guid shortUrlId, DateTime? fromDate = null, DateTime? toDate = null)
        {
            // Try pre-aggregated data first (if available)
            var dailyStatsQuery = appDbContext.ShortUrlDailyClickStats
                .Where(d => d.ShortUrlId == shortUrlId);

            if (fromDate.HasValue)
                dailyStatsQuery = dailyStatsQuery.Where(d => d.Date >= DateOnly.FromDateTime(fromDate.Value));
            if (toDate.HasValue)
                dailyStatsQuery = dailyStatsQuery.Where(d => d.Date <= DateOnly.FromDateTime(toDate.Value));

            var preAggregated = await dailyStatsQuery.AnyAsync();

            if (preAggregated)
            {
                return await dailyStatsQuery
                    .Select(d => new DailyClickDto
                    {
                        Date = d.Date.ToDateTime(TimeOnly.MinValue),
                        TotalClicks = d.TotalClicks,
                        UniqueClicks = d.UniqueClicks,
                        BotClicks = d.BotClicks,
                        AverageResponseTime = d.AvgResponseTimeMs
                    })
                    .OrderBy(d => d.Date)
                    .ToListAsync();
            }

            // Fall back to raw analytics data
            var analyticsQuery = appDbContext.ShortUrlClickAnalytics
                .Where(a => a.ShortUrlId == shortUrlId);

            if (fromDate.HasValue)
                analyticsQuery = analyticsQuery.Where(a => a.CreatedAt >= fromDate.Value);
            if (toDate.HasValue)
                analyticsQuery = analyticsQuery.Where(a => a.CreatedAt <= toDate.Value);

            return await analyticsQuery
                .GroupBy(a => a.CreatedAt.Date)
                .Select(g => new DailyClickDto
                {
                    Date = g.Key,
                    TotalClicks = g.Count(),
                    UniqueClicks = g.Count(a => a.IsUniqueDaily),
                    BotClicks = g.Count(a => a.BotScore > 0.5m),
                    AverageResponseTime = 0 // Could join with clicks table if needed
                })
                .OrderBy(d => d.Date)
                .ToListAsync();
        }

        public async Task<List<HourlyClickDto>> GetHourlyStatsAsync(Guid shortUrlId, DateTime? fromDate = null, DateTime? toDate = null)
        {
            var analyticsQuery = appDbContext.ShortUrlClickAnalytics
                .Where(a => a.ShortUrlId == shortUrlId);

            if (fromDate.HasValue)
                analyticsQuery = analyticsQuery.Where(a => a.CreatedAt >= fromDate.Value);
            if (toDate.HasValue)
                analyticsQuery = analyticsQuery.Where(a => a.CreatedAt <= toDate.Value);

            var totalClicks = await analyticsQuery.CountAsync();

            return await analyticsQuery
                .GroupBy(a => a.LocalHour)
                .Select(g => new HourlyClickDto
                {
                    Hour = g.Key,
                    TotalClicks = g.Count(),
                    UniqueClicks = g.Count(a => a.IsUniqueClick),
                    Percentage = totalClicks > 0 ? Math.Round((decimal)g.Count() / totalClicks * 100, 2) : 0
                })
                .OrderBy(h => h.Hour)
                .ToListAsync();
        }

        public async Task<bool> HasAnalyticsDataAsync(Guid shortUrlId)
        {
            return await appDbContext.ShortUrlClickAnalytics
                .AnyAsync(a => a.ShortUrlId == shortUrlId);
        }
    }
}
