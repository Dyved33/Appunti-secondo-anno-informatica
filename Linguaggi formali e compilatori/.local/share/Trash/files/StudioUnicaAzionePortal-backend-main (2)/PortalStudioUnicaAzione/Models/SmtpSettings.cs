namespace PortalStudioUnicaAzione.Models
{
    public class SmtpSettings
    {
        public String Host { get; set; } = String.Empty;
        public Int32 Port { get; set; }
        public Boolean UseSsl { get; set; }
        public String Username { get; set; } = String.Empty;
        public String Password { get; set; } = String.Empty;
        public String FromEmail { get; set; } = String.Empty;
        public String FromName { get; set; } = String.Empty;
        public Int32 RetryCount { get; set; } = 3;
        public Int32 RetryDelaySeconds { get; set; } = 5;
    }
}
