using PortalStudioUnicaAzione.Services.Interfaces;
using System.Security.Claims;

namespace PortalStudioUnicaAzione.Services
{
    public class CurrentUser(IHttpContextAccessor accessor) : ICurrentUser
    {
        private readonly HttpContext? _ctx = accessor.HttpContext;
        private ClaimsPrincipal? U => _ctx?.User;

        public Guid? UserId => Guid.TryParse(U?.FindFirst("userid")?.Value, out var id) ? id : null;
        public string? Fullname => U?.FindFirst("fullname")?.Value ?? U?.Identity?.Name;
        public string? Username => U?.FindFirst("username")?.Value ?? U?.Identity?.Name;
        public string? Email => U?.FindFirst("email")?.Value;
        public string? Role => U?.FindFirst(ClaimTypes.Role)?.Value ?? U?.FindFirst("role")?.Value;
        public string? CompanyId => U?.FindFirst("companyid")?.Value;
        public string? CompanyName => U?.FindFirst("companyname")?.Value;
        public string? IpAddress
            => _ctx?.Connection?.RemoteIpAddress?.ToString()
               ?? _ctx?.Request?.Headers["X-Forwarded-For"].ToString();
    }
}
