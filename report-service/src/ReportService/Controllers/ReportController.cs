using Microsoft.AspNetCore.Mvc;
using ReportService.Interfaces;
using ReportService.Models;
using ReportService.Services;
using System.Text.Json;

namespace ReportService.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class ReportController : ControllerBase
    {
        // Dependency injection for the report repository, RabbitMQ service, and logger
        private readonly IReportRepository _reportRepository;
        private readonly RabbitMQService _rabbitMQService;
        private readonly ILogger<ReportController> _logger;

        public ReportController(IReportRepository reportRepository, RabbitMQService rabbitMQService, ILogger<ReportController> logger)
        {
            _reportRepository = reportRepository;
            _rabbitMQService = rabbitMQService;
            _logger = logger;
        }

        // GET: Retrieve all reports
        [HttpGet]
        public async Task<IActionResult> GetAll()
        {
            var reports = await _reportRepository.GetAllAsync();
            return Ok(reports);
        }

        // GET: Retrieve a specific report by ID
        [HttpGet("{id}")]
        public async Task<IActionResult> Get(string id)
        {
            var report = await _reportRepository.GetByIdAsync(id);
            if (report == null)
            {
                return NotFound();
            }
            return Ok(report);
        }

        // POST: Request a new report
        [HttpPost]
        public IActionResult RequestReport([FromBody] ReportRequest request)
        {
            var message = JsonSerializer.Serialize(request);
            _rabbitMQService.PublishMessage(message);
            _logger.LogInformation($"Report requested for location: {request.Location}");
            return Accepted();
        }
    }

    // Model for report request containing location information
    public class ReportRequest
    {
        public string Location { get; set; } = null!;
    }
}
