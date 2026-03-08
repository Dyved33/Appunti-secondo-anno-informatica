namespace PortalStudioUnicaAzione.Services.Interfaces
{
    public interface IDbLogger
    {
        Task LogEvent(String message, Object? rawData = null);
        Task LogException(Exception ex, Object? rawData = null);
    }
}
