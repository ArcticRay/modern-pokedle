package observability

import "go.uber.org/zap"

type Logger struct {
	zap *zap.Logger
}

func NewLogger(env string) (*Logger, error) {
	var zapLogger *zap.Logger
	var err error

	switch env {
	case "production":
		zapLogger, err = zap.NewProduction()
	case "test":
		zapLogger = zap.NewNop()
	default:
		zapLogger, err = zap.NewDevelopment()
	}
	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(1))

	if err != nil {
		return nil, err
	}

	return &Logger{zap: zapLogger}, nil
}

func (l *Logger) Info(msg string, fields map[string]any) {
	l.zap.Info(msg, toZapFields(fields)...)
}

func (l *Logger) Error(msg string, fields map[string]any) {
	l.zap.Error(msg, toZapFields(fields)...)
}

func (l *Logger) Fatal(msg string, fields map[string]any) {
	l.zap.Fatal(msg, toZapFields(fields)...)
}

func (l *Logger) With(fields map[string]any) *Logger {
	return &Logger{zap: l.zap.With(toZapFields(fields)...)}
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}

func toZapFields(fields map[string]any) []zap.Field {
	result := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		result = append(result, zap.Any(k, v))
	}
	return result
}
