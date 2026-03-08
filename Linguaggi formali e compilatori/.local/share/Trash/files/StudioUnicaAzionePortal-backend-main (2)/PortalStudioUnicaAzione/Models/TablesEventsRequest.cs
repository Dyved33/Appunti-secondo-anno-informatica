namespace PortalStudioUnicaAzione.Models
{
    public class TablesEventsRequest : TablesRequest
    {
        public String SearchTitle { get; set; } = String.Empty;
        public DateTime? SearchDateEventFrom { get; set; }
        public DateTime? SearchDateEventTo { get; set; }
        public Int32 SearchEnabled { get; set; } = -1;
    }
}