using ReportService.Interfaces;
using ReportService.Models;
using System.Text.Json;

namespace ReportService.Services
{
    public class ReportGenerationService : IReportGenerationService
    {
        // Repository for accessing report data
        private readonly IReportRepository _reportRepository;
        // Logger for logging information and errors
        private readonly ILogger<ReportGenerationService> _logger;

        // Constructor to initialize the report repository and logger
        public ReportGenerationService(IReportRepository reportRepository, ILogger<ReportGenerationService> logger)
        {
            _reportRepository = reportRepository;
            _logger = logger;
        }

        // Method to generate a report based on the incoming message
        public async Task GenerateReport(string message)
        {
            // Log the received message
            _logger.LogInformation($"Received message: {message}");

            try
            {
                // Deserialize the message into a ReportRequest object
                var reportRequest = JsonSerializer.Deserialize<ReportRequest>(message);

                // Check if the deserialization was successful
                if (reportRequest != null)
                {
                    // Create a new report object with dummy values
                    var report = new Report
                    {
                        RequestDate = DateTime.UtcNow,
                        Status = "Completed",
                        Location = reportRequest.Location,
                        HotelCount = 10, // Dummy value
                        PhoneNumberCount = 20 // Dummy value
                    };

                    // Add the report to the repository
                    await _reportRepository.AddAsync(report);

                    // Log the successful report generation
                    _logger.LogInformation($"Report generated for location: {report.Location}");
                }
                else
                {
                    // Log a warning if deserialization fails
                    _logger.LogWarning("Failed to deserialize report request");
                }
            }
            catch (JsonException ex)
            {
                _logger.LogError(ex, "Error processing message");
                _logger.LogWarning("Failed to deserialize report request");
            }
            catch (Exception ex)
            {
                // Log any exceptions that occur during processing
                _logger.LogError(ex, "Error processing message");
            }
        }
    }

    // Class representing a report request
    public class ReportRequest
    {
        // Location for the report request
        public string Location { get; set; } = null!;
    }
}
