package main

import (
	"fmt"
	"github.com/codex-team/hawk.workers.go/lib/worker"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Supported log levels:
// "debug", "DEBUG"
// "info", "INFO", ""
// "warn", "WARN"
// "error", "ERROR"
// "dpanic", "DPANIC"
// "panic", "PANIC"
// "fatal", "FATAL"
func main() {
	if len(os.Args) > 1 {
		viper.SetConfigFile(os.Args[1])
	} else {
		viper.SetConfigName("default-processor")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	viper.SetDefault("registry_url", getEnv("REGISTRY_URL", "amqp://127.0.0.1:5672"))
	viper.SetDefault("log_level", getEnv("LOG_LEVEL", "info"))
	viper.SetDefault("queue", getEnv("QUEUE", "errors/default"))
	viper.SetDefault("target_queue", getEnv("TARGET_QUEUE", "grouper"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config was not found, using default values")
		} else {
			panic(fmt.Sprint("Failed to read the config: ", err.Error()))
		}
	}

	logLevel, _ := zapcore.ParseLevel(viper.GetString("log_level"))
	targetQueue = viper.GetString("target_queue")
	workerInstance := worker.New(
		viper.GetString("registry_url"),
		viper.GetString("queue"),
		Handler,
		worker.CreateDefaultLoggerWithLevel(logLevel))
	defer workerInstance.Stop() // TODO gracefully close connections on exit
	<-workerInstance.Run()
}
