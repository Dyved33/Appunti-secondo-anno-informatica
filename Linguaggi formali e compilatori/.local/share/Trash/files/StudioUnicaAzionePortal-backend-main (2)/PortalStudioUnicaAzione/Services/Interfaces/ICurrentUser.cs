namespace PortalStudioUnicaAzione.Services.Interfaces
{
    public interface ICurrentUser
    {
        Guid? UserId { get; }
        String? Fullname { get; }
        String? Username { get; }
        String? Email { get; }
        String? Role { get; }
        String? CompanyId { get; }
        String? CompanyName { get; }
        String? IpAddress { get; }
    }
}
