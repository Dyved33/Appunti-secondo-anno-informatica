namespace PortalStudioUnicaAzione.Services.Interfaces
{
    public interface IEmailService
    {
        Task SendEmailAsync(
                List<String> to,
                List<String>? cc = null,
                List<String>? bcc = null,
                String subject = "",
                String body = "",
                List<(String fileName, Byte[] content)>? attachments = null);
    }
}