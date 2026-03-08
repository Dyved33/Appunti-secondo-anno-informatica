namespace PortalStudioUnicaAzione.Entities
{
    public class MediaIndex
    {
        public Guid FileId { get; set; }
        public Guid EventId { get; set; }
        public byte[] Sha256 { get; set; } = Array.Empty<byte>();
        public long? SizeBytes { get; set; }
        public DateTime CreatedUtc { get; set; }
    }
}
