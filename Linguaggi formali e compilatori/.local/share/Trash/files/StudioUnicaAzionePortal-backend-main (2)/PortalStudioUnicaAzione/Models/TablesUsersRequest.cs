namespace PortalStudioUnicaAzione.Models
{
    public class TablesUsersRequest : TablesRequest
    {
        public String SearchLastname { get; set; } = String.Empty;
        public String SearchFirstname { get; set; } = String.Empty;
        public String SearchUsername { get; set; } = String.Empty;
        public String SearchEmail { get; set; } = String.Empty;
        public String SearchRole { get; set; } = String.Empty;
        public String SearchCompanyName { get; set; } = String.Empty;
        public Int32 SearchEnabled { get; set; } = -1;
    }
}