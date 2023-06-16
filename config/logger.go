package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger(fg flags) *zap.Logger {
	logger, _ := zap.NewDevelopment()

	if *fg.logProd {
		atom := zap.NewAtomicLevel()
		if *fg.debug {
			atom.SetLevel(zap.DebugLevel)
		}

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stderr),
			atom,
		))
	}

	if *fg.logService != "" {
		logger = logger.With(zap.String("service", *fg.logService))
	}

	return logger
}
