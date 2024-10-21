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
            // Initialize mock objects and set up their behavior
            _mockConfiguration = new Mock<IConfiguration>();
            _mockConnectionFactory = new Mock<IConnectionFactory>();
            _mockConnection = new Mock<IConnection>();
            _mockChannel = new Mock<IModel>();

            // Configure mock to return a test queue name
            _mockConfiguration.Setup(c => c["RabbitMQ:QueueName"]).Returns("test_queue");
            // Set up connection factory to return a mock connection
            _mockConnectionFactory.Setup(cf => cf.CreateConnection()).Returns(_mockConnection.Object);
            // Set up connection to return a mock channel
            _mockConnection.Setup(c => c.CreateModel()).Returns(_mockChannel.Object);
        }

        [Fact]
        public void Constructor_InitializesRabbitMQConnection()
        {
            // Arrange & Act
            var rabbitMQService = new RabbitMQService(_mockConfiguration.Object, _mockConnectionFactory.Object);

            // Assert that the connection and channel are created once
            _mockConnectionFactory.Verify(cf => cf.CreateConnection(), Times.Once);
            _mockConnection.Verify(c => c.CreateModel(), Times.Once);
            // Verify that the queue is declared once
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
            // Arrange
            var rabbitMQService = new RabbitMQService(_mockConfiguration.Object, _mockConnectionFactory.Object);
            var message = "Test message";

            // Act
            rabbitMQService.PublishMessage(message);

            // Assert that the message is published to the queue
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
            // Arrange
            var rabbitMQService = new RabbitMQService(_mockConfiguration.Object, _mockConnectionFactory.Object);

            // Act
            rabbitMQService.StartConsuming(_ => { });

            // Assert that the consumer is set up correctly
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
