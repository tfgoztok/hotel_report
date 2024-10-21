using System.Collections.Generic;
using System.Threading.Tasks;
using System.Linq; // Bu satırı ekleyin
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Moq;
using ReportService.Controllers;
using ReportService.Interfaces;
using ReportService.Models;
using ReportService.Services;
using Xunit;

namespace ReportService.Tests
{
    public class ReportControllerTests
    {
        private readonly Mock<IReportRepository> _mockRepo;
        private readonly Mock<IRabbitMQService> _mockRabbitMQ;
        private readonly Mock<ILogger<ReportController>> _mockLogger;
        private readonly ReportController _controller;

        public ReportControllerTests()
        {
            _mockRepo = new Mock<IReportRepository>();
            _mockRabbitMQ = new Mock<IRabbitMQService>();
            _mockLogger = new Mock<ILogger<ReportController>>();
            _controller = new ReportController(_mockRepo.Object, _mockRabbitMQ.Object, _mockLogger.Object);
        }

        [Fact]
        public async Task GetAll_ReturnsAllReports()
        {
            // Arrange
            var expectedReports = new List<Report>
            {
                new Report { Id = "1", Location = "Istanbul" },
                new Report { Id = "2", Location = "Ankara" }
            };
            _mockRepo.Setup(repo => repo.GetAllAsync()).ReturnsAsync(expectedReports);

            // Act
            var result = await _controller.GetAll();

            // Assert
            var okResult = Assert.IsType<OkObjectResult>(result);
            var returnedReports = Assert.IsAssignableFrom<IEnumerable<Report>>(okResult.Value);
            Assert.Equal(expectedReports.Count, returnedReports.Count());
        }

        [Fact]
        public async Task Get_WithValidId_ReturnsReport()
        {
            // Arrange
            var expectedReport = new Report { Id = "1", Location = "Istanbul" };
            _mockRepo.Setup(repo => repo.GetByIdAsync("1")).ReturnsAsync(expectedReport);

            // Act
            var result = await _controller.Get("1");

            // Assert
            var okResult = Assert.IsType<OkObjectResult>(result);
            var returnedReport = Assert.IsType<Report>(okResult.Value);
            Assert.Equal(expectedReport.Id, returnedReport.Id);
        }

        [Fact]
        public void RequestReport_PublishesMessageAndReturnsAccepted()
        {
            // Arrange
            var request = new ReportService.Controllers.ReportRequest { Location = "Istanbul" };

            // Act
            var result = _controller.RequestReport(request);

            // Assert
            Assert.IsType<AcceptedResult>(result);
            _mockRabbitMQ.Verify(r => r.PublishMessage(It.IsAny<string>()), Times.Once);
        }
    }
}