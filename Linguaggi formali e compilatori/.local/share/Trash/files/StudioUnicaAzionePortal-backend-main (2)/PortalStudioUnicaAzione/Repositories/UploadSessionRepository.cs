using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Repositories.Interfaces;

namespace PortalStudioUnicaAzione.Repositories
{
    public class UploadSessionRepository(PortalStudioUnicaAzioneContext appDbContext) : IUploadSessionRepository
    {
        public Task<UploadSession?> GetAsync(Guid uploadId, CancellationToken ct) =>
            appDbContext.UploadSessions.FirstOrDefaultAsync(x => x.UploadId == uploadId, ct);

        public Task<UploadSession?> GetByFileIdAsync(Guid fileId, CancellationToken ct) =>
            appDbContext.UploadSessions.FirstOrDefaultAsync(x => x.FileId == fileId, ct);

        public async Task InsertAsync(UploadSession session, CancellationToken ct)
        {
            appDbContext.UploadSessions.Add(session);
            await appDbContext.SaveChangesAsync(ct);
        }

        public async Task UpdateAsync(UploadSession session, CancellationToken ct)
        {
            appDbContext.UploadSessions.Update(session);
            await appDbContext.SaveChangesAsync(ct);
        }

        public async Task<int> ExpireOldAsync(DateTime utcNow, CancellationToken ct)
        {
            var old = await appDbContext.UploadSessions
                .Where(x => x.Status == 0 && x.ExpiresUtc < utcNow)
                .ExecuteUpdateAsync(x => x.SetProperty(s => s.Status, (byte)2), ct);
            return old;
        }
    }
}
