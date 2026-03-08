using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Models;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using PortalStudioUnicaAzione.Services.Interfaces;
using System.Data;

namespace PortalStudioUnicaAzione.Repositories
{
    public class MailTemplatesRepository : BaseCrudRepository<MailTemplates>, IMailTemplatesRepository
    {
        private readonly PortalStudioUnicaAzioneContext _dbContext;
        private readonly IDbLogger _logger;

        public MailTemplatesRepository(PortalStudioUnicaAzioneContext db, IDbLogger logger) : base(db)
        {
            _dbContext = db;
            _logger = logger;
        }

        public async Task<TablesResponse<MailTemplateDto?>> GetTableList(TablesMailTemplatesRequest request)
        {
            var response = new TablesResponse<MailTemplateDto?>();
            try
            {
                request ??= new TablesMailTemplatesRequest();
                int pagenum = pagenum = (request?.PageNumber ?? 0) - 1;
                if (pagenum < 0) pagenum = 0;
                int pagerecords = pagerecords = request?.RecordsPerPage ?? 0;
                if (pagerecords < 1) pagerecords = 100;
                int start = pagenum * pagerecords;
                int length = pagerecords;
                // esegui la query
                var list = _dbContext.MailTemplates.Where(x =>
                                                (String.IsNullOrEmpty(request!.SearchTitle) || String.IsNullOrEmpty(x.Title) || x.Title.Contains(request.SearchTitle)) &&
                                                (String.IsNullOrEmpty(request!.SearchSubject) || String.IsNullOrEmpty(x.Subject) || x.Subject.Contains(request.SearchSubject)) &&
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
                response.Records = results.Select(MailTemplateDto.ConvertToTableDto).ToList() ?? [];
            }
            catch (Exception ex)
            {
                await _logger.LogException(ex);
            }
            return response;
        }

        public async Task<List<SelectOption<Guid>>> GetSelectList(SelectRequest request)
        {
            List<SelectOption<Guid>> list = [];
            try
            {
                list = await _dbContext.MailTemplates.Where(x => x.Deleted == false)
                                            .Select(x => new SelectOption<Guid>
                                            {
                                                Id = x.Id,
                                                Label = x.Title ?? ""
                                            })
                                            .ToListAsync();
            }
            catch (Exception ex)
            {
                await _logger.LogException(ex);
            }
            return list;
        }

        public async Task<MailTemplates?> GetByTitle(String title)
        {
            return await _dbContext.MailTemplates
                                    .FirstOrDefaultAsync(e => e.Title == title);
        }

    }
}
