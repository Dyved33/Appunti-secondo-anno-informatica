using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Models;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using PortalStudioUnicaAzione.Services.Interfaces;

namespace PortalStudioUnicaAzione.Services
{
    public class CompanyService(ICompanyRepository companyRepository) : BaseCrudService<Companies>(companyRepository), ICompanyService
    {
        public async Task<TablesResponse<CompanyDto?>> GetTableList(TablesCompaniesRequest request)
        {
            return await companyRepository.GetTableList(request);
        }

        public Task<IEnumerable<UserDto>> GetUsersByCompanyIdAsync(Guid companyId)
        {
            return companyRepository.GetUsersByCompanyIdAsync(companyId);
        }

        public async Task Enable(Guid id)
        {
            await companyRepository.Enable(id);
        }

        public async Task Disable(Guid id)
        {
            await companyRepository.Disable(id);
        }
    }
}
