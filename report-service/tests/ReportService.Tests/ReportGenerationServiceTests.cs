using System;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Moq;
using ReportService.Interfaces;
using ReportService.Models;
using ReportService.Services;
using Xunit;

namespace ReportService.Tests
{
    public class ReportGenerationServiceTests
    {
        // Mock objects for testing
        private readonly Mock<IReportRepository> _mockRepo; // Mock for the report repository
        private readonly Mock<ILogger<ReportGenerationService>> _mockLogger; // Mock for the logger
        private readonly ReportGenerationService _service; // Instance of the service being tested

        public ReportGenerationServiceTests()
        {
            // Initialize mocks and the service instance
            _mockRepo = new Mock<IReportRepository>();
            _mockLogger = new Mock<ILogger<ReportGenerationService>>();
            _service = new ReportGenerationService(_mockRepo.Object, _mockLogger.Object);
        }

        [Fact]
        public async Task GenerateReport_WithValidMessage_AddsReportToRepository()
        {
            // Arrange
            var message = "{\"Location\":\"Istanbul\"}"; // Valid JSON message

            // Act
            await _service.GenerateReport(message); // Call the method under test

            // Assert
            // Verify that the AddAsync method was called once with any Report object
            _mockRepo.Verify(repo => repo.AddAsync(It.IsAny<Report>()), Times.Once);
        }

        [Fact]
        public async Task GenerateReport_WithInvalidMessage_LogsWarning()
        {
            // Arrange
            var message = "invalid json";

            // Act
            await _service.GenerateReport(message);

            // Assert
            _mockLogger.Verify(
                x => x.Log(
                    LogLevel.Warning,
                    It.IsAny<EventId>(),
                    It.Is<It.IsAnyType>((o, t) => o.ToString().Contains("Failed to deserialize report request")),
                    It.IsAny<Exception>(),
                    It.IsAny<Func<It.IsAnyType, Exception, string>>()),
                Times.Once);

            _mockLogger.Verify(
                x => x.Log(
                    LogLevel.Error,
                    It.IsAny<EventId>(),
                    It.Is<It.IsAnyType>((o, t) => o.ToString().Contains("Error processing message")),
                    It.IsAny<Exception>(),
                    It.IsAny<Func<It.IsAnyType, Exception, string>>()),
                Times.Once);
        }
    }
}
