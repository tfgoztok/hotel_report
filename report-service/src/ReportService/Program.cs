using Microsoft.EntityFrameworkCore;
using ReportService.Data;
using ReportService.Services;
using ReportService.Interfaces;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers(); // Adds MVC controllers to the service container
builder.Services.AddEndpointsApiExplorer(); // Adds support for API endpoint exploration
builder.Services.AddSwaggerGen(); // Adds Swagger generation for API documentation

// Add DbContext
builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection"))); // Configures the DbContext to use SQL Server with the connection string

// Add services
builder.Services.AddSingleton<RabbitMQService>(); // Registers RabbitMQService as a singleton
builder.Services.AddScoped<IReportRepository, ReportRepository>(); // Registers the report repository with scoped lifetime
builder.Services.AddScoped<IReportGenerationService, ReportGenerationService>(); // Registers the report generation service with scoped lifetime

var app = builder.Build(); // Builds the application

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger(); // Enables Swagger in development mode
    app.UseSwaggerUI(); // Enables the Swagger UI in development mode
}

app.UseHttpsRedirection(); // Redirects HTTP requests to HTTPS
app.UseAuthorization(); // Enables authorization middleware
app.MapControllers(); // Maps attribute-routed controllers

// Start consuming RabbitMQ messages
var rabbitMQService = app.Services.GetRequiredService<RabbitMQService>(); // Retrieves the RabbitMQ service from the service provider
rabbitMQService.StartConsuming(message => 
{
    using var scope = app.Services.CreateScope(); // Creates a new scope for dependency injection
    var reportService = scope.ServiceProvider.GetRequiredService<IReportGenerationService>(); // Retrieves the report generation service
    reportService.GenerateReport(message); // Generates a report based on the received message
});

app.Run(); // Runs the application
