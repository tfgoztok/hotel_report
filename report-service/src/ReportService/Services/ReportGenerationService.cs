using ReportService.Interfaces;
using ReportService.Models;
using System.Text.Json;

namespace ReportService.Services
{
    public class ReportGenerationService : IReportGenerationService
    {
        private readonly IReportRepository _reportRepository;

        public ReportGenerationService(IReportRepository reportRepository)
        {
            _reportRepository = reportRepository;
        }

        public async Task GenerateReport(string message)
        {
            // Deserialize the message
            var reportRequest = JsonSerializer.Deserialize<ReportRequest>(message);

            // TODO: Implement actual report generation logic here
            // For now, we'll just create a dummy report
            var report = new Report
            {
                Id = Guid.NewGuid(),
                RequestDate = DateTime.UtcNow,
                Status = "Completed",
                Location = reportRequest.Location,
                HotelCount = 10, // Dummy value
                PhoneNumberCount = 20 // Dummy value
            };

            await _reportRepository.AddAsync(report);
        }
    }

    public class ReportRequest
    {
        public string Location { get; set; }
    }
}