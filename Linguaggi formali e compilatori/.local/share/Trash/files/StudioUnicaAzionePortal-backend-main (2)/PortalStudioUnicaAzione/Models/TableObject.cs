namespace PortalStudioUnicaAzione.Models
{
    public class TablesRequest
    {
        public String OrderColumn { get; set; } = "Id";
        public String OrderDirection { get; set; } = "ASC";
        public Int32 PageNumber { get; set; } = 1;
        public Int32 RecordsPerPage { get; set; } = 100;
    }

    public class TablesResponse<T>
    {
        public List<T> Records { get; set; } = [];
        public Int32 RecordsTotal { get; set; }
    }
}
