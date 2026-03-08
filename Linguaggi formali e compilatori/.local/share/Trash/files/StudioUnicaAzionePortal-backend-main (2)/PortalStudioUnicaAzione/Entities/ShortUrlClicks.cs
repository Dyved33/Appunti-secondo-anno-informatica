using System.Text.Json.Serialization;

namespace PortalStudioUnicaAzione.Entities;

public partial class ShortUrlClicks
{
    public Guid Id { get; set; }

    public Guid ShortUrlId { get; set; }

    public DateTime ClickedAt { get; set; }

    public string IpAddress { get; set; } = null!;

    public string UserAgent { get; set; } = null!;

    public string? RefererUrl { get; set; }

    public string? AcceptLanguage { get; set; }

    public string SessionId { get; set; } = null!;

    public int ClickSequence { get; set; }

    public int ResponseTimeMs { get; set; }

    public bool IsProcessed { get; set; }

    public DateTime? ProcessedAt { get; set; }

    public DateTime? InsertTimestamp { get; set; }

    public Guid? InsertUser { get; set; }

    public DateTime? UpdateTimestamp { get; set; }

    public Guid? UpdateUser { get; set; }

    public bool Deleted { get; set; }

    [JsonIgnore]  public virtual ShortUrls ShortUrl { get; set; } = null!;

    [JsonIgnore]  public virtual ShortUrlClickAnalytics? ShortUrlClickAnalytics { get; set; }
}
