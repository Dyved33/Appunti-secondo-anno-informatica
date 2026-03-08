using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Models;

namespace PortalStudioUnicaAzione.Repositories.Interfaces
{
    public interface IEventsRepository : IBaseCrudRepository<Events>
    {
        Task<TablesResponse<EventDto?>> GetTableList(TablesEventsRequest request);
        Task<Boolean> Enable(Guid id);
        Task<Boolean> Disable(Guid id);
    }
}