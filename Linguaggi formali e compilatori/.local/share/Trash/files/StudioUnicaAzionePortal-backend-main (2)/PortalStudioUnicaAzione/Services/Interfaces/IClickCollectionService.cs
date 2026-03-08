using PortalStudioUnicaAzione.Dtos;

namespace PortalStudioUnicaAzione.Services.Interfaces
{
    public interface IClickCollectionService
    {
        Task CollectClickAsync(ClickCollectionData clickData);
        Task<int> GetPendingClicksCountAsync();
    }
}