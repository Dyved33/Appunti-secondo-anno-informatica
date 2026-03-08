using PortalStudioUnicaAzione.Dtos;

namespace PortalStudioUnicaAzione.Repositories.Interfaces
{
    public interface IClickCollectionRepository
    {
        Task CollectClickAsync(ClickCollectionData clickData);
        Task<int> GetPendingClicksCountAsync();
    }
}