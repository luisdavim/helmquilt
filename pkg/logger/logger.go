package logger

import (
	"context"
	"fmt"
	"log"
	"os"
)

type ctxKey string

const (
	loggerKey ctxKey = "loggerKey"
)

func New(name string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", name), log.LstdFlags)
}

func NewContext(ctx context.Context, name string) context.Context {
	return IntoContext(ctx, New(name))
}

func IntoContext(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerKey).(*log.Logger); ok {
		return logger
	}

	return New("helmquilt")
}
