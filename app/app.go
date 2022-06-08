package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
	"strings"
)

// App 应用基本结构.
type App struct {
	basename    string     // 应用的二进制文件名
	name        string     // 应用名
	description string     // 应用描述
	options     CliOptions // 应用配置
	runFunc     RunFunc    // 应用启动函数
	// todo
	silence   bool
	noVersion bool
	noConfig  bool
	// todo commands
	args cobra.PositionalArgs // args 用于对 cmd 进行参数校验
	cmd  *cobra.Command
}

// Option 应用初始化时使用的可选参数.
type Option func(a *App)

func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

type RunFunc func(basename string) error

func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

func NewApp(name, basename string, options ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range options {
		o(a)
	}

	// todo
	// a.buildCommond
	return a
}

func (a *App) buildCommand() {

}

// FormatBaseName 根据不同的系统将传入的二进制文件名格式化为可执行文件名.
func FormatBaseName(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}
