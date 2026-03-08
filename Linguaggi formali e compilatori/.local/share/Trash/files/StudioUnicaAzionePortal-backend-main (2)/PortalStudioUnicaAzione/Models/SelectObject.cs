namespace PortalStudioUnicaAzione.Models
{
    public class SelectRequest
    {

    }

    public class SelectOption<T>
    {
        public required T Id { get; set; }
        public required String Label { get; set; }
        public String? FilterProp1 { get; set; }
        public String? FilterProp2 { get; set; }
    }

}
