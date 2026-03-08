using System.Text.Json.Serialization;

namespace PortalStudioUnicaAzione.Entities;

public partial class ShortUrlGeoInfo
{
    public Guid Id { get; set; }

    public string IpAddress { get; set; } = null!;

    public string? Country { get; set; }

    public string? CountryCode { get; set; }

    public string? Region { get; set; }

    public string? City { get; set; }

    public string? PostalCode { get; set; }

    public string? TimeZone { get; set; }

    public decimal? Latitude { get; set; }

    public decimal? Longitude { get; set; }

    public string? Isp { get; set; }

    public string? Organization { get; set; }

    public bool IsVpn { get; set; }

    public bool IsTor { get; set; }

    public bool IsProxy { get; set; }

    public bool IsDataCenter { get; set; }

    public DateTime CreatedAt { get; set; }

    public DateTime UpdatedAt { get; set; }

    [JsonIgnore]  public virtual ICollection<ShortUrlClickAnalytics> ShortUrlClickAnalytics { get; set; } = new List<ShortUrlClickAnalytics>();
}
