namespace PortalStudioUnicaAzione.Models
{
    public class TablesShortUrlsRequest : TablesRequest
    {
        public String SearchCompanyName { get; set; } = String.Empty;
        public String SearchOriginalUrl { get; set; } = String.Empty;
        public String SearchShortCode { get; set; } = String.Empty;
        public DateTime? ValidityFromDateStart { get; set; }
        public DateTime? ValidityFromDateEnd { get; set; }
        public DateTime? ValidityToDateStart { get; set; }
        public DateTime? ValidityToDateEnd { get; set; }
        public Int32 SearchEnabled { get; set; } = -1;
    }
}