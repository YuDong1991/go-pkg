package log

import "context"

type key int

const (
	logContextKey key = iota
)

func WithContext(ctx context.Context) context.Context {
	return _logger.WithContext(ctx)
}

func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		if logger := ctx.Value(logContextKey); logger != nil {
			return logger.(*zapLogger)
		}
	}

	return WithName("Unknown-Context")
}
