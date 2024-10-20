using System;
using System.Text.Json;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using ReportService.Interfaces;
using ReportService.Models;

namespace ReportService.Services
{
    public class ReportGenerationService : IReportGenerationService
    {
        // Dependencies for report generation
        private readonly IReportRepository _reportRepository; // Repository for report data
        private readonly ILogger<ReportGenerationService> _logger; // Logger for logging information and errors
        private readonly GraphQLClient _graphQLClient; // Client for making GraphQL queries

        // Constructor to initialize dependencies
        public ReportGenerationService(
            IReportRepository reportRepository,
            ILogger<ReportGenerationService> logger,
            GraphQLClient graphQLClient)
        {
            _reportRepository = reportRepository; // Assigning the report repository
            _logger = logger; // Assigning the logger
            _graphQLClient = graphQLClient; // Assigning the GraphQL client
        }

        // Method to generate a report based on the incoming message
        public async Task GenerateReport(string message)
        {
            _logger.LogInformation($"Received message: {message}"); // Log the received message

            try
            {
                // Deserialize the incoming message to a ReportRequest object
                var reportRequest = JsonSerializer.Deserialize<ReportRequest>(message);

                if (reportRequest != null) // Check if deserialization was successful
                {
                    // GraphQL queries to fetch hotels and contacts by location
                    var hotelsQuery = @"
                        query($location: String!) {
                            hotelsByLocation(location: $location) {
                                id
                            }
                        }";

                    var contactsQuery = @"
                        query($location: String!) {
                            contactsByLocation(location: $location) {
                                id
                                type
                            }
                        }";

                    // Send queries and get results
                    var hotelsResult = await _graphQLClient.SendQueryAsync(hotelsQuery, new { location = reportRequest.Location });
                    var contactsResult = await _graphQLClient.SendQueryAsync(contactsQuery, new { location = reportRequest.Location });

                    // Deserialize results into data models
                    var hotelsData = JsonSerializer.Deserialize<HotelsQueryResult>(hotelsResult);
                    var contactsData = JsonSerializer.Deserialize<ContactsQueryResult>(contactsResult);

                    // Count hotels and phone numbers
                    var hotelCount = hotelsData?.Data?.HotelsByLocation?.Count ?? 0; // Count of hotels
                    var phoneNumberCount = contactsData?.Data?.ContactsByLocation?.Count(c => c.Type == "PHONE") ?? 0; // Count of phone numbers

                    // Create a report object
                    var report = new Report
                    {
                        Id = reportRequest.Id, // Set report ID
                        RequestDate = DateTime.UtcNow, // Set the current date as request date
                        Status = "Completed", // Set report status
                        Location = reportRequest.Location, // Set report location
                        HotelCount = hotelCount, // Set hotel count
                        PhoneNumberCount = phoneNumberCount // Set phone number count
                    };

                    // Save the report to the repository
                    await _reportRepository.AddAsync(report); // Save the report asynchronously

                    _logger.LogInformation($"Report generated for location: {report.Location}"); // Log successful report generation
                }
                else
                {
                    _logger.LogWarning("Failed to deserialize report request"); // Log warning if deserialization fails
                }
            }
            catch (JsonException ex)
            {
                _logger.LogError(ex, "Error processing message"); // Log error for JSON processing issues
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error generating report"); // Log error for general report generation issues
            }
        }
    }

    // Class representing the report request structure
    public class ReportRequest
    {
        public Guid Id { get; set; } // Unique identifier for the report request
        public string Status { get; set; } // Status of the report request
        public string Location { get; set; } // Location for the report
    }

    // Class representing the result of hotel queries
    public class HotelsQueryResult
    {
        public HotelsData Data { get; set; } // Data containing hotel information
    }

    // Class representing the data structure for hotels
    public class HotelsData
    {
        public List<HotelResult> HotelsByLocation { get; set; } // List of hotels by location
    }

    // Class representing individual hotel results
    public class HotelResult
    {
        public string Id { get; set; } // Unique identifier for the hotel
    }

    // Class representing the result of contact queries
    public class ContactsQueryResult
    {
        public ContactsData Data { get; set; } // Data containing contact information
    }

    // Class representing the data structure for contacts
    public class ContactsData
    {
        public List<ContactResult> ContactsByLocation { get; set; } // List of contacts by location
    }

    // Class representing individual contact results
    public class ContactResult
    {
        public string Id { get; set; } // Unique identifier for the contact
        public string Type { get; set; } // Type of the contact (e.g., PHONE)
    }
}
