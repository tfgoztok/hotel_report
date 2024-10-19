package config

import (
	"github.com/spf13/viper"
)

// Config struct holds the configuration values for the application.
type Config struct {
	Port        string `mapstructure:"PORT"`         // Port for the application
	DatabaseURL string `mapstructure:"DATABASE_URL"` // Database connection URL
	LogLevel    string `mapstructure:"LOG_LEVEL"`    // Logging level for the application
	RabbitMQURL string `mapstructure:"RABBITMQ_URL"` // RabbitMQ connection URL
}

// Load function initializes the configuration by reading environment variables and setting defaults.
func Load() (*Config, error) {
	viper.AutomaticEnv() // Automatically read environment variables

	// Set default values for configuration
	viper.SetDefault("PORT", "8080")                                                                    // Default port
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/hoteldb?sslmode=disable") // Default database URL
	viper.SetDefault("LOG_LEVEL", "info")                                                               // Default log level
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")                              // Default Rabbitmq URL

	var config Config                                // Create an instance of Config to hold the values
	if err := viper.Unmarshal(&config); err != nil { // Unmarshal environment variables into the config struct
		return nil, err // Return nil and the error if unmarshalling fails
	}

	return &config, nil // Return the populated config struct and nil error
}
