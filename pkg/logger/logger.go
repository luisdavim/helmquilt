package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

type ctxKey string

const (
	loggerKey ctxKey = "loggerKey"
)

// New creates a new logger with the given name and output
func New(name string, out io.Writer) *log.Logger {
	return log.New(out, fmt.Sprintf("[%s] ", name), log.LstdFlags)
}

// NewContext returns a derived contex with a logger as a child of the given context
func NewContext(ctx context.Context, name string, out io.Writer) context.Context {
	return IntoContext(ctx, New(name, out))
}

// IntoContext returns a derived context that points to the given parent contex and the provided logger
func IntoContext(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves a logger from the provided context
// if a logger is not found in the context a new one is created so the return logger is always functional
func FromContext(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerKey).(*log.Logger); ok {
		return logger
	}

	return New("helmquilt", os.Stderr)
}
