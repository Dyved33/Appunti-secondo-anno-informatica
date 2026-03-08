using Microsoft.EntityFrameworkCore;
using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Models;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using PortalStudioUnicaAzione.Services.Interfaces;

namespace PortalStudioUnicaAzione.Repositories
{
    public class EventsRepository : BaseCrudRepository<Events>, IEventsRepository
    {
        private readonly PortalStudioUnicaAzioneContext _dbContext;
        private readonly IDbLogger _logger;

        public EventsRepository(PortalStudioUnicaAzioneContext db, IDbLogger logger) : base(db)
        {
            _dbContext = db;
            _logger = logger;
        }

        public async Task<TablesResponse<EventDto?>> GetTableList(TablesEventsRequest request)
        {
            var response = new TablesResponse<EventDto?>();
            try
            {
                request ??= new TablesEventsRequest();
                int pagenum = pagenum = (request?.PageNumber ?? 0) - 1;
                if (pagenum < 0) pagenum = 0;
                int pagerecords = pagerecords = request?.RecordsPerPage ?? 0;
                if (pagerecords < 1) pagerecords = 100;
                int start = pagenum * pagerecords;
                int length = pagerecords;
                // esegui la query
                var list = _dbContext.Events
                                            .Where(x =>
                                                (String.IsNullOrEmpty(request!.SearchTitle) || String.IsNullOrEmpty(x.Title) || x.Title.Contains(request.SearchTitle)) &&
                                                (request.SearchDateEventFrom == null || x.DateEvent >= request.SearchDateEventFrom) &&
                                                (request.SearchDateEventTo == null || x.DateEvent <= request.SearchDateEventTo) &&
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
                response.Records = results.Select(EventDto.ConvertToTableDto).ToList() ?? [];
            }
            catch (Exception ex)
            {
                await _logger.LogException(ex);
            }
            return response;
        }

        public async Task<Boolean> Enable(Guid id)
        {
            int res = -1;
            try
            {
                Events? obj = await GetById(id);
                if (obj == null) return res > 0;
                obj.Enabled = true;
                _dbContext.Update(obj);
                res = await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                await _logger.LogException(ex);
            }
            return res > 0;
        }
        public async Task<Boolean> Disable(Guid id)
        {
            int res = -1;
            try
            {
                Events? obj = await GetById(id);
                if (obj == null) return res > 0;
                obj.Enabled = false;
                _dbContext.Update(obj);
                res = await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                await _logger.LogException(ex);
            }
            return res > 0;
        }
    }
}