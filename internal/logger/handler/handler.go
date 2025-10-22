package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// после Close() при попытке использования будет ошибка
type BHandler struct {
	eCh   chan slog.Record
	iCh   chan slog.Record
	eFile *bufio.Writer
	iFile *bufio.Writer
	wg    sync.WaitGroup
}

func New(errF *bufio.Writer, infoF *bufio.Writer, bSize int) *BHandler {
	h := &BHandler{
		eCh:   make(chan slog.Record, bSize),
		iCh:   make(chan slog.Record, bSize),
		eFile: errF,
		iFile: infoF,
		wg:    sync.WaitGroup{},
	}

	h.wg.Add(2)
	go h.fileWriter(h.eCh, h.eFile)
	go h.fileWriter(h.iCh, h.iFile)

	return h
}

func (h *BHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return true
}

func (h *BHandler) Handle(ctx context.Context, r slog.Record) error {

	switch r.Level {
	case slog.LevelError:
		select {
		case h.eCh <- r:
			return nil
		default:
			return fmt.Errorf("error chan is full")
		}
	default:
		select {
		case h.iCh <- r:
			return nil
		default:
			return fmt.Errorf("info chan is full")
		}
	}
}

func (h *BHandler) fileWriter(ch chan slog.Record, f *bufio.Writer) {
	defer h.wg.Done()
	for {
		r, ok := <-ch
		if !ok {
			return
		}

		s, err := h.formatRecord(r)
		if err != nil {
			slog.Error("formatting record error", slog.Any("error", err))
			continue
		}

		if _, err := f.WriteString(s + "\n"); err != nil {
			slog.Error("file write error", slog.Any("error", err))
		}
		f.Flush()
	}
}

func (h *BHandler) formatRecord(r slog.Record) (string, error) {
	m := map[string]interface{}{
		"time":    r.Time.Format(time.RFC3339),
		"level":   r.Level.String(),
		"source":  r.Source().File,
		"message": r.Message,
	}

	r.Attrs(func(a slog.Attr) bool {
		m[a.Key] = a.Value.Any()
		return true
	})

	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// мб доделать бы
func (h *BHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// и это тоже
func (h *BHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *BHandler) Close() {
	close(h.eCh)
	close(h.iCh)
	h.wg.Wait()
}
