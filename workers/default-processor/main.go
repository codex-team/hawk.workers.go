package main

import (
	"fmt"
	"github.com/codex-team/hawk.workers.go/lib/worker"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"os"
)

const targetQueue string = "grouper"
const sourceQueue string = "errors/default"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

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
	var config worker.CommonConfig
	if err := worker.ReadConfig(&config); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config was not found, using default values")
		} else {
			panic(err.Error())
		}
	}

	logLevel, _ := zapcore.ParseLevel(config.LogLevel)
	workerInstance := worker.New(
		config.RegistryUrl,
		sourceQueue,
		Handler,
		worker.CreateDefaultLoggerWithLevel(logLevel))
	defer workerInstance.Stop() // TODO gracefully close connections on exit
	<-workerInstance.Run()
}
