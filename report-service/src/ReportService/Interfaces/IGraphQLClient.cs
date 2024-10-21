// This interface defines a contract for sending GraphQL queries.
public interface IGraphQLClient
{
    Task<string> SendQueryAsync(string query, object? variables = null);
}
