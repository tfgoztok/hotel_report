public interface IRabbitMQService
{
    void StartConsuming(Action<string> processMessage);
    void Dispose();
    void PublishMessage(string message);
}