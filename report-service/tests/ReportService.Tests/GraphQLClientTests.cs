using System;
using System.Net;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Moq;
using Moq.Protected;
using ReportService.Services;
using Xunit;

namespace ReportService.Tests
{
    public class GraphQLClientTests
    {
        // Mocked HTTP message handler for testing
        private readonly Mock<HttpMessageHandler> _mockHttpMessageHandler;
        // HttpClient instance for making requests
        private readonly HttpClient _httpClient;
        // Endpoint for the GraphQL service
        private readonly string _endpoint = "http://test-graphql-endpoint.com/graphql";

        public GraphQLClientTests()
        {
            // Initialize the mock HTTP message handler and HttpClient
            _mockHttpMessageHandler = new Mock<HttpMessageHandler>();
            _httpClient = new HttpClient(_mockHttpMessageHandler.Object);
        }

        [Fact]
        public async Task SendQueryAsync_SuccessfulRequest_ReturnsResponse()
        {
            // Arrange: Set up expected response and mock behavior
            var expectedResponse = "{\"data\":{\"test\":\"result\"}}";
            SetupMockHttpMessageHandler(HttpStatusCode.OK, expectedResponse);

            // Create an instance of the GraphQL client
            var graphQLClient = new GraphQLClient(_httpClient, _endpoint);

            // Act: Send a query to the GraphQL client
            var result = await graphQLClient.SendQueryAsync("query { test }", null);

            // Assert: Verify that the result matches the expected response
            Assert.Equal(expectedResponse, result);
        }

        [Fact]
        public async Task SendQueryAsync_ErrorResponse_ThrowsException()
        {
            // Arrange: Set up mock to return an error response
            SetupMockHttpMessageHandler(HttpStatusCode.InternalServerError, "Server Error");

            // Create an instance of the GraphQL client
            var graphQLClient = new GraphQLClient(_httpClient, _endpoint);

            // Act & Assert: Verify that an exception is thrown for error responses
            await Assert.ThrowsAsync<HttpRequestException>(() => 
                graphQLClient.SendQueryAsync("query { test }", null));
        }

        // Helper method to set up the mock HTTP message handler
        private void SetupMockHttpMessageHandler(HttpStatusCode statusCode, string content)
        {
            _mockHttpMessageHandler.Protected()
                .Setup<Task<HttpResponseMessage>>(
                    "SendAsync",
                    ItExpr.IsAny<HttpRequestMessage>(),
                    ItExpr.IsAny<CancellationToken>()
                )
                .ReturnsAsync(new HttpResponseMessage
                {
                    StatusCode = statusCode,
                    Content = new StringContent(content)
                });
        }
    }
}
