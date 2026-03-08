using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Repositories.Interfaces;

namespace PortalStudioUnicaAzione.Repositories
{
    public class MediaIndexRepository(PortalStudioUnicaAzioneContext appDbContext) : IMediaIndexRepository
    {
        public Task<MediaIndex?> GetByFileIdAsync(Guid fileId, CancellationToken ct) =>
            appDbContext.MediaIndex.AsNoTracking().FirstOrDefaultAsync(x => x.FileId == fileId, ct);

        public Task<MediaIndex?> GetByEventAndHashAsync(Guid eventId, byte[] sha256, CancellationToken ct) =>
            appDbContext.MediaIndex.AsNoTracking().FirstOrDefaultAsync(x => x.EventId == eventId && x.Sha256 == sha256, ct);

        public async Task InsertAsync(MediaIndex entity, CancellationToken ct)
        {
            appDbContext.MediaIndex.Add(entity);
            await appDbContext.SaveChangesAsync(ct);
        }

        public Task<List<MediaIndex>> ListByEventAsync(Guid eventId, CancellationToken ct) =>
            appDbContext.MediaIndex.AsNoTracking()
              .Where(x => x.EventId == eventId)
              .OrderByDescending(x => x.CreatedUtc)
              .ToListAsync(ct);
    }
}
