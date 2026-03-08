namespace PortalStudioUnicaAzione.Dtos
{
    public class DeviceAnalyticsDto
    {
        public List<DeviceTypeDto> Types { get; set; } = new();
        public List<BrowserDto> Browsers { get; set; } = new();
        public List<OperatingSystemDto> OperatingSystems { get; set; } = new();
        public List<DeviceModelDto> Models { get; set; } = new();
    }

    public class DeviceTypeDto
    {
        public string DeviceType { get; set; } = null!;
        public int TotalClicks { get; set; }
        public int UniqueClicks { get; set; }
        public decimal Percentage { get; set; }
    }

    public class BrowserDto
    {
        public string Browser { get; set; } = null!;
        public string? Version { get; set; }
        public int TotalClicks { get; set; }
        public decimal Percentage { get; set; }
    }

    public class OperatingSystemDto
    {
        public string OS { get; set; } = null!;
        public string? Version { get; set; }
        public int TotalClicks { get; set; }
        public decimal Percentage { get; set; }
    }

    public class DeviceModelDto
    {
        public string? Brand { get; set; }
        public string? Model { get; set; }
        public string DeviceType { get; set; } = null!;
        public int TotalClicks { get; set; }
        public decimal Percentage { get; set; }
    }
}
