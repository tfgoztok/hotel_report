using Microsoft.EntityFrameworkCore;
using ReportService.Data;
using ReportService.Interfaces;
using ReportService.Models;

namespace ReportService.Data.Repositories
{
    public class ReportRepository : IReportRepository
    {
        private readonly ApplicationDbContext _context;

        public ReportRepository(ApplicationDbContext context)
        {
            _context = context;
        }

        /// <summary>
        /// Retrieves a report by its unique identifier.
        /// </summary>
        public async Task<Report> GetByIdAsync(Guid id)
        {
            return await _context.Reports.FindAsync(id);
        }

        /// <summary>
        /// Retrieves all reports from the database.
        /// </summary>
        public async Task<IEnumerable<Report>> GetAllAsync()
        {
            return await _context.Reports.ToListAsync();
        }

        /// <summary>
        /// Adds a new report to the database.
        /// </summary>
        public async Task AddAsync(Report report)
        {
            await _context.Reports.AddAsync(report);
            await _context.SaveChangesAsync();
        }

        /// <summary>
        /// Updates an existing report in the database.
        /// </summary>
        public async Task UpdateAsync(Report report)
        {
            _context.Reports.Update(report);
            await _context.SaveChangesAsync();
        }
    }
}
