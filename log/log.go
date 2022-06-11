package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
)

type InfoLogger interface {
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Enabled() bool
}

type Logger interface {
	InfoLogger

	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	Panic(msg string, fields ...Field)
	Panicf(format string, v ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	V(level int) InfoLogger

	Write(p []byte) (n int, err error)

	WithValues(keysAndValues ...interface{}) Logger

	WithName(name string) Logger

	WithContext(ctx context.Context) context.Context

	Flush()
}

// noopInfoLogger TODO
type noopInfoLogger struct{}

var _ InfoLogger = &noopInfoLogger{}

func (l *noopInfoLogger) Enabled() bool                    { return false }
func (l *noopInfoLogger) Info(_ string, _ ...Field)        {}
func (l *noopInfoLogger) Infof(_ string, _ ...interface{}) {}
func (l *noopInfoLogger) Infow(_ string, _ ...interface{}) {}

var disabledInfoLogger = &noopInfoLogger{}

// infoLogger TODO
type infoLogger struct {
	level zapcore.Level
	log   *zap.Logger
}

var _ InfoLogger = &infoLogger{}

func Info(msg string, field ...Field) {
	_logger.Info(msg, field...)
}

func Infof(format string, v ...interface{}) {
	_logger.Infof(format, v...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	_logger.Infow(msg, keysAndValues...)
}

func (l *infoLogger) Info(msg string, fields ...Field) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(fields...)
	}
}

func (l *infoLogger) Infof(format string, v ...interface{}) {
	if checkedEntry := l.log.Check(l.level, fmt.Sprintf(format, v...)); checkedEntry != nil {
		checkedEntry.Write()
	}
}

func (l *infoLogger) Infow(msg string, keysAndValues ...interface{}) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(l.log, keysAndValues)...)
	}
}

func (l *infoLogger) Enabled() bool {
	return true
}

func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		return additional
	}

	fileds := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); i += 2 {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}

		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}

		var (
			key, val  = args[i], args[i+1]
			keyString string
			isKey     bool
		)

		if keyString, isKey = key.(string); !isKey {
			l.DPanic("non-string key argument passed to logging, ignoring all later arguments", zap.Any("invalid key", key))
			break
		}

		fileds = append(fileds, zap.Any(keyString, val))
	}

	return append(fileds, additional...)
}

// TODO
type zapLogger struct {
	// NB: this looks very similar to zap.SugaredLogger, but
	// deals with our desire to have multiple verbosity levels.
	zapLogger        *zap.Logger
	zapSugaredLogger *zap.SugaredLogger
	infoLogger
}

var _ Logger = &zapLogger{}

// 记录 Debug 级别的日志.

func Debug(msg string, fields ...Field) {
	_logger.zapLogger.Debug(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	_logger.zapSugaredLogger.Debugf(format, v...)
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.zapSugaredLogger.Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	_logger.zapSugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Debugw(msg, keysAndValues...)
}

// 记录 Warn 级别的日志.

func Warn(msg string, fields ...Field) {
	_logger.zapLogger.Warn(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	_logger.zapSugaredLogger.Warnf(format, v...)
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.zapSugaredLogger.Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	_logger.zapSugaredLogger.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Warnw(msg, keysAndValues...)
}

// 记录 Error 级别的日志.

func Error(msg string, fields ...Field) {
	_logger.zapLogger.Error(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	_logger.zapSugaredLogger.Errorf(format, v...)
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.zapSugaredLogger.Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	_logger.zapSugaredLogger.Errorf(msg, keysAndValues...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Errorf(msg, keysAndValues...)
}

// 记录 Panic 级别的日志.

func Panic(msg string, fields ...Field) {
	_logger.zapLogger.Panic(msg, fields...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	_logger.zapSugaredLogger.Panicf(format, v...)
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.zapSugaredLogger.Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	_logger.zapSugaredLogger.Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Panicw(msg, keysAndValues...)
}

// 记录 Fatal 级别的日志.

func Fatal(msg string, fields ...Field) {
	_logger.zapLogger.Fatal(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	_logger.zapSugaredLogger.Fatalf(format, v...)
}

func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.zapSugaredLogger.Fatalf(format, v...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	_logger.zapSugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Fatalw(msg, keysAndValues...)
}

func V(level int) InfoLogger {
	return _logger.V(level)
}

func (l *zapLogger) V(level int) InfoLogger {
	lvl := zapcore.Level(-1 * level)

	if ok := l.zapLogger.Core().Enabled(lvl); ok {
		return &infoLogger{
			log:   l.zapLogger,
			level: lvl,
		}
	}

	return disabledInfoLogger
}

func (l *zapLogger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))

	return len(p), nil
}

func WithValues(keysAndValues ...interface{}) Logger {
	return _logger.WithValues(keysAndValues...)
}

func (l *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	newL := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)

	return NewLogger(newL)
}

func WithName(name string) Logger {
	return _logger.WithName(name)
}

func (l *zapLogger) WithName(name string) Logger {
	newL := l.zapLogger.Named(name)

	return NewLogger(newL)
}

func L(ctx context.Context) *zapLogger {
	return _logger.L(ctx)
}

func (l *zapLogger) L(ctx context.Context) *zapLogger {
	lg := l.clone()

	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyRequestID, requestID))
	}
	if username := ctx.Value(KeyUsername); username != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyUsername, username))
	}
	if watcherName := ctx.Value(KeyWatcherName); watcherName != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyWatcherName, watcherName))
	}

	return lg
}

