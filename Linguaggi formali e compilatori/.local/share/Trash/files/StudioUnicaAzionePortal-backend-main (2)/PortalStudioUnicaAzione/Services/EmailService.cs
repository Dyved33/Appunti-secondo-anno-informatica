using PortalStudioUnicaAzione.Models;
using PortalStudioUnicaAzione.Services.Interfaces;
using MailKit.Net.Smtp;
using MailKit.Security;
using Microsoft.Extensions.Options;
using MimeKit;

namespace PortalStudioUnicaAzione.Services
{
    public class EmailService(IOptions<SmtpSettings> smtpOptions) : IEmailService
    {
        private readonly SmtpSettings _settings = smtpOptions.Value;

        public async Task SendEmailAsync(
                List<String> to,
                List<String>? cc = null,
                List<String>? bcc = null,
                String subject = "",
                String body = "",
                List<(String fileName, Byte[] content)>? attachments = null)
        {
            var message = new MimeMessage();
            message.From.Add(new MailboxAddress(_settings.FromName, _settings.FromEmail));
            message.To.AddRange(to.Select(MailboxAddress.Parse));
            if (cc != null)
                message.Cc.AddRange(cc.Select(MailboxAddress.Parse));
            if (bcc != null)
                message.Bcc.AddRange(bcc.Select(MailboxAddress.Parse));
            message.Subject = subject;
            var builder = new BodyBuilder { HtmlBody = body };
            if (attachments != null)
            {
                foreach (var (fileName, content) in attachments)
                {
                    builder.Attachments.Add(fileName, content);
                }
            }
            message.Body = builder.ToMessageBody();
            int retries = 0;
            while (true)
            {
                try
                {
                    using var client = new SmtpClient();
                    await client.ConnectAsync(_settings.Host, _settings.Port, _settings.UseSsl ? SecureSocketOptions.SslOnConnect : SecureSocketOptions.StartTls);
                    await client.AuthenticateAsync(_settings.Username, _settings.Password);
                    await client.SendAsync(message);
                    await client.DisconnectAsync(true);
                    break;
                }
                catch (Exception ex)
                {
                    retries++;
                    if (retries >= _settings.RetryCount)
                        throw new InvalidOperationException("Failed to send email after retries.", ex);
                    await Task.Delay(_settings.RetryDelaySeconds * 1000);
                }
            }
        }
    }
}
