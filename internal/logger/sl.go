package logger

import (
	"bufio"
	"cmd/obyavigo/main.go/internal/config"
	"cmd/obyavigo/main.go/internal/logger/handler"
	"log"
	"log/slog"
	"os"
)

func Init(l config.Logger, bSize int) (*slog.Logger, *handler.BHandler) {

	fe, err := os.OpenFile(l.ErrLogPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("error log file open error: ", err)
	}
	fi, err := os.OpenFile(l.InfoLogPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("info log file open error: ", err)
	}

	h := handler.New(bufio.NewWriter(fe), bufio.NewWriter(fi), bSize)

	return slog.New(h), h
}
