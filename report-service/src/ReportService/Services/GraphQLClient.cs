using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;

namespace ReportService.Services
{
    public class GraphQLClient
    {
        private readonly HttpClient _httpClient; // HttpClient for making requests
        private readonly string _endpoint; // Endpoint for the GraphQL service

        // Constructor to initialize HttpClient and endpoint
        public GraphQLClient(HttpClient httpClient, string endpoint)
        {
            _httpClient = httpClient;
            _endpoint = endpoint;
        }

        // Method to send a GraphQL query asynchronously
        public async Task<string> SendQueryAsync(string query, object variables = null)
        {
            // Create a request object containing the query and variables
            var request = new
            {
                query,
                variables
            };

            // Serialize the request object to JSON
            var json = JsonSerializer.Serialize(request);
            // Create content for the HTTP request
            var content = new StringContent(json, Encoding.UTF8, "application/json");

            // Send the POST request to the GraphQL endpoint
            var response = await _httpClient.PostAsync(_endpoint, content);
            // Ensure the response indicates success
            response.EnsureSuccessStatusCode();

            // Read and return the response content as a string
            return await response.Content.ReadAsStringAsync();
        }
    }
}
