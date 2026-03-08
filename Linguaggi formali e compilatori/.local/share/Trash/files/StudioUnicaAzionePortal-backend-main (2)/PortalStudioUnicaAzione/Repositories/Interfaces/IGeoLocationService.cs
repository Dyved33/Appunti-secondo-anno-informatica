using PortalStudioUnicaAzione.Dtos;

namespace PortalStudioUnicaAzione.Repositories.Interfaces
{
    public interface IGeoLocationRepository
    {
        Task<ParsedGeoLocation?> GetGeoLocationAsync(string ipAddress);
        Task<Guid> GetOrCreateGeoInfoAsync(string ipAddress);
    }
}