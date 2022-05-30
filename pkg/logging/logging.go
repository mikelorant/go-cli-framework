package logging

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/leaanthony/go-ansi-parser"
)

type Logger struct {
	buffer  *bytes.Buffer
	logger  *log.Logger
	level   Level
	options Options
}

type Level uint32

type Options struct {
	level bool
	time  bool
	field string
}

type OptionsAdapter func(*Options)

const (
	ErrorLevel Level = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

func New(ctx context.Context, opts ...OptionsAdapter) *Logger {
	var buffer bytes.Buffer

	logger := &Logger{
		buffer:  &buffer,
		logger:  log.New(&buffer, "", 0),
		level:   InfoLevel,
		options: Options{},
	}

	for _, o := range opts {
		o(&logger.options)
	}

	return logger
}

func WithLevel(isEnabled bool) func(*Options) {
	return func(o *Options) {
		o.level = isEnabled
	}
}

func WithTime(isEnabled bool) func(*Options) {
	return func(o *Options) {
		o.time = isEnabled
	}
}

func WithField(field string) func(*Options) {
	return func(o *Options) {
		o.field = field
	}
}

func (lw *Logger) Error(v ...any) { lw.print(ErrorLevel, v...) }
func (lw *Logger) Warn(v ...any)  { lw.print(WarnLevel, v...) }
func (lw *Logger) Info(v ...any)  { lw.print(InfoLevel, v...) }
func (lw *Logger) Debug(v ...any) { lw.print(DebugLevel, v...) }

func (lw *Logger) Errorf(format string, v ...any) { lw.printf(ErrorLevel, format, v...) }
func (lw *Logger) Warnf(format string, v ...any)  { lw.printf(WarnLevel, format, v...) }
func (lw *Logger) Infof(format string, v ...any)  { lw.printf(InfoLevel, format, v...) }
func (lw *Logger) Debugf(format string, v ...any) { lw.printf(DebugLevel, format, v...) }

func (lw *Logger) SetLevel(level string) error {
	lvl, err := parseLevel(level)
	if err != nil {
		return fmt.Errorf("log level must be: %v received %w", []Level{
			ErrorLevel,
			WarnLevel,
			InfoLevel,
			DebugLevel,
		}, err)
	}

	switch lvl {
	case ErrorLevel:
		lw.level = ErrorLevel
	case WarnLevel:
		lw.level = WarnLevel
	case InfoLevel:
		lw.level = InfoLevel
	case DebugLevel:
		lw.level = DebugLevel
	}

	return nil
}

func (lw *Logger) SetOption(opts ...OptionsAdapter) {
	for _, o := range opts {
		o(&lw.options)
	}
}

func (lw *Logger) print(lvl Level, v ...any) {
	if lw.level < lvl {
		return
	}

	lw.logger.Print(v...)
	lw.output(lvl)
}

func (lw *Logger) printf(lvl Level, format string, v ...any) {
	if lw.level < lvl {
		return
	}

	lw.logger.Printf(format, v...)
	lw.output(lvl)
}

func (lw *Logger) output(lvl Level) {
	scanner := bufio.NewScanner(lw.buffer)
	prefix := lw.buildPrefix(lvl)
	length, _ := ansi.Length(prefix)
	stdout := log.New(os.Stdout, "", 0)

	i := 0
	for scanner.Scan() {
		switch i {
		case 0:
			stdout.SetPrefix(prefix)
		default:
			stdout.SetPrefix(strings.Repeat(" ", length))
		}

		stdout.Print(scanner.Text())
		color.Unset()
		i++
	}
}

func (lw *Logger) buildPrefix(lvl Level) string {
	var prefix string

	if lw.options.level {
		prefix += fmt.Sprintf("%v │ ", formatLevel(lvl))
	}

	if lw.options.time {
		prefix += fmt.Sprintf("%v │ ", formatTime())
	}

	if lw.options.field != "" {
		prefix += fmt.Sprintf("%v │ ", lw.options.field)
	}

	return prefix
}

func (ll Level) String() string {
	return [4]string{
		"error",
		"warn",
		"info",
		"debug",
	}[ll]
}

func parseLevel(level string) (lvl Level, err error) {
	levels := map[string]Level{
		ErrorLevel.String(): ErrorLevel,
		WarnLevel.String():  WarnLevel,
		InfoLevel.String():  InfoLevel,
		DebugLevel.String(): DebugLevel,
	}

	lvl, ok := levels[strings.ToLower(level)]
	if !ok {
		return 0, fmt.Errorf("unknown log level: %v", level)
	}

	return lvl, nil
}

func formatLevel(lvl Level) string {
	colours := map[Level]func(format string, a ...interface{}) string{
		ErrorLevel: color.New(color.FgWhite, color.BgRed).SprintfFunc(),
		WarnLevel:  color.New(color.FgBlack, color.BgYellow).SprintfFunc(),
		InfoLevel:  color.New(color.FgBlack, color.BgWhite).SprintfFunc(),
		DebugLevel: color.New(color.FgWhite, color.BgBlue).SprintfFunc(),
	}

	colour := colours[lvl]

	formats := map[Level]string{
		ErrorLevel: colour("  ERROR  "),
		WarnLevel:  colour("  WARN   "),
		InfoLevel:  colour("  INFO   "),
		DebugLevel: colour("  DEBUG  "),
	}

	return formats[lvl]
}

func formatTime() string {
	return time.Now().In(time.Local).Format("15:04:05")
}
