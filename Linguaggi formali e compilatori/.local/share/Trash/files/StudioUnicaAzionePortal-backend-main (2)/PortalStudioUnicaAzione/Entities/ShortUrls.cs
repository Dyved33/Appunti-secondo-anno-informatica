using System.Text.Json.Serialization;

namespace PortalStudioUnicaAzione.Entities;

public partial class ShortUrls : BaseEntity
{
    public Guid? CompanyId { get; set; }

    public string OriginalUrl { get; set; } = null!;

    public string ShortCode { get; set; } = null!;

    public DateTime? ValidityFromDate { get; set; }

    public DateTime? ValidityToDate { get; set; }

    public Boolean Enabled { get; set; }

    [JsonIgnore]  public virtual Users? InsertUser { get; set; }

    [JsonIgnore]  public virtual Companies? Company { get; set; }

    [JsonIgnore]  public virtual ICollection<ShortUrlClickAnalytics> ShortUrlClickAnalytics { get; set; } = new List<ShortUrlClickAnalytics>();

    [JsonIgnore]  public virtual ICollection<ShortUrlClickHeatMap> ShortUrlClickHeatMap { get; set; } = new List<ShortUrlClickHeatMap>();

    [JsonIgnore]  public virtual ICollection<ShortUrlClicks> ShortUrlClicks { get; set; } = new List<ShortUrlClicks>();

    [JsonIgnore]  public virtual ICollection<ShortUrlDailyClickStats> ShortUrlDailyClickStats { get; set; } = new List<ShortUrlDailyClickStats>();

    [JsonIgnore]  public virtual ICollection<ShortUrlHourlyClickStats> ShortUrlHourlyClickStats { get; set; } = new List<ShortUrlHourlyClickStats>();

    [JsonIgnore]  public virtual Users? UpdateUser { get; set; }
}
