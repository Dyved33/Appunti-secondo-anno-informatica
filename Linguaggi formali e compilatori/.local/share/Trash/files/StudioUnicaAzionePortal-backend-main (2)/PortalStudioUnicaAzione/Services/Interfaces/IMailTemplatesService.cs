using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Models;

namespace PortalStudioUnicaAzione.Services.Interfaces
{
    public interface IMailTemplatesService : IBaseCrudService<MailTemplates>
    {
        Task<TablesResponse<MailTemplateDto?>> GetTableList(TablesMailTemplatesRequest request);
        Task<List<SelectOption<Guid>>> GetSelectList(SelectRequest request);
        Task<MailTemplates?> GetByTitle(String title);
    }
}