using Microsoft.EntityFrameworkCore;
using ReportService.Models;

namespace ReportService.Data
{
    public class ApplicationDbContext : DbContext
    {
        // This class represents the database context for the application, managing the Reports DbSet.
        public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
            : base(options)
        {
        }

        public DbSet<Report> Reports { get; set; }
    }
}
