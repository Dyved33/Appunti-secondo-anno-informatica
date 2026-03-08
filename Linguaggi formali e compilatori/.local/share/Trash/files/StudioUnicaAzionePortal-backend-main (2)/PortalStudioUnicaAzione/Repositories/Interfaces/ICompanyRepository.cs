using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Models;

namespace PortalStudioUnicaAzione.Repositories.Interfaces
{
    public interface ICompanyRepository : IBaseCrudRepository<Companies>
    {
        Task<TablesResponse<CompanyDto?>> GetTableList(TablesCompaniesRequest request);
        Task<IEnumerable<UserDto>> GetUsersByCompanyIdAsync(Guid companyId);
        Task Enable(Guid id);
        Task Disable(Guid id);
    }
}