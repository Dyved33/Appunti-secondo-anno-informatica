using PortalStudioUnicaAzione.Dtos;

namespace PortalStudioUnicaAzione.Repositories.Interfaces
{
    public interface IAnalyticsRepository
    {
        Task<AnalyticsResponseDto> GetAnalyticsSummaryAsync(Guid shortUrlId, DateTime? fromDate = null, DateTime? toDate = null);
        Task<AnalyticsResponseDto> GetAnalyticsDetailedAsync(Guid shortUrlId, DateTime? fromDate = null, DateTime? toDate = null);
        Task<AnalyticsResponseDto> GetAnalyticsByShortCodeAsync(string shortCode, bool includeDetails = false, DateTime? fromDate = null, DateTime? toDate = null);
        Task<List<DailyClickDto>> GetDailyStatsAsync(Guid shortUrlId, DateTime? fromDate = null, DateTime? toDate = null);
        Task<List<HourlyClickDto>> GetHourlyStatsAsync(Guid shortUrlId, DateTime? fromDate = null, DateTime? toDate = null);
        Task<bool> HasAnalyticsDataAsync(Guid shortUrlId);
    }
}
