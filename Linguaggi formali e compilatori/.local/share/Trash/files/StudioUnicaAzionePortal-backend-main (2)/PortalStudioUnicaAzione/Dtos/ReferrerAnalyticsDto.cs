namespace PortalStudioUnicaAzione.Dtos
{
    public class ReferrerAnalyticsDto
    {
        public List<ReferrerTypeDto> Types { get; set; } = new();
        public List<ReferrerDomainDto> Domains { get; set; } = new();
        public List<SocialPlatformDto> SocialPlatforms { get; set; } = new();
        public List<CampaignDto> Campaigns { get; set; } = new();
        public List<SearchEngineDto> SearchEngines { get; set; } = new();
    }

    public class ReferrerTypeDto
    {
        public string ReferrerType { get; set; } = null!;
        public int TotalClicks { get; set; }
        public decimal Percentage { get; set; }
    }

    public class ReferrerDomainDto
    {
        public string Domain { get; set; } = null!;
        public string ReferrerType { get; set; } = null!;
        public int TotalClicks { get; set; }
        public decimal Percentage { get; set; }
    }

    public class SocialPlatformDto
    {
        public string Platform { get; set; } = null!;
        public int TotalClicks { get; set; }
        public decimal Percentage { get; set; }
    }

    public class CampaignDto
    {
        public string? UTMSource { get; set; }
        public string? UTMMedium { get; set; }
        public string? UTMCampaign { get; set; }
        public int TotalClicks { get; set; }
        public int UniqueClicks { get; set; }
    }

    public class SearchEngineDto
    {
        public string SearchEngine { get; set; } = null!;
        public string? Keywords { get; set; }
        public int TotalClicks { get; set; }
        public decimal Percentage { get; set; }
    }
}
