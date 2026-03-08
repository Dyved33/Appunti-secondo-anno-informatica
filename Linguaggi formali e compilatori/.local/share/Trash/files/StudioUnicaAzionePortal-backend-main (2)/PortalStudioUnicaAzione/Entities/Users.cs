using System.Text.Json.Serialization;

namespace PortalStudioUnicaAzione.Entities;

public partial class Users : BaseEntity
{
    public string Username { get; set; } = null!;

    public string Password { get; set; } = null!;

    public string? Firstname { get; set; }

    public string? Lastname { get; set; }

    public string? Picture { get; set; }
    
    public string Email { get; set; } = null!;

    public Guid? CompanyId { get; set; }

    public Guid RoleId { get; set; }

    public bool Enabled { get; set; }

    public string? PasswordResetToken { get; set; }

    public DateTime? PasswordResetExpires { get; set; }

    [JsonIgnore]  public virtual ICollection<Companies> CompaniesInsertUser { get; set; } = new List<Companies>();

    [JsonIgnore]  public virtual ICollection<Companies> CompaniesUpdateUser { get; set; } = new List<Companies>();

    [JsonIgnore]  public virtual Companies? Company { get; set; }

    [JsonIgnore]  public virtual ICollection<Events> EventsInsertUser { get; set; } = new List<Events>();

    [JsonIgnore]  public virtual ICollection<Events> EventsUpdateUser { get; set; } = new List<Events>();

    [JsonIgnore]  public virtual Users? InsertUser { get; set; }

    [JsonIgnore]  public virtual ICollection<Users> InverseInsertUser { get; set; } = new List<Users>();

    [JsonIgnore]  public virtual ICollection<Users> InverseUpdateUser { get; set; } = new List<Users>();

    [JsonIgnore]  public virtual ICollection<Logs> Logs { get; set; } = new List<Logs>();

    [JsonIgnore]  public virtual ICollection<MailTemplates> MailTemplatesInsertUser { get; set; } = new List<MailTemplates>();

    [JsonIgnore]  public virtual ICollection<MailTemplates> MailTemplatesUpdateUser { get; set; } = new List<MailTemplates>();

    [JsonIgnore]  public virtual Roles Role { get; set; } = null!;

    [JsonIgnore]  public virtual ICollection<ShortUrls> ShortUrlsInsertUser { get; set; } = new List<ShortUrls>();

    [JsonIgnore]  public virtual ICollection<ShortUrls> ShortUrlsUpdateUser { get; set; } = new List<ShortUrls>();

    [JsonIgnore]  public virtual Users? UpdateUser { get; set; }
}
