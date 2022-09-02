package worker

// MainConfig Supported log levels:
// "debug", "DEBUG"
// "info", "INFO", ""
// "warn", "WARN"
// "error", "ERROR"
// "dpanic", "DPANIC"
// "panic", "PANIC"
// "fatal", "FATAL"
type MainConfig struct {
	RegistryUrl string `mapstructure:"registry_url"` // URL to the RabbitMQ, for example "amqp://127.0.0.1:5672"
	LogLevel    string `mapstructure:"log_level"`    // Log level applied
}
