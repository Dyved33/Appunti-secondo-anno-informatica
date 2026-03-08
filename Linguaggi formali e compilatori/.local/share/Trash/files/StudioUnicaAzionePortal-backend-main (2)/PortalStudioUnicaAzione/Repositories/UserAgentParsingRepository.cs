using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Handlers;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using PortalStudioUnicaAzione.Services.Interfaces;
using UAParser;

namespace PortalStudioUnicaAzione.Repositories
{
    public class UserAgentParsingRepository(PortalStudioUnicaAzioneContext appDbContext, IDbLogger logger) : IUserAgentParsingRepository
    {
        public async Task<ParsedUserAgent> ParseUserAgent(string userAgent)
        {
            try
            {
                var parser = Parser.GetDefault();
                var clientInfo = parser.Parse(userAgent);
                var result = new ParsedUserAgent
                {
                    OS = clientInfo.OS.Family,
                    OSVersion = $"{clientInfo.OS.Major}.{clientInfo.OS.Minor}".TrimEnd('.'),
                    Browser = clientInfo.UA.Family,
                    BrowserVersion = $"{clientInfo.UA.Major}.{clientInfo.UA.Minor}".TrimEnd('.'),
                    DeviceModel = clientInfo.Device.Model,
                    DeviceBrand = clientInfo.Device.Brand,
                    // Determine device type
                    DeviceType = DetermineDeviceType(clientInfo, userAgent),
                    // Bot detection
                    IsBot = DetectBot(userAgent, clientInfo)
                };
                if (result.IsBot)
                {
                    result.BotType = DetectBotType(userAgent);
                }
                return result;
            }
            catch (Exception ex)
            {
                await logger.LogException(ex);
                return new ParsedUserAgent();
            }
        }

        public async Task<Guid> GetOrCreateDeviceInfoAsync(string userAgent)
        {
            // Check if we already have this user agent
            var existing = await appDbContext.ShortUrlDeviceInfo.FirstOrDefaultAsync(d => d.UserAgent == userAgent);
            if (existing != null)
            {
                return existing.Id;
            }
            // Parse and create new device info
            var parsed = await ParseUserAgent(userAgent);
            var deviceInfo = new ShortUrlDeviceInfo
            {
                Id = Uuid7.NewGuid(),
                UserAgent = userAgent,
                DeviceType = parsed.DeviceType,
                DeviceModel = parsed.DeviceModel,
                DeviceBrand = parsed.DeviceBrand,
                Os = parsed.OS,
                Osversion = parsed.OSVersion,
                Browser = parsed.Browser,
                BrowserVersion = parsed.BrowserVersion,
                IsBot = parsed.IsBot,
                BotType = parsed.BotType,
                CreatedAt = DateTime.UtcNow
            };
            await appDbContext.ShortUrlDeviceInfo.AddAsync(deviceInfo);
            await appDbContext.SaveChangesAsync();
            return deviceInfo.Id;
        }

        private string DetermineDeviceType(ClientInfo clientInfo, string userAgent)
        {
            // Check for TV/Smart TV
            if (userAgent.Contains("TV", StringComparison.OrdinalIgnoreCase) ||
                userAgent.Contains("SmartTV", StringComparison.OrdinalIgnoreCase))
                return "TV";
            // Check for tablet
            if (userAgent.Contains("iPad", StringComparison.OrdinalIgnoreCase) ||
                userAgent.Contains("tablet", StringComparison.OrdinalIgnoreCase) ||
                (clientInfo.Device.Family?.Contains("tablet", StringComparison.OrdinalIgnoreCase) == true))
                return "Tablet";
            // Check for mobile
            if (userAgent.Contains("Mobile", StringComparison.OrdinalIgnoreCase) ||
                userAgent.Contains("iPhone", StringComparison.OrdinalIgnoreCase) ||
                userAgent.Contains("Android", StringComparison.OrdinalIgnoreCase) ||
                clientInfo.Device.IsSpider ||
                !string.IsNullOrEmpty(clientInfo.Device.Model))
                return "Mobile";
            // Check for bot
            if (DetectBot(userAgent, clientInfo))
                return "Bot";
            return "Desktop";
        }

        private bool DetectBot(string userAgent, ClientInfo clientInfo)
        {
            if (clientInfo.Device.IsSpider) return true;
            var botIndicators = new[]
            {
                "bot", "crawler", "spider", "scraper", "facebook", "twitter",
                "googlebot", "bingbot", "slurp", "duckduckbot", "baiduspider"
            };
            return botIndicators.Any(indicator =>
                userAgent.Contains(indicator, StringComparison.OrdinalIgnoreCase));
        }

        private string? DetectBotType(string userAgent)
        {
            var botTypes = new Dictionary<string, string>
                {
                    { "googlebot", "GoogleBot" },
                    { "bingbot", "BingBot" },
                    { "facebookexternalhit", "FacebookBot" },
                    { "twitterbot", "TwitterBot" },
                    { "linkedinbot", "LinkedInBot" },
                    { "slurp", "YahooBot" },
                    { "duckduckbot", "DuckDuckGoBot" },
                    { "baiduspider", "BaiduBot" }
                };
            foreach (var botType in botTypes)
            {
                if (userAgent.Contains(botType.Key, StringComparison.OrdinalIgnoreCase))
                    return botType.Value;
            }
            return "Unknown";
        }
    }
}