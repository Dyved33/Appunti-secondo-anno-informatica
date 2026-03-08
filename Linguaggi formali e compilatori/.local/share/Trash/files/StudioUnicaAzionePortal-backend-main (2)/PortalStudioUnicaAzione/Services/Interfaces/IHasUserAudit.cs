namespace PortalStudioUnicaAzione.Services.Interfaces
{
    public interface IHasUserAudit
    {
        Guid? InsertUserId { get; set; }
        Guid? UpdateUserId { get; set; }
    }
}
