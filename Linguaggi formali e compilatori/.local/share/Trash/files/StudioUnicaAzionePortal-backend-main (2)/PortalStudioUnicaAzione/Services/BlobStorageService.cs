using Azure.Identity;
using Azure.Storage.Blobs;
using Azure.Storage.Blobs.Models;
using Azure.Storage.Sas;
using PortalStudioUnicaAzione.Services.Interfaces;

namespace PortalStudioUnicaAzione.Services
{
    public class BlobStorageService : IBlobStorageService
    {
        private readonly BlobServiceClient _svc;
        private readonly string _containerName = "main";

        public BlobStorageService(IConfiguration config)
        {
            // config["Storage:AccountUrl"] e.g. https://<account>.blob.core.windows.net
            _svc = new BlobServiceClient(new Uri(config["Storage:AccountUrl"]!), new DefaultAzureCredential());
        }

        public string GetOriginalPath(Guid eventId, Guid fileId) => $"events/{eventId}/{fileId}";
        public string GetThumbPath(Guid eventId, Guid fileId, string ext = "webp") => $"events/{eventId}/thumbs/{fileId}.{ext}";

        private BlobContainerClient GetContainer() => _svc.GetBlobContainerClient(_containerName);

        public async Task<Uri> GetWriteSasForBlobAsync(string blobPath, TimeSpan ttl, CancellationToken ct)
        {
            var container = GetContainer();
            var blob = container.GetBlobClient(blobPath);

            // Ensure container exists (or make sure it’s provisioned elsewhere)
            await container.CreateIfNotExistsAsync(PublicAccessType.None, cancellationToken: ct);

            // Use User Delegation SAS (Managed Identity → get Ud Key)
            var key = await _svc.GetUserDelegationKeyAsync(DateTimeOffset.UtcNow.AddMinutes(-5), DateTimeOffset.UtcNow.Add(ttl), ct);

            var sasBuilder = new BlobSasBuilder
            {
                BlobContainerName = container.Name,
                BlobName = blob.Name,
                Resource = "b",
                StartsOn = DateTimeOffset.UtcNow.AddMinutes(-2),
                ExpiresOn = DateTimeOffset.UtcNow.Add(ttl)
            };
            sasBuilder.SetPermissions(BlobSasPermissions.Create | BlobSasPermissions.Write | BlobSasPermissions.Add);

            var sas = sasBuilder.ToSasQueryParameters(key, _svc.AccountName);
            var uriBuilder = new UriBuilder(blob.Uri) { Query = sas.ToString() };
            return uriBuilder.Uri;
        }

        public async Task<(bool Exists, long? Size, string? ContentType, DateTimeOffset? LastModified)> TryGetBlobAsync(string blobPath, CancellationToken ct)
        {
            var blob = GetContainer().GetBlobClient(blobPath);
            if (!await blob.ExistsAsync(ct)) return (false, null, null, null);
            var props = await blob.GetPropertiesAsync(cancellationToken: ct);
            return (true, props.Value.ContentLength, props.Value.ContentType, props.Value.LastModified);
        }

        public async Task<(Stream? Stream, long Length, string ContentType, string ETag, DateTimeOffset LastModified)> OpenReadAsync(string blobPath, CancellationToken ct)
        {
            var blob = GetContainer().GetBlobClient(blobPath);
            if (!await blob.ExistsAsync(ct)) return (null, 0, "application/octet-stream", "", DateTimeOffset.UtcNow);
            var resp = await blob.DownloadStreamingAsync(cancellationToken: ct);
            return (resp.Value.Content, resp.Value.Details.ContentLength, resp.Value.Details.ContentType ?? "application/octet-stream", resp.Value.Details.ETag.ToString(), resp.Value.Details.LastModified);
        }

        public async Task<(Stream? Stream, long From, long To, long Length, string ETag, string ContentType)> OpenReadRangeAsync(string blobPath, long from, long to, CancellationToken ct)
        {
            var blob = GetContainer().GetBlobClient(blobPath);
            if (!await blob.ExistsAsync(ct)) return (null, 0, 0, 0, "", "application/octet-stream");

            var httpRange = new Azure.HttpRange(from, to - from + 1);
            var options = new BlobDownloadOptions
            {
                Range = httpRange
            };
            var resp = await blob.DownloadStreamingAsync(options, ct);
            var len = (await blob.GetPropertiesAsync(cancellationToken: ct)).Value.ContentLength;

            return (resp.Value.Content, from, to, len, resp.Value.Details.ETag.ToString(), resp.Value.Details.ContentType ?? "application/octet-stream");
        }

        public async Task SetBlobContentTypeAsync(string blobPath, string contentType, CancellationToken ct)
        {
            var blob = GetContainer().GetBlobClient(blobPath);
            var props = await blob.GetPropertiesAsync(cancellationToken: ct);
            await blob.SetHttpHeadersAsync(new BlobHttpHeaders
            {
                ContentType = contentType,
                ContentDisposition = props.Value.ContentDisposition
            }, cancellationToken: ct);
        }

        public Task SetBlobMetadataAsync(string blobPath, IDictionary<string, string> metadata, CancellationToken ct)
        {
            var blob = GetContainer().GetBlobClient(blobPath);
            return blob.SetMetadataAsync(metadata, cancellationToken: ct);
        }

        public async Task<bool> ThumbnailExistsAsync(Guid eventId, Guid fileId, CancellationToken ct)
        {
            var blob = GetContainer().GetBlobClient(GetThumbPath(eventId, fileId, "webp"));
            if (await blob.ExistsAsync(ct)) return true;
            var jpg = GetContainer().GetBlobClient(GetThumbPath(eventId, fileId, "jpg"));
            return await jpg.ExistsAsync(ct);
        }
    }
}
