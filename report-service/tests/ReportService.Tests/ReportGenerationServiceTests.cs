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

       /*  [Fact]
        public async Task GenerateReport_WithValidMessage_AddsReportToRepository()
        {
            // Arrange
            var message = "{\"Id\":\"12345\",\"Location\":\"Istanbul\",\"Status\":\"Pending\"}";
            var hotelsResponse = "{\"data\":{\"hotelsByLocation\":[{\"id\":\"1\"}]}}";
            var contactsResponse = "{\"data\":{\"contactsByLocation\":[{\"id\":\"1\",\"type\":\"PHONE\"}]}}";

            _mockGraphQLClient.SetupSequence(client => client.SendQueryAsync(It.IsAny<string>(), It.IsAny<object>()))
                .ReturnsAsync(hotelsResponse)
                .ReturnsAsync(contactsResponse);

            Report capturedInitialReport = null;
            Report capturedFinalReport = null;
            _mockRepo.Setup(repo => repo.AddAsync(It.IsAny<Report>()))
                .Callback<Report>(r => capturedInitialReport = r)
                .Returns(Task.CompletedTask);

            _mockRepo.Setup(repo => repo.UpdateAsync(It.IsAny<Report>()))
                .Callback<Report>(r => capturedFinalReport = r)
                .Returns(Task.CompletedTask);

            // Act
            await _service.GenerateReport(message);

            // Assert
            Assert.NotNull(capturedInitialReport);
            Assert.NotNull(capturedFinalReport);

            Console.WriteLine("Initial Report:");
            Console.WriteLine(JsonSerializer.Serialize(capturedInitialReport));

            Console.WriteLine("Final Report:");
            Console.WriteLine(JsonSerializer.Serialize(capturedFinalReport));

            Console.WriteLine("Logged messages:");
            foreach (var logMessage in _logMessages)
            {
                Console.WriteLine(logMessage);
            }

            Assert.Equal("Istanbul", capturedInitialReport.Location);
            Assert.Equal("Istanbul", capturedFinalReport.Location);
            Assert.Equal("Completed", capturedFinalReport.Status);
            Assert.Equal(1, capturedFinalReport.HotelCount);
            Assert.Equal(1, capturedFinalReport.PhoneNumberCount);

            Assert.Contains(_logMessages, m => m.Contains("Starting GenerateReport with message:"));
            Assert.Contains(_logMessages, m => m.Contains("Deserialized report request:"));
            Assert.Contains(_logMessages, m => m.Contains("Creating initial report with Location: Istanbul"));
            Assert.Contains(_logMessages, m => m.Contains("Added initial report to repository"));
            Assert.DoesNotContain(_logMessages, m => m.Contains("Error"));
        } */

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
    }
}
