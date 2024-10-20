namespace ReportService.Interfaces
{
    public interface IReportGenerationService
    {
        /// <summary>
        /// Generates a report based on the provided message.
        /// </summary>
        Task GenerateReport(string message);
    }
}
