package main

import (
	"github.com/codex-team/hawk.workers.go/lib/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	rabbitmqUrl := getEnv("REGISTRY_URL", "amqp://127.0.0.1:5672")
	queue := "errors/default"

	cfg := zap.Config{
		Encoding:         "console",
		Development:      true,
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:          "message",
			LevelKey:            "level",
			TimeKey:             "name",
			NameKey:             "ts",
			CallerKey:           "caller",
			FunctionKey:         "func",
			StacktraceKey:       "stacktrace",
			SkipLineEnding:      false,
			LineEnding:          "\n",
			EncodeLevel:         zapcore.CapitalColorLevelEncoder,
			EncodeTime:          zapcore.ISO8601TimeEncoder,
			EncodeDuration:      zapcore.MillisDurationEncoder,
			EncodeCaller:        zapcore.FullCallerEncoder,
			EncodeName:          zapcore.FullNameEncoder,
			NewReflectedEncoder: nil,
			ConsoleSeparator:    "\t",
		},
	}
	logger, _ := cfg.Build()

	workerInstance := worker.New(rabbitmqUrl, queue, Handler, logger.Sugar())
	defer workerInstance.Stop() // TODO gracefully close connections on exit
	<-workerInstance.Run()
}
