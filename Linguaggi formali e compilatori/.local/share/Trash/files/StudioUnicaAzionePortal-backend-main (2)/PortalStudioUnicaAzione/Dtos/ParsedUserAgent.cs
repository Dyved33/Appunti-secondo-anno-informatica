namespace PortalStudioUnicaAzione.Dtos
{
    public class ParsedUserAgent
    {
        public string DeviceType { get; set; } = "Unknown";
        public string? DeviceModel { get; set; }
        public string? DeviceBrand { get; set; }
        public string? OS { get; set; }
        public string? OSVersion { get; set; }
        public string? Browser { get; set; }
        public string? BrowserVersion { get; set; }
        public bool IsBot { get; set; }
        public string? BotType { get; set; }
    }
}
