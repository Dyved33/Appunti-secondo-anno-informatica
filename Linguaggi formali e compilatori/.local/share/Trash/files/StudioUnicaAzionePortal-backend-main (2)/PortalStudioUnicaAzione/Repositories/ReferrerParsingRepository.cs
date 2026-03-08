using System.Web;
using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Handlers;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using PortalStudioUnicaAzione.Services.Interfaces;

namespace PortalStudioUnicaAzione.Repositories
{
    public class ReferrerParsingRepository(PortalStudioUnicaAzioneContext appDbContext, IDbLogger logger) : IReferrerParsingRepository
    {
        public async Task<ParsedReferrer> ParseReferrer(string? referrerUrl)
        {
            var result = new ParsedReferrer();
            if (string.IsNullOrEmpty(referrerUrl))
            {
                result.ReferrerType = "Direct";
                return result;
            }
            try
            {
                var uri = new Uri(referrerUrl);
                result.ReferrerDomain = uri.Host.ToLower();
                // Parse UTM parameters
                var query = HttpUtility.ParseQueryString(uri.Query);
                result.UTMSource = query["utm_source"];
                result.UTMMedium = query["utm_medium"];
                result.UTMCampaign = query["utm_campaign"];
                result.UTMTerm = query["utm_term"];
                result.UTMContent = query["utm_content"];
                // Determine referrer type and platform
                result.ReferrerType = DetermineReferrerType(result.ReferrerDomain);
                result.SocialPlatform = GetSocialPlatform(result.ReferrerDomain);
                result.SearchEngine = GetSearchEngine(result.ReferrerDomain);
                // Extract search keywords if from search engine
                if (result.SearchEngine != null)
                {
                    result.SearchKeywords = ExtractSearchKeywords(uri, result.SearchEngine);
                }

                return result;
            }
            catch (Exception ex)
            {
                await logger.LogException(ex);
                result.ReferrerType = "Unknown";
                return result;
            }
        }

        public async Task<Guid> GetOrCreateReferrerInfoAsync(string? referrerUrl, string? acceptLanguage)
        {
            // Create a hash key for lookup
            var hashKey = $"{referrerUrl}|{acceptLanguage}".GetHashCode().ToString();
            var existing = await appDbContext.ShortUrlReferrerInfo.FirstOrDefaultAsync(r => r.ReferrerUrl == referrerUrl && r.AcceptLanguage == acceptLanguage);
            if (existing != null)
            {
                return existing.Id;
            }
            var parsed = await ParseReferrer(referrerUrl);
            var referrerInfo = new ShortUrlReferrerInfo
            {
                Id = Uuid7.NewGuid(),
                ReferrerUrl = referrerUrl,
                ReferrerDomain = parsed.ReferrerDomain,
                ReferrerType = parsed.ReferrerType,
                SocialPlatform = parsed.SocialPlatform,
                SearchEngine = parsed.SearchEngine,
                SearchKeywords = parsed.SearchKeywords,
                Utmsource = parsed.UTMSource,
                Utmmedium = parsed.UTMMedium,
                Utmcampaign = parsed.UTMCampaign,
                Utmterm = parsed.UTMTerm,
                Utmcontent = parsed.UTMContent,
                AcceptLanguage = acceptLanguage,
                CreatedAt = DateTime.UtcNow
            };
            await appDbContext.ShortUrlReferrerInfo.AddAsync(referrerInfo);
            await appDbContext.SaveChangesAsync();
            return referrerInfo.Id;
        }

        private string DetermineReferrerType(string domain)
        {
            var socialPlatforms = new[] { "facebook.com", "twitter.com", "instagram.com", "linkedin.com", "tiktok.com", "youtube.com" };
            var searchEngines = new[] { "google.com", "bing.com", "yahoo.com", "duckduckgo.com", "baidu.com" };
            if (socialPlatforms.Any(p => domain.Contains(p))) return "Social";
            if (searchEngines.Any(s => domain.Contains(s))) return "Search";
            if (domain.Contains("mail") || domain.Contains("gmail") || domain.Contains("outlook")) return "Email";
            return "Website";
        }

        private string? GetSocialPlatform(string domain)
        {
            var platforms = new Dictionary<string, string>
                {
                    { "facebook.com", "Facebook" },
                    { "twitter.com", "Twitter" },
                    { "instagram.com", "Instagram" },
                    { "linkedin.com", "LinkedIn" },
                    { "tiktok.com", "TikTok" },
                    { "youtube.com", "YouTube" },
                    { "pinterest.com", "Pinterest" },
                    { "snapchat.com", "Snapchat" },
                    { "reddit.com", "Reddit" }
                };
            return platforms.FirstOrDefault(p => domain.Contains(p.Key)).Value;
        }

        private string? GetSearchEngine(string domain)
        {
            var engines = new Dictionary<string, string>
                {
                    { "google.com", "Google" },
                    { "bing.com", "Bing" },
                    { "yahoo.com", "Yahoo" },
                    { "duckduckgo.com", "DuckDuckGo" },
                    { "baidu.com", "Baidu" },
                    { "yandex.com", "Yandex" }
                };
            return engines.FirstOrDefault(e => domain.Contains(e.Key)).Value;
        }

        private string? ExtractSearchKeywords(Uri uri, string searchEngine)
        {
            var query = HttpUtility.ParseQueryString(uri.Query);
            return searchEngine.ToLower() switch
            {
                "google" => query["q"],
                "bing" => query["q"],
                "yahoo" => query["p"],
                "duckduckgo" => query["q"],
                "baidu" => query["wd"],
                _ => null
            };
        }
    }
}