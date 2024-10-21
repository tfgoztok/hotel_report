using MongoDB.Driver;
using ReportService.Data.Repositories;
using ReportService.Interfaces;
using ReportService.Services;
using Serilog;
using Serilog.Sinks.Elasticsearch;
using Microsoft.Extensions.Logging;

var builder = WebApplication.CreateBuilder(args);

// Configure Serilog for logging
Log.Logger = new LoggerConfiguration()
    .Enrich.FromLogContext() // Enrich logs with contextual information
    .WriteTo.Elasticsearch(new ElasticsearchSinkOptions(new Uri(builder.Configuration["ELASTICSEARCH_URL"] ?? "http://elasticsearch:9200"))
    {
        AutoRegisterTemplate = true, // Automatically register the template in Elasticsearch
        IndexFormat = $"report-service-logs-{DateTime.UtcNow:yyyy-MM}" // Format for log index
    })
    .WriteTo.Console() // Also write logs to the console
    .CreateLogger();

builder.Host.UseSerilog(); // Use Serilog for logging in the application

// Add services to the container.
builder.Services.AddControllers(); // Add MVC controllers
builder.Services.AddEndpointsApiExplorer(); // Add support for API endpoint exploration
builder.Services.AddSwaggerGen(); // Add Swagger for API documentation

// Configure MongoDB
var mongoConnectionString = builder.Configuration["MongoDbSettings:ConnectionString"]; // Get MongoDB connection string from configuration
var mongoDatabaseName = builder.Configuration["MongoDbSettings:DatabaseName"]; // Get MongoDB database name from configuration
var mongoClient = new MongoClient(mongoConnectionString); // Create a MongoDB client
var mongoDatabase = mongoClient.GetDatabase(mongoDatabaseName); // Get the specified database
builder.Services.AddSingleton<IMongoDatabase>(mongoDatabase); // Register MongoDB database as a singleton service

// Add services
builder.Services.AddSingleton<RabbitMQService>(); // Register RabbitMQ service as a singleton
builder.Services.AddScoped<IReportRepository, ReportRepository>(); // Register report repository with scoped lifetime
builder.Services.AddScoped<IReportGenerationService, ReportGenerationService>(); // Register report generation service with scoped lifetime

// Add HTTP client services for making HTTP requests
builder.Services.AddHttpClient();

builder.Services.AddScoped<IGraphQLClient>(sp =>
{
    var httpClientFactory = sp.GetRequiredService<IHttpClientFactory>();
    var httpClient = httpClientFactory.CreateClient();

    var configuration = sp.GetRequiredService<IConfiguration>();
    var endpoint = configuration["GraphQL:Endpoint"] 
        ?? throw new InvalidOperationException("GraphQL:Endpoint is not configured");

    return new GraphQLClient(httpClient, endpoint);
});



var app = builder.Build(); // Build the application

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger(); // Enable Swagger in development mode
    app.UseSwaggerUI(); // Enable Swagger UI in development mode
}

app.UseHttpsRedirection(); // Redirect HTTP requests to HTTPS
app.UseAuthorization(); // Enable authorization middleware
app.MapControllers(); // Map controller routes

// Start consuming RabbitMQ messages
var rabbitMQService = app.Services.GetRequiredService<RabbitMQService>(); // Get the RabbitMQ service
rabbitMQService.StartConsuming(message =>
{
    try
    {
        using var scope = app.Services.CreateScope(); // Create a scope for dependency injection
        var logger = scope.ServiceProvider.GetRequiredService<ILogger<RabbitMQService>>(); // Get the logger instance
        logger.LogInformation("Consuming message and creating scope."); // Use the logger instance
        var reportService = scope.ServiceProvider.GetRequiredService<IReportGenerationService>(); // Get the report generation service
        logger.LogInformation("Report service resolved, calling GenerateReport."); // Use the logger instance
        reportService.GenerateReport(message); // Generate a report based on the received message
        logger.LogInformation("GenerateReport called."); // Use the logger instance
    }
    catch (Exception ex)
    {
        Console.WriteLine("An error occurred while consuming the message: " + ex.Message);
    }

});

app.Run(); // Run the application
