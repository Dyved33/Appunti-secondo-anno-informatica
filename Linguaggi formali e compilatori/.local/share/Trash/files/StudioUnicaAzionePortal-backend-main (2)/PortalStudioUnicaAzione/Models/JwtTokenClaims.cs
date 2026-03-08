namespace PortalStudioUnicaAzione.Models
{
    public class JwtTokenClaims
    {
        public Guid? UserId { get; set; } = null!;
        public String? Fullname { get; set; }
        public String? Username { get; set; } = null!;
        public String? Email { get; set; }
        public String? Role { get; set; }
        public String? CompanyId { get; set; }
        public String? CompanyName { get; set; }
    }
}
