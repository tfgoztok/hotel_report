using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System.Text;

namespace ReportService.Services
{
    public class RabbitMQService : IRabbitMQService
    {
        private readonly IConnection _connection; // Connection to the RabbitMQ server
        private readonly IModel _channel; // Channel for communication with RabbitMQ
        private readonly string _queueName; // Name of the queue to consume messages from

        // Constructor that initializes the RabbitMQ connection and declares the queue
        public RabbitMQService(IConfiguration configuration, IConnectionFactory connectionFactory)
        {
            _queueName = configuration["RabbitMQ:QueueName"] ?? throw new ArgumentNullException(nameof(configuration)); // Retrieve queue name from config
            _connection = connectionFactory.CreateConnection(); // Establish RabbitMQ connection
            _channel = _connection.CreateModel(); // Create RabbitMQ channel
            _channel.QueueDeclare(queue: _queueName, durable: false, exclusive: false, autoDelete: false, arguments: null); // Declare queue
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
            _channel.BasicConsume(queue: _queueName, autoAck: true, consumer: consumer); // Start consuming messages
        }

        // Dispose method to clean up resources
        public void Dispose()
        {
            _channel?.Dispose(); // Dispose of the channel if it exists
            _connection?.Dispose(); // Dispose of the connection if it exists
        }

        // Method to publish a message to the RabbitMQ queue
        public void PublishMessage(string message)
        {
            var body = Encoding.UTF8.GetBytes(message); // Convert the message to a byte array
            _channel.BasicPublish(exchange: "", routingKey: _queueName, basicProperties: null, body: body); // Publish the message to the queue
        }
    }
}
