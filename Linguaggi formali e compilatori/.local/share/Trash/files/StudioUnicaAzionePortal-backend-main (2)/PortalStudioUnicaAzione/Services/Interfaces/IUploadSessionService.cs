using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;

namespace PortalStudioUnicaAzione.Services.Interfaces
{
    public interface IUploadSessionService
    {
        Task<(bool Exists, Guid? FileId)> ExistsByHashAsync(Guid eventId, string sha256Hex, CancellationToken ct);
        Task<UploadInitResponse> CreateAsync(Guid eventId, UploadInitRequest req, CancellationToken ct);
        Task<(long SizeBytes, string MimeType)> CompleteAsync(Guid eventId, UploadCompleteRequest req, string? ip, string? userAgent, CancellationToken ct);
    }
}
