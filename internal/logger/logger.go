package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var L = log.Logger

func InitLogger(dataDir string) {
	zerolog.TimeFieldFormat = "2006/01/02 15:04:05"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05"}

	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(dataDir, "logs.json"),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     28,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(consoleWriter, fileWriter)
	L = zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Logger()
}
