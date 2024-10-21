using System;
using Microsoft.Extensions.Configuration;
using Moq;
using RabbitMQ.Client;
using ReportService.Services;
using Xunit;
using System.Collections.Generic;

namespace ReportService.Tests
{
    public class RabbitMQServiceTests
    {
        // Mock objects for dependencies
        private readonly Mock<IConfiguration> _mockConfiguration;
        private readonly Mock<IConnectionFactory> _mockConnectionFactory;
        private readonly Mock<IConnection> _mockConnection;
        private readonly Mock<IModel> _mockChannel;

        public RabbitMQServiceTests()
        {
            // Initialize mock objects
            _mockConfiguration = new Mock<IConfiguration>();
            _mockConnectionFactory = new Mock<IConnectionFactory>();
            _mockConnection = new Mock<IConnection>();
            _mockChannel = new Mock<IModel>();

            // Setup mock behavior for configuration and connection factory
            _mockConfiguration.Setup(c => c["RabbitMQ:HostName"]).Returns("localhost");
            _mockConnectionFactory.Setup(cf => cf.CreateConnection()).Returns(_mockConnection.Object);
            _mockConnection.Setup(c => c.CreateModel()).Returns(_mockChannel.Object);
        }

        [Fact]
        public void Constructor_InitializesRabbitMQConnection()
        {
            // Arrange: Setup queue name for testing
            _mockConfiguration.Setup(c => c["RabbitMQ:QueueName"]).Returns("test_queue");

            // Act: Create instance of RabbitMQService
            var rabbitMQService = new RabbitMQService(_mockConfiguration.Object, _mockConnectionFactory.Object);

            // Assert: Verify that connection and channel are created
            _mockConnectionFactory.Verify(cf => cf.CreateConnection(), Times.Once);
            _mockConnection.Verify(c => c.CreateModel(), Times.Once);
            _mockChannel.Verify(ch => ch.QueueDeclare(
                It.IsAny<string>(),
                It.IsAny<bool>(),
                It.IsAny<bool>(),
                It.IsAny<bool>(),
                It.IsAny<IDictionary<string, object>>()
            ), Times.Once);
        }

        [Fact]
        public void PublishMessage_SendsMessageToQueue()
        {
            // Arrange: Setup queue name and message for testing
            _mockConfiguration.Setup(c => c["RabbitMQ:QueueName"]).Returns("test_queue");
            var rabbitMQService = new RabbitMQService(_mockConfiguration.Object, _mockConnectionFactory.Object);
            var message = "Test message";

            // Act: Publish message to the queue
            rabbitMQService.PublishMessage(message);

            // Assert: Verify that the message was published
            _mockChannel.Verify(ch => ch.BasicPublish(
                It.IsAny<string>(),
                It.IsAny<string>(),
                It.IsAny<bool>(),
                It.IsAny<IBasicProperties>(),
                It.Is<ReadOnlyMemory<byte>>(b => System.Text.Encoding.UTF8.GetString(b.ToArray()) == message)
            ), Times.Once);
        }

        [Fact]
        public void StartConsuming_SetupConsumerCorrectly()
        {
            // Arrange: Setup queue name for testing
            _mockConfiguration.Setup(c => c["RabbitMQ:QueueName"]).Returns("test_queue");
            var rabbitMQService = new RabbitMQService(_mockConfiguration.Object, _mockConnectionFactory.Object);

            // Act: Start consuming messages from the queue
            rabbitMQService.StartConsuming(_ => { });

            // Assert: Verify that the consumer is set up correctly
            _mockChannel.Verify(ch => ch.BasicConsume(
                It.IsAny<string>(),
                It.IsAny<bool>(),
                It.IsAny<string>(),
                It.IsAny<bool>(),
                It.IsAny<bool>(),
                It.IsAny<IDictionary<string, object>>(),
                It.IsAny<IBasicConsumer>()
            ), Times.Once);
        }
    }
}
