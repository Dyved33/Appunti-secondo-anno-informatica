namespace PortalStudioUnicaAzione.Dtos
{
    public class FileDto
    {
        public Guid FileId { get; set; } = Guid.Empty;
        public String FileName { get; set; } = String.Empty;
        public String FileThumbPath { get; set; } = String.Empty;
        public String FileOriginalPath { get; set; } = String.Empty;
        public DateTime? FileLinkExpiration { get; set; }
    }
}
