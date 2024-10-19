using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System.Text;

namespace ReportService.Services
{
    public class RabbitMQService
    {
        private readonly IConnection _connection; // Connection to the RabbitMQ server
        private readonly IModel _channel; // Channel for communication with RabbitMQ
        private const string QueueName = "report_requests"; // Name of the queue to consume messages from

        // Constructor that initializes the RabbitMQ connection and declares the queue
        public RabbitMQService(IConfiguration configuration)
        {
            var factory = new ConnectionFactory() { HostName = configuration["RabbitMQ:HostName"] };
            _connection = factory.CreateConnection(); // Create a connection to RabbitMQ
            _channel = _connection.CreateModel(); // Create a channel
            // Declare the queue with specified parameters
            _channel.QueueDeclare(queue: QueueName, durable: false, exclusive: false, autoDelete: false, arguments: null);
        }

        // Method to start consuming messages from the queue
        public void StartConsuming(Action<string> processMessage)
        {
            var consumer = new EventingBasicConsumer(_channel); // Create a consumer
            consumer.Received += (model, ea) =>
            {
                var body = ea.Body.ToArray(); // Get the message body
                var message = Encoding.UTF8.GetString(body); // Convert the message body to a string
                processMessage(message); // Process the received message
            };
            _channel.BasicConsume(queue: QueueName, autoAck: true, consumer: consumer); // Start consuming messages
        }

        // Dispose method to clean up resources
        public void Dispose()
        {
            _channel?.Dispose(); // Dispose of the channel if it exists
            _connection?.Dispose(); // Dispose of the connection if it exists
        }
    }
}
