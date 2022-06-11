package log

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	ConsoleFormat = "console"
	JsonFormat    = "json"
)

type Options struct {
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level             string   `json:"level"              mapstructure:"level"`
	Format            string   `json:"format"             mapstructure:"format"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`
	Development       bool     `json:"development"        mapstructure:"development"`
	Name              string   `json:"name"               mapstructure:"name"`
}

func NewOptions() *Options {
	return &Options{
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		Level:             zapcore.InfoLevel.String(),
		Format:            ConsoleFormat,
		DisableCaller:     false,
		DisableStacktrace: false,
		EnableColor:       false,
		Development:       false,
		Name:              "",
	}
}

// Validate 用于对 Options 的各字段进行合法性验证.
func (o *Options) Validate() []error {
	var errs []error

	// 验证 Options.Level 字段是否合法.
	var level Level
	if err := level.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	// 验证 Options.Format 字段是否合法.
	format := strings.ToLower(o.Format)
	switch format {
	case ConsoleFormat, JsonFormat:
		break
	default:
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}

// Json 用于将 Options 序列化为 Json 字符串.
func (o *Options) Json() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func (o *Options) Build() error {
	var (
		level Level
		err   error
	)
	if err = level.UnmarshalText([]byte(o.Level)); err != nil {
		return err
	}

	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == ConsoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zc := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: o.Format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:          "message",
			LevelKey:            "level",
			TimeKey:             "timestamp",
			NameKey:             "logger",
			CallerKey:           "caller",
			StacktraceKey:       "stacktrace",
			SkipLineEnding:      false,
			LineEnding:          zapcore.DefaultLineEnding,
			EncodeLevel:         encodeLevel,
			EncodeTime:          timeEncoder,
			EncodeDuration:      milliSecondsDurationEncoder,
			EncodeCaller:        zapcore.ShortCallerEncoder,
			EncodeName:          zapcore.FullNameEncoder,
			NewReflectedEncoder: nil,
			ConsoleSeparator:    "",
		},
		OutputPaths:      o.OutputPaths,
		ErrorOutputPaths: o.ErrorOutputPaths,
		InitialFields:    nil,
	}

	var logger *zap.Logger
	if logger, err = zc.Build(); err != nil {
		return err
	}
	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}
