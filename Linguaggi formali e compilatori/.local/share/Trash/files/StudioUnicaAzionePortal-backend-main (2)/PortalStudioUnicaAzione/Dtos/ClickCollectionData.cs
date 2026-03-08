namespace PortalStudioUnicaAzione.Dtos
{
    public class ClickCollectionData
    {
        public Guid ShortUrlId { get; set; }
        public string IpAddress { get; set; } = null!;
        public string UserAgent { get; set; } = null!;
        public string? RefererUrl { get; set; }
        public string? AcceptLanguage { get; set; }
        public string SessionId { get; set; } = null!;
        public int ResponseTimeMs { get; set; }
        public DateTime ClickedAt { get; set; }
    }
}
