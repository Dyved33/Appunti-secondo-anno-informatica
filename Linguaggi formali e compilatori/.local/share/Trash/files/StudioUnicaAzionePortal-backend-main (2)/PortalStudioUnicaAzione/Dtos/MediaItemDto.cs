namespace PortalStudioUnicaAzione.Dtos
{
    public sealed record MediaItemDto(
        Guid FileId,
        string OriginalFilename, /// <summary>"pending" | "processed"</summary>
        string Status,
        string MimeType,
        long SizeBytes,
        DateTimeOffset UploadedAtUtc
    );
}
