using System;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Moq;
using ReportService.Interfaces;
using ReportService.Models;
using ReportService.Services;
using System.Collections.Generic;
using Xunit;
using System.Text.Json;

namespace ReportService.Tests
{
    public class ReportGenerationServiceTests
    {
        private readonly Mock<IReportRepository> _mockRepo;
        private readonly Mock<ILogger<ReportGenerationService>> _mockLogger;
        private readonly Mock<IGraphQLClient> _mockGraphQLClient;
        private readonly ReportGenerationService _service;
        private readonly List<string> _logMessages;

        public ReportGenerationServiceTests()
        {
            _mockRepo = new Mock<IReportRepository>();
            _mockLogger = new Mock<ILogger<ReportGenerationService>>();
            _mockGraphQLClient = new Mock<IGraphQLClient>();
            _logMessages = new List<string>();

            _mockLogger.Setup(x => x.Log(
                It.IsAny<LogLevel>(),
                It.IsAny<EventId>(),
                It.Is<It.IsAnyType>((v, t) => true),
                It.IsAny<Exception>(),
                (Func<It.IsAnyType, Exception, string>)It.IsAny<object>()))
            .Callback(new InvocationAction(invocation =>
            {
                var logLevel = (LogLevel)invocation.Arguments[0];
                var eventId = (EventId)invocation.Arguments[1];
                var state = invocation.Arguments[2];
                var exception = (Exception)invocation.Arguments[3];
                var formatter = invocation.Arguments[4];

                var invokeMethod = formatter.GetType().GetMethod("Invoke");
                var logMessage = (string)invokeMethod.Invoke(formatter, new[] { state, exception });
                _logMessages.Add($"{logLevel}: {logMessage}");
            }));

            _service = new ReportGenerationService(_mockRepo.Object, _mockLogger.Object, _mockGraphQLClient.Object);
        }


        [Fact]
        public async Task GenerateReport_WithInvalidMessage_LogsError()
        {
            // Arrange: Set up an invalid message
            var message = "invalid json";

            // Act: Call the method under test
            await _service.GenerateReport(message);

            // Assert: Verify that an error was logged
            _mockLogger.Verify(
                x => x.Log(
                    LogLevel.Error,
                    It.IsAny<EventId>(),
                    It.Is<It.IsAnyType>((o, t) => o.ToString().Contains("Error generating report")),
                    It.IsAny<Exception>(),
                    It.IsAny<Func<It.IsAnyType, Exception, string>>()),
                Times.Once);

            // Assert: Verify that an error report was added to the repository
            _mockRepo.Verify(repo => repo.AddAsync(It.Is<Report>(r =>
                r.Status == "Error" &&
                r.Location == "Unknown")), Times.Once);
        }

        [Fact]
        public async Task GenerateReport_WithInvalidMessage_CreatesErrorReport()
        {
            // Arrange: Prepare an invalid message
            var invalidMessage = "invalid json";

            // Capture the report that will be added to the repository
            Report capturedReport = null;
            _mockRepo.Setup(repo => repo.AddAsync(It.IsAny<Report>()))
                .Callback<Report>(r => capturedReport = r)
                .Returns(Task.CompletedTask);

            // Act: Call the method under test
            await _service.GenerateReport(invalidMessage);

            // Assert: Verify the error report was created with expected values
            Assert.NotNull(capturedReport);
            Assert.Equal("Unknown", capturedReport.Location);
            Assert.Equal("Error", capturedReport.Status);
        }

        [Fact]
        public async Task GenerateReport_WithValidMessage_GeneratesReportAndAddsToRepository()
        {
            // Arrange: Create a valid message for report generation
            var validMessage = JsonSerializer.Serialize(new ReportRequest
            {
                Id = Guid.NewGuid(),
                Status = "Pending",
                Location = "New York"
            });

            // Set up mock responses for GraphQL queries
            var hotelsQueryResult = JsonSerializer.Serialize(new HotelsQueryResult
            {
                Data = new HotelsData
                {
                    HotelsByLocation = new List<HotelResult>
                    {
                        new HotelResult { Id = "1" },
                        new HotelResult { Id = "2" }
                    }
                }
            });

            var contactsQueryResult = JsonSerializer.Serialize(new ContactsQueryResult
            {
                Data = new ContactsData
                {
                    ContactsByLocation = new List<ContactResult>
                    {
                        new ContactResult { Id = "1", Type = "PHONE" },
                        new ContactResult { Id = "2", Type = "EMAIL" },
                        new ContactResult { Id = "3", Type = "PHONE" }
                    }
                }
            });

            // Mock the GraphQL client to return predefined results
            _mockGraphQLClient.Setup(client => client.SendQueryAsync(It.IsAny<string>(), It.IsAny<object>()))
                .ReturnsAsync((string query, object variables) =>
                {
                    if (query.Contains("hotelsByLocation"))
                        return hotelsQueryResult;
                    else if (query.Contains("contactsByLocation"))
                        return contactsQueryResult;
                    return null;
                });

            // Capture the report that will be added to the repository
            Report capturedReport = null;
            _mockRepo.Setup(repo => repo.AddAsync(It.IsAny<Report>()))
                .Callback<Report>(r => capturedReport = r)
                .Returns(Task.CompletedTask);

            _mockRepo.Setup(repo => repo.UpdateAsync(It.IsAny<Report>()))
                .Callback<Report>(r => capturedReport = r)
                .Returns(Task.CompletedTask);

            // Act: Generate the report
            await _service.GenerateReport(validMessage);

            // Assert: Verify the report was created with expected values
            Assert.NotNull(capturedReport);
            Assert.Equal("New York", capturedReport.Location);
            Assert.Equal("Completed", capturedReport.Status);
            Assert.Equal(2, capturedReport.HotelCount);
            Assert.Equal(2, capturedReport.PhoneNumberCount);

            // Verify interactions with the repository and GraphQL client
            _mockRepo.Verify(repo => repo.AddAsync(It.IsAny<Report>()), Times.Once);
            _mockRepo.Verify(repo => repo.UpdateAsync(It.IsAny<Report>()), Times.Once);
            _mockGraphQLClient.Verify(client => client.SendQueryAsync(It.IsAny<string>(), It.IsAny<object>()), Times.Exactly(2));

            // Verify that the expected log messages were generated
            Assert.Contains(_logMessages, msg => msg.Contains("Deserialized report request"));
            Assert.Contains(_logMessages, msg => msg.Contains("Created initial report"));
            Assert.Contains(_logMessages, msg => msg.Contains("GraphQL hotels result"));
            Assert.Contains(_logMessages, msg => msg.Contains("GraphQL contacts result"));
            Assert.Contains(_logMessages, msg => msg.Contains("Updated report"));
        }
    }
}
