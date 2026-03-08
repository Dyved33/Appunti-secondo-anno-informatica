using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Models;
using PortalStudioUnicaAzione.Services.Interfaces;

namespace PortalStudioUnicaAzione.Repositories
{
    public class CompanyRepository : BaseCrudRepository<Companies>, ICompanyRepository
    {
        private readonly PortalStudioUnicaAzioneContext _dbContext;
        private readonly IDbLogger _logger;

        public CompanyRepository(PortalStudioUnicaAzioneContext db, IDbLogger logger) : base(db)
        {
            _dbContext = db;
            _logger = logger;
        }

        public async Task<TablesResponse<CompanyDto?>> GetTableList(TablesCompaniesRequest request)
        {
            var response = new TablesResponse<CompanyDto?>();
            try
            {
                request ??= new TablesCompaniesRequest();
                int pagenum = pagenum = (request?.PageNumber ?? 0) - 1;
                if (pagenum < 0) pagenum = 0;
                int pagerecords = pagerecords = request?.RecordsPerPage ?? 0;
                if (pagerecords < 1) pagerecords = 100;
                int start = pagenum * pagerecords;
                int length = pagerecords;
                // esegui la query
                var list = _dbContext.Companies
                                            .Where(x =>
                                                (String.IsNullOrEmpty(request!.SearchCompanyName) || String.IsNullOrEmpty(x.CompanyName) || x.CompanyName.Contains(request.SearchCompanyName)) &&
                                                (request.SearchEnabled == -1 || (request.SearchEnabled == 1 && x.Enabled) || (request.SearchEnabled == 0 && !x.Enabled)) &&
                                                x.Deleted == false)
                                            .AsQueryable();
                if (request!.OrderDirection.Equals("ASC", StringComparison.CurrentCultureIgnoreCase))
                {
                    list = list.OrderBy(x => EF.Property<object>(x, request!.OrderColumn));
                }
                else
                {
                    list = list.OrderByDescending(x => EF.Property<object>(x, request!.OrderColumn));
                }
                var results = await list.Skip(start).Take(length).ToListAsync();
                response.RecordsTotal = results.Count;
                response.Records = results.Select(CompanyDto.ConvertToTableDto).ToList() ?? [];
            }
            catch (Exception ex)
            {
                await _logger.LogException(ex);
            }
            return response;
        }

        public async Task<IEnumerable<UserDto>> GetUsersByCompanyIdAsync(Guid companyId)
        {
            var users = await _dbContext.Users
                                            .Where(u => u.CompanyId == companyId && u.Enabled && !u.Deleted)
                                            .OrderBy(u => u.Lastname)
                                            .ThenBy(u => u.Firstname)
                                            .ToListAsync();
            return users
                .Select(UserDto.ConvertToDto)
                .Where(dto => dto != null)!
                .Cast<UserDto>();
        }

        public async Task Enable(Guid id)
        {
            try
            {
                var entity = await _dbContext.Companies.FirstOrDefaultAsync(e => e.Id == id);
                if (entity == null)
                    return;
                entity.Enabled = true;
                _dbContext.Update(entity);
                await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                await _logger.LogException(ex);
            }
        }

        public async Task Disable(Guid id)
        {
            try
            {
                var entity = await _dbContext.Companies.FirstOrDefaultAsync(e => e.Id == id);
                if (entity == null)
                    return;
                entity.Enabled = false;
                _dbContext.Update(entity);
                await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                await _logger.LogException(ex);
            }
        }
    }
}