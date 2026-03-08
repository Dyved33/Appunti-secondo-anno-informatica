using PortalStudioUnicaAzione.Dtos;

namespace PortalStudioUnicaAzione.Repositories.Interfaces
{
    public interface IReferrerParsingRepository
    {
        Task<ParsedReferrer> ParseReferrer(string? referrerUrl);
        Task<Guid> GetOrCreateReferrerInfoAsync(string? referrerUrl, string? acceptLanguage);
    }
}