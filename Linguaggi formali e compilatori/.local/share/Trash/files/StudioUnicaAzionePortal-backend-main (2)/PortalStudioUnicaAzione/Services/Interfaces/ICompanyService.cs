using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Models;

namespace PortalStudioUnicaAzione.Services.Interfaces
{
    public interface ICompanyService : IBaseCrudService<Companies>
    {
        Task<TablesResponse<CompanyDto?>> GetTableList(TablesCompaniesRequest request);
        Task<IEnumerable<UserDto>> GetUsersByCompanyIdAsync(Guid companyId);
        Task Enable(Guid id);
        Task Disable(Guid id);
    }
}
