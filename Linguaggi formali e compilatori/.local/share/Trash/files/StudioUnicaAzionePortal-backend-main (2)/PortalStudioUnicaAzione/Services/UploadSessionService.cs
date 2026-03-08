using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Handlers;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using PortalStudioUnicaAzione.Services.Interfaces;

namespace PortalStudioUnicaAzione.Services
{
    public class UploadSessionService(
    IUploadSessionRepository uploadRepo,
    IMediaIndexRepository mediaRepo,
    IBlobStorageService blob,
    IDbLogger dbLogger // your existing logger
) : IUploadSessionService
    {
        public async Task<(bool Exists, Guid? FileId)> ExistsByHashAsync(Guid eventId, string sha256Hex, CancellationToken ct)
        {
            var sha = HexToBytes(sha256Hex);
            var existing = await mediaRepo.GetByEventAndHashAsync(eventId, sha, ct);
            return (existing != null, existing?.FileId);
        }

        public async Task<UploadInitResponse> CreateAsync(Guid eventId, UploadInitRequest req, CancellationToken ct)
        {
            var fileId = Uuid7.NewGuid(); // UUIDv7 in prod: use a UUIDv7 generator
            var uploadId = Uuid7.NewGuid();
            var blobPath = blob.GetOriginalPath(eventId, fileId);

            var sas = await blob.GetWriteSasForBlobAsync(blobPath, TimeSpan.FromMinutes(15), ct);

            var session = new UploadSession
            {
                UploadId = uploadId,
                EventId = eventId,
                FileId = fileId,
                Sha256 = HexToBytes(req.Sha256),
                SizeBytes = req.SizeBytes,
                MimeType = req.MimeType,
                OriginalFilename = req.OriginalFilename,
                CreatedUtc = DateTime.UtcNow,
                ExpiresUtc = DateTime.UtcNow.AddMinutes(20),
                Status = 0
            };
            await uploadRepo.InsertAsync(session, ct);

            await dbLogger.LogEvent(
    "UploadInitIssued",
    new { eventId, uploadId, fileId }
);

            return new UploadInitResponse()
            {
                UploadId = uploadId,
                FileId = fileId,
                PutSasUrl = sas,
                RequiredHeaders = new Dictionary<string, string>
                {
                    ["x-ms-blob-type"] = "BlockBlob",
                    ["Content-Type"] = req.MimeType
                }
            };
        }

        public async Task<(long SizeBytes, string MimeType)> CompleteAsync(Guid eventId, UploadCompleteRequest req, string? ip, string? userAgent, CancellationToken ct)
        {
            var session = await uploadRepo.GetAsync(req.UploadId, ct);
            if (session == null || session.FileId != req.FileId || session.EventId != eventId)
                throw new InvalidOperationException("Upload session mismatch.");

            // Blob exist & properties
            var path = blob.GetOriginalPath(eventId, req.FileId);
            var probe = await blob.TryGetBlobAsync(path, ct);
            if (!probe.Exists) throw new FileNotFoundException("Blob not found.");

            // MIME sniff (first bytes)
            var (stream, length, _, _, _) = await blob.OpenReadAsync(path, ct);
            if (stream == null) throw new IOException("Cannot open blob.");
            await using (stream.ConfigureAwait(false))
            {
                var sniffed = await MimeSniffer.SniffAsync(stream);
                var mime = sniffed ?? session.MimeType;
                await blob.SetBlobContentTypeAsync(path, mime, ct);
                await blob.SetBlobMetadataAsync(path, new Dictionary<string, string>
                {
                    ["originalName"] = session.OriginalFilename ?? "",
                    ["sha256"] = BytesToHex(session.Sha256)
                }, ct);

                session.Status = 1; // Completed
                await uploadRepo.UpdateAsync(session, ct);

                await dbLogger.LogEvent(
                    "UploadCompleted",
                    new { eventId, req.FileId, ip, userAgent, size = length, mime }
                );

                return (length, mime);
            }
        }

        private static byte[] HexToBytes(string hex)
        {
            if (hex.Length % 2 != 0) throw new ArgumentException("HEX length must be even.");
            var bytes = new byte[hex.Length / 2];
            for (int i = 0; i < bytes.Length; i++)
                bytes[i] = Convert.ToByte(hex.Substring(i * 2, 2), 16);
            return bytes;
        }
        private static string BytesToHex(byte[] bytes) => BitConverter.ToString(bytes).Replace("-", "").ToLowerInvariant();
    }
}
