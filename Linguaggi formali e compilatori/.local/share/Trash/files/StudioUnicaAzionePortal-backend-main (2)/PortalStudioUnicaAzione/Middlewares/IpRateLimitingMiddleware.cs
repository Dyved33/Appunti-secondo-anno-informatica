using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Caching.Memory;
using PortalStudioUnicaAzione.Entities;

namespace PortalStudioUnicaAzione.Middlewares
{
    public class IpRateLimitingMiddleware
    {
        private readonly RequestDelegate _next;
        private readonly HttpContext? _httpContextObj;
        private readonly IDbContextFactory<PortalStudioUnicaAzioneContext> _factory;
        private static readonly MemoryCache _memoryCache = new(new MemoryCacheOptions());
        private const int Limit = 50; // Max requests
        private static readonly TimeSpan Period = TimeSpan.FromMinutes(1); // Time window

        public IpRateLimitingMiddleware(RequestDelegate next,
            IHttpContextAccessor httpContextAccessor,
            IDbContextFactory<PortalStudioUnicaAzioneContext> factory)
        {
            _next = next;
            _httpContextObj = httpContextAccessor.HttpContext;
            _factory = factory;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            var ip = context.Connection.RemoteIpAddress?.MapToIPv4().ToString() ?? "unknown";

            var cacheKey = $"ip:{ip}";
            var requestCount = _memoryCache.GetOrCreate(cacheKey, entry =>
            {
                entry.AbsoluteExpirationRelativeToNow = Period;
                return 0;
            });

            if (requestCount >= Limit)
            {
                await using var db = await _factory.CreateDbContextAsync();
                Logs entry = new()
                {
                    UserId = null,
                    IpAddress = _httpContextObj?.Connection?.RemoteIpAddress?.ToString() ?? _httpContextObj?.Request?.Headers["X-Forwarded-For"].ToString(),
                    Endpoint = $"{_httpContextObj?.Request?.Scheme}://{_httpContextObj?.Request?.Host}",
                    Message = $"IP {ip} blocked due to too many requests.",
                    RawData = null,
                    Timestamp = DateTime.Now
                };
                await db.Logs.AddAsync(entry);
                await db.SaveChangesAsync();
                context.Response.StatusCode = StatusCodes.Status429TooManyRequests;
                await context.Response.WriteAsync("Too many requests. Try again later.");
                return;
            }

            _memoryCache.Set(cacheKey, requestCount + 1, Period);

            await _next(context);
        }
    }
}