func Flush() {
	_logger.Flush()
}

func (l *zapLogger) Flush() {
	_ = l.zapLogger.Sync()
}

func (l *zapLogger) clone() *zapLogger {
	copy := *l

	return &copy
}

var (
	_logger = New(NewOptions())
	_mu     sync.Mutex
)

func InitGlobalLogger(opt *Options) {
	_mu.Lock()
	defer _mu.Unlock()

	_logger = New(opt)
}

func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger:        l,
		zapSugaredLogger: l.Sugar(),
		infoLogger: infoLogger{
			log:   l,
			level: zapcore.InfoLevel,
		},
	}
}

func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var (
		level          Level
		err            error
		encodeLevel    = zapcore.CapitalLevelEncoder
		encoderConfig  zapcore.EncoderConfig
		SamplingConfig *zap.SamplingConfig
		loggerConfig   *zap.Config
	)

	if err = level.UnmarshalText([]byte(opts.Level)); err != nil {
		level = zap.InfoLevel
	}

	if opts.Format == ConsoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig = zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "level",
		TimeKey:             "timestamp",
		NameKey:             "logger",
		CallerKey:           "caller",
		FunctionKey:         "",
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         encodeLevel,
		EncodeTime:          timeEncoder,
		EncodeDuration:      milliSecondsDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		EncodeName:          zapcore.FullNameEncoder,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    " ",
	}

	SamplingConfig = &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	loggerConfig = &zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       opts.Development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Sampling:          SamplingConfig,
		Encoding:          opts.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  opts.ErrorOutputPaths,
		InitialFields:     nil,
	}

	// 配置完成, 根据配置生成 zap.logger
	var (
		l      *zap.Logger
		logger *zapLogger
	)

	if l, err = loggerConfig.Build(zap.AddStacktrace(zap.PanicLevel), zap.AddCallerSkip(1)); err != nil {
		panic(err)
	}

	l = l.Named(opts.Name)
	logger = &zapLogger{
		zapLogger:        l,
		zapSugaredLogger: l.Sugar(),
		infoLogger: infoLogger{
			level: zap.InfoLevel,
			log:   l,
		},
	}

	zap.RedirectStdLog(l)

	return logger
}

// ZapLogger 返回全局 logger 的 zap logger.
func ZapLogger() *zap.Logger {
	return _logger.zapLogger
}

// SugaredLogger 返回全局 logger 的 sugared logger.
func SugaredLogger() *zap.SugaredLogger {
	return _logger.zapSugaredLogger
}

// StdErrLogger 返回全局对象 Error 级别的 logger.
func StdErrLogger() *log.Logger {
	if _logger == nil {
		return nil
	}

	if l, err := zap.NewStdLogAt(_logger.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}

	return nil
}

// StdInfoLogger 返回全局对象 Info 级别的 logger.
func StdInfoLogger() *log.Logger {
	if _logger == nil {
		return nil
	}

	if l, err := zap.NewStdLogAt(_logger.zapLogger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}
