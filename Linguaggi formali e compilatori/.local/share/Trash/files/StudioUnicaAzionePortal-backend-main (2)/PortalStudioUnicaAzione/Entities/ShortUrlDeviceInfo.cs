using System.Text.Json.Serialization;

namespace PortalStudioUnicaAzione.Entities;

public partial class ShortUrlDeviceInfo
{
    public Guid Id { get; set; }

    public string UserAgent { get; set; } = null!;

    public string DeviceType { get; set; } = null!;

    public string? DeviceModel { get; set; }

    public string? DeviceBrand { get; set; }

    public string? Os { get; set; }

    public string? Osversion { get; set; }

    public string? Browser { get; set; }

    public string? BrowserVersion { get; set; }

    public bool IsBot { get; set; }

    public string? BotType { get; set; }

    public DateTime CreatedAt { get; set; }

    [JsonIgnore]  public virtual ICollection<ShortUrlClickAnalytics> ShortUrlClickAnalytics { get; set; } = new List<ShortUrlClickAnalytics>();
}
