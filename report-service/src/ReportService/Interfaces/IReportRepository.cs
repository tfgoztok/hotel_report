using ReportService.Models;

namespace ReportService.Interfaces
{
    public interface IReportRepository
    {
        // This interface defines the contract for report repository operations
        Task<Report> GetByIdAsync(string  id);
        Task<IEnumerable<Report>> GetAllAsync();
        Task AddAsync(Report report);
        Task UpdateAsync(Report report);
    }
}
