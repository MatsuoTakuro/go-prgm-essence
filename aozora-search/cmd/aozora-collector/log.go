package main

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func setupLogger(logFile string) (*zap.Logger, func() error, error) {
	if logger != nil {
		// in case of testing, already set it up
		return logger, nil, nil
	}

	// Open the log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// Create a zapcore.WriteSyncer from the opened file
	writeSyncer := zapcore.AddSync(file)

	// Create a new production encoder config
	encoderConfig := zap.NewProductionEncoderConfig()

	// Create a new zapcore.Core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zap.InfoLevel,
	)

	// Create a new zap.Logger
	logger := zap.New(core)

	return logger, func() error {
		if err := logger.Sync(); err != nil {
			log.Println(err)
		}
		return file.Close()
	}, nil
}
