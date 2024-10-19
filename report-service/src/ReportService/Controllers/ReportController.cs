using Microsoft.AspNetCore.Mvc;
using ReportService.Interfaces;
using ReportService.Models;

namespace ReportService.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class ReportController : ControllerBase
    {
        // Dependency injection for the report repository and logger
        private readonly IReportRepository _reportRepository;
        private readonly ILogger<ReportController> _logger;

        public ReportController(IReportRepository reportRepository, ILogger<ReportController> logger)
        {
            _reportRepository = reportRepository;
            _logger = logger;
        }

        // GET: /report - Retrieves all reports
        [HttpGet]
        public async Task<IActionResult> GetAll()
        {
            var reports = await _reportRepository.GetAllAsync();
            return Ok(reports);
        }

        // GET: /report/{id} - Retrieves a report by its ID
        [HttpGet("{id}")]
        public async Task<IActionResult> Get(Guid id)
        {
            var report = await _reportRepository.GetByIdAsync(id);
            if (report == null)
            {
                return NotFound();
            }
            return Ok(report);
        }

        // POST: /report - Requests a new report
        [HttpPost]
        public IActionResult RequestReport([FromBody] ReportRequest request)
        {
            // TODO: Send message to RabbitMQ
            _logger.LogInformation($"Report requested for location: {request.Location}");
            return Accepted();
        }
    }

    // Model for report request
    public class ReportRequest
    {
        public string Location { get; set; }
    }
}
