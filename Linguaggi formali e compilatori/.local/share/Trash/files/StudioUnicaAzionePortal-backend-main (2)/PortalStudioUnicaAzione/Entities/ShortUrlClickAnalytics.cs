using System.Text.Json.Serialization;

namespace PortalStudioUnicaAzione.Entities;

public partial class ShortUrlClickAnalytics
{
    public Guid Id { get; set; }

    public Guid ClickId { get; set; }

    public Guid ShortUrlId { get; set; }

    public string SessionId { get; set; } = null!;

    public Guid? DeviceInfoId { get; set; }

    public Guid? GeoInfoId { get; set; }

    public Guid? ReferrerInfoId { get; set; }

    public string ClickHash { get; set; } = null!;

    public bool IsUniqueClick { get; set; }

    public bool IsUniqueDaily { get; set; }

    public bool IsUniqueSession { get; set; }

    public int ClickVelocity { get; set; }

    public decimal BotScore { get; set; }

    public decimal FraudScore { get; set; }

    public bool IsWeekend { get; set; }

    public bool IsBusinessHours { get; set; }

    public int LocalHour { get; set; }

    public int DayOfWeek { get; set; }

    public DateTime CreatedAt { get; set; }

    [JsonIgnore] public virtual ShortUrlClicks Click { get; set; } = null!;

    [JsonIgnore] public virtual ShortUrlDeviceInfo? DeviceInfo { get; set; }

    [JsonIgnore] public virtual ShortUrlGeoInfo? GeoInfo { get; set; }

    [JsonIgnore] public virtual ShortUrlReferrerInfo? ReferrerInfo { get; set; }

    [JsonIgnore] public virtual ShortUrls ShortUrl { get; set; } = null!;
}
