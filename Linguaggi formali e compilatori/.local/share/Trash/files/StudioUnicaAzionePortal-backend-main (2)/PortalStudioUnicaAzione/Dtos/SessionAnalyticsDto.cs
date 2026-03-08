namespace PortalStudioUnicaAzione.Dtos
{
    public class SessionAnalyticsDto
    {
        public int TotalSessions { get; set; }
        public int AverageSessionDuration { get; set; }
        public decimal BounceRate { get; set; }
        public int AverageClicksPerSession { get; set; }
        public List<SessionDetailDto> TopSessions { get; set; } = new();
    }

    public class SessionDetailDto
    {
        public string SessionId { get; set; } = null!;
        public DateTime FirstClickAt { get; set; }
        public DateTime LastClickAt { get; set; }
        public int TotalClicks { get; set; }
        public int SessionDuration { get; set; }
        public string? SessionSource { get; set; }
    }
}
