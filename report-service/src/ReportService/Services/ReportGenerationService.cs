using System;
using System.Text.Json;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using ReportService.Interfaces;
using ReportService.Models;

namespace ReportService.Services
{
    // Class responsible for generating reports based on incoming messages
    public class ReportGenerationService : IReportGenerationService
    {
        // Dependencies for report generation
        private readonly IReportRepository _reportRepository; // Repository for report data
        private readonly ILogger<ReportGenerationService> _logger; // Logger for logging information and errors
        private readonly IGraphQLClient _graphQLClient; // Client for sending GraphQL queries

        // Constructor to initialize dependencies
        public ReportGenerationService(
            IReportRepository reportRepository,
            ILogger<ReportGenerationService> logger,
            IGraphQLClient graphQLClient)
        {
            _reportRepository = reportRepository; // Assigning the report repository
            _logger = logger; // Assigning the logger
            _graphQLClient = graphQLClient; // Assigning the GraphQL client
        }

        // Method to generate a report based on the incoming message
        // Parameters:
        // - message: JSON string containing report request details
        public async Task GenerateReport(string message)
        {
            Console.WriteLine("Starting report generation process");
            _logger.LogInformation($"Received message: {message}");

            try
            {
                var options = new JsonSerializerOptions
                {
                    PropertyNameCaseInsensitive = true
                };

                var reportRequest = JsonSerializer.Deserialize<ReportRequest>(message, options);
                _logger.LogInformation($"Deserialized report request: {JsonSerializer.Serialize(reportRequest)}");

                if (reportRequest != null)
                {
                    var report = new Report
                    {
                        ReportRequestId = reportRequest.Id,
                        RequestDate = DateTime.UtcNow,
                        Status = "Processing",
                        Location = reportRequest.Location ?? "Unknown",
                        HotelCount = 0,
                        PhoneNumberCount = 0
                    };

                    _logger.LogInformation($"Created initial report: {JsonSerializer.Serialize(report)}");

                    await _reportRepository.AddAsync(report);

                    try
                    {
                        var hotelsQuery = @"query($location: String!) { hotelsByLocation(location: $location) { id } }";
                        var contactsQuery = @"query($location: String!) { contactsByLocation(location: $location) { id Type } }";

                        var hotelsResult = await _graphQLClient.SendQueryAsync(hotelsQuery, new { location = reportRequest.Location });
                        var contactsResult = await _graphQLClient.SendQueryAsync(contactsQuery, new { location = reportRequest.Location });

                        _logger.LogInformation($"GraphQL hotels result: {hotelsResult}");
                        _logger.LogInformation($"GraphQL contacts result: {contactsResult}");

                        var hotelsData = JsonSerializer.Deserialize<HotelsQueryResult>(hotelsResult, options);
                        var contactsData = JsonSerializer.Deserialize<ContactsQueryResult>(contactsResult, options);

                        // Ensure hotelsData and contactsData are not null before accessing their properties
                        report.HotelCount = hotelsData?.Data?.HotelsByLocation?.Count ?? 0;
                        report.PhoneNumberCount = contactsData?.Data?.ContactsByLocation?.Count(c => c.Type?.Equals("PHONE", StringComparison.OrdinalIgnoreCase) == true) ?? 0;


                        report.Status = "Completed";

                        await _reportRepository.UpdateAsync(report);

                        _logger.LogInformation($"Updated report: {JsonSerializer.Serialize(report)}");
                    }
                    catch (Exception ex)
                    {
                        _logger.LogError(ex, "Error processing GraphQL queries");
                        report.Status = "Error";
                        await _reportRepository.UpdateAsync(report);
                    }
                }
                else
                {
                    _logger.LogWarning("Failed to deserialize report request: null result");
                    await CreateErrorReport("Failed to deserialize report request", "Unknown");
                }
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error generating report: {Message}", ex.Message);
                await CreateErrorReport($"Error generating report: {ex.Message}", "Unknown");
            }
        }

        // Method to create an error report
        private async Task CreateErrorReport(string errorMessage, string location = "Unknown")
        {
            var errorReport = new Report
            {
                Id = Guid.NewGuid().ToString(),
                RequestDate = DateTime.UtcNow,
                Status = "Error",
                Location = location,
                HotelCount = 0,
                PhoneNumberCount = 0
            };
            await _reportRepository.AddAsync(errorReport);
        }
    }

    // Class representing the incoming report request
    public class ReportRequest
    {
        public Guid Id { get; set; } // Unique identifier for the report
        public string? Status { get; set; } // Status of the report
        public string? Location { get; set; } // Location for the report
    }

    // Class representing the result of the hotels query
    public class HotelsQueryResult
    {
        public HotelsData? Data { get; set; } // Data containing hotel information
    }

    // Class representing the data structure for hotels
    public class HotelsData
    {
        public List<HotelResult>? HotelsByLocation { get; set; } // List of hotels by location
    }

    // Class representing individual hotel results
    public class HotelResult
    {
        public string? Id { get; set; } // Unique identifier for the hotel
    }

    // Class representing the result of the contacts query
    public class ContactsQueryResult
    {
        public ContactsData? Data { get; set; } // Data containing contact information
    }

    // Class representing the data structure for contacts
    public class ContactsData
    {
        public List<ContactResult>? ContactsByLocation { get; set; } // List of contacts by location
    }

    // Class representing individual contact results
    public class ContactResult
    {
        public string? Id { get; set; } // Unique identifier for the contact
        public string? Type { get; set; } // Type of the contact (e.g., PHONE)
    }
}
