using System.Text.Json.Serialization;

namespace PortalStudioUnicaAzione.Entities;

public partial class Logs
{
    public long Id { get; set; }

    public Guid? UserId { get; set; }

    public string? IpAddress { get; set; }

    public string? Endpoint { get; set; }

    public string? Message { get; set; }

    public string? RawData { get; set; }

    public DateTime Timestamp { get; set; }

    [JsonIgnore] public virtual Users? User { get; set; }
}
