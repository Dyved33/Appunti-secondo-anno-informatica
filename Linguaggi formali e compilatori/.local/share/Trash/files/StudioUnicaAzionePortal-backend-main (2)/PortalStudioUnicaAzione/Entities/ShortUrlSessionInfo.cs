namespace PortalStudioUnicaAzione.Entities;

public partial class ShortUrlSessionInfo
{
    public Guid Id { get; set; }

    public string SessionId { get; set; } = null!;

    public DateTime FirstClickAt { get; set; }

    public DateTime LastClickAt { get; set; }

    public int TotalClicks { get; set; }

    public int UniqueUrlsClicked { get; set; }

    public int SessionDurationSeconds { get; set; }

    public bool IsNewSession { get; set; }

    public string? SessionSource { get; set; }

    public bool DeviceConsistency { get; set; }

    public bool LocationConsistency { get; set; }

    public DateTime CreatedAt { get; set; }

    public DateTime UpdatedAt { get; set; }
}
