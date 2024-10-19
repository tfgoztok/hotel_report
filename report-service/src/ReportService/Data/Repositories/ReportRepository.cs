using MongoDB.Driver;
using ReportService.Interfaces;
using ReportService.Models;

namespace ReportService.Data.Repositories
{
    public class ReportRepository : IReportRepository
    {
        private readonly IMongoCollection<Report> _reports;

        // Constructor that initializes the MongoDB collection for reports
        public ReportRepository(IMongoDatabase database)
        {
            _reports = database.GetCollection<Report>("Reports");
        }

        // Asynchronously retrieves a report by its ID
        public async Task<Report> GetByIdAsync(string id)
        {
            return await _reports.Find(r => r.Id == id).FirstOrDefaultAsync();
        }

        // Asynchronously retrieves all reports
        public async Task<IEnumerable<Report>> GetAllAsync()
        {
            return await _reports.Find(_ => true).ToListAsync();
        }

        // Asynchronously adds a new report to the collection
        public async Task AddAsync(Report report)
        {
            await _reports.InsertOneAsync(report);
        }

        // Asynchronously updates an existing report in the collection
        public async Task UpdateAsync(Report report)
        {
            await _reports.ReplaceOneAsync(r => r.Id == report.Id, report);
        }
    }
}
