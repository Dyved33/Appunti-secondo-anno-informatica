using System.Text.Json.Serialization;

namespace PortalStudioUnicaAzione.Entities;

public partial class MailTemplates : BaseEntity
{
    public string? Title { get; set; }

    public string? Subject { get; set; }

    public string? Body { get; set; }

    [JsonIgnore] public virtual Users? InsertUser { get; set; }

    [JsonIgnore] public virtual Users? UpdateUser { get; set; }
}
