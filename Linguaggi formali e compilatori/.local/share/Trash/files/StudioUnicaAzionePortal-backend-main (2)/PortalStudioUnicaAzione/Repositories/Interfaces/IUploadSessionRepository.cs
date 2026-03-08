using PortalStudioUnicaAzione.Entities;

namespace PortalStudioUnicaAzione.Repositories.Interfaces
{
    public interface IUploadSessionRepository
    {
        Task<UploadSession?> GetAsync(Guid uploadId, CancellationToken ct);
        Task<UploadSession?> GetByFileIdAsync(Guid fileId, CancellationToken ct);
        Task InsertAsync(UploadSession session, CancellationToken ct);
        Task UpdateAsync(UploadSession session, CancellationToken ct);
        Task<int> ExpireOldAsync(DateTime utcNow, CancellationToken ct);
    }
}
