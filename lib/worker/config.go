package worker

import "github.com/spf13/viper"

// CommonConfig Supported log levels:
// "debug", "DEBUG"
// "info", "INFO", ""
// "warn", "WARN"
// "error", "ERROR"
// "dpanic", "DPANIC"
// "panic", "PANIC"
// "fatal", "FATAL"
type CommonConfig struct {
	RegistryUrl string // URL to the RabbitMQ, for example "amqp://127.0.0.1:5672"
	LogLevel    string // Log level applied
}

// ReadConfig reads config to the receiver, needs viper to be configured with paths and defaults
// see viper docs
func ReadConfig(receiver interface{}) error {
	var result error
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			result = err
		} else {
			return err
		}
	}
	if err := viper.Unmarshal(&receiver); err != nil {
		return err
	}
	return result
}
