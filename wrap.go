// Package zwrap provides a wrapper for zerolog.Logger that implements the standard library's log.Logger methods,
// as well as other common logging methods as an attempt to provide compatibility with other logging libraries.
package zwrap

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
	*sync.RWMutex

	prefix     string
	printLevel zerolog.Level
}

func (l *Logger) Warning(args ...any) {
	l.Logger.Warn().Msg(fmt.Sprint(args...))
}

func (l *Logger) Warningln(args ...any) {
	l.Logger.Warn().Msg(fmt.Sprintln(args...))
}

func (l *Logger) V(level int) bool {
	if level > 127 || level < 0 {
		return false
	}
	if l.Logger.GetLevel() > zerolog.Level(int8(level)) {
		return true
	}
	return false
}

func (l *Logger) SetPrefix(prefix string) {
	l.Lock()
	l.prefix = prefix
	l.Unlock()
}

func (l *Logger) SetPrintLevel(level zerolog.Level) {
	l.Lock()
	l.printLevel = level
	l.Unlock()
}

func (l *Logger) Prefix() string {
	l.RLock()
	defer l.RUnlock()
	return l.prefix
}

func (l *Logger) Println(v ...interface{}) {
	l.RLock()
	l.Logger.WithLevel(l.printLevel).Msg(fmt.Sprint(v...))
	l.RUnlock()
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.RLock()
	l.Logger.WithLevel(l.printLevel).Msgf(format, v...)
	l.RUnlock()
}

func (l *Logger) Print(v ...interface{}) {
	l.RLock()
	l.Logger.WithLevel(l.printLevel).Msg(fmt.Sprint(v...))
	l.RUnlock()
}

func (l *Logger) Fatal(v ...interface{}) {
	// Don't check mutex here because we're exiting anyway.
	printLn(l.Logger.Fatal(), v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	// Don't check mutex here because we're exiting anyway.
	l.Logger.Fatal().Msgf(format, v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	// Don't check mutex here because we're exiting anyway.
	printLn(l.Logger.Fatal(), v...)
}

func (l *Logger) Panic(v ...interface{}) {
	// Don't check mutex here because we're panicking anyway.
	printLn(l.Logger.Panic(), v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	// Don't check mutex here because we're panicking anyway.
	l.Logger.Panic().Msgf(format, v...)
}

func (l *Logger) Panicln(v ...interface{}) {
	// Don't check mutex here because we're panicking anyway.
	printLn(l.Logger.Panic(), v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.RLock()
	l.Logger.Error().Msgf(format, v...)
	l.RUnlock()
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.RLock()
	l.Logger.Warn().Msgf(format, v...)
	l.RUnlock()
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.RLock()
	l.Logger.Info().Msgf(format, v...)
	l.RUnlock()
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.RLock()
	l.Logger.Debug().Msgf(format, v...)
	l.RUnlock()
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.RLock()
	l.Logger.Trace().Msgf(format, v...)
	l.RUnlock()
}

func (l *Logger) Error(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Error(), v...)
	l.RUnlock()
}

func (l *Logger) Warn(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Warn(), v...)
	l.RUnlock()
}

func (l *Logger) Info(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Info(), v...)
	l.RUnlock()
}

func (l *Logger) Debug(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Debug(), v...)
	l.RUnlock()
}

func (l *Logger) Trace(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Trace(), v...)
	l.RUnlock()
}

func (l *Logger) Errorln(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Error(), v...)
	l.RUnlock()
}

func (l *Logger) Warnln(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Warn(), v...)
	l.RUnlock()
}

func (l *Logger) Infoln(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Info(), v...)
	l.RUnlock()
}

func (l *Logger) Debugln(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Debug(), v...)
	l.RUnlock()
}

func (l *Logger) Traceln(v ...interface{}) {
	l.RLock()
	printLn(l.Logger.Trace(), v...)
	l.RUnlock()
}

func (l *Logger) Verbosef(format string, v ...interface{}) {
	l.RLock()
	l.Logger.Trace().Msgf(format, v...)
	l.RUnlock()
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.RLock()
	l.Logger.Info().Msgf(format, v...)
	l.RUnlock()
}
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.RLock()
	l.Logger.Warn().Msgf(format, v...)
	l.RUnlock()
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	l.SetPrefix(prefix)
	return l
}

func (l *Logger) Logf(format string, v ...interface{}) {
	l.Printf(format, v...)
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.RLock()
	nl := l.Logger.With().Fields(fields).Logger()
	l.Logger = &nl
	l.RUnlock()
	return l
}

// SetLevel is compatibility for ghettovoice/gosip/log.Logger
func (l *Logger) SetLevel(level uint32) {
	l.Lock()
	nl := l.Logger.Level(gosipLevelToZerologLevel(level))
	l.Logger = &nl
	l.Unlock()
}

func gosipLevelToZerologLevel(level uint32) zerolog.Level {
	switch level {
	case 0:
		return zerolog.PanicLevel
	case 1:
		return zerolog.FatalLevel
	case 2:
		return zerolog.ErrorLevel
	case 3:
		return zerolog.WarnLevel
	case 4:
		return zerolog.InfoLevel
	case 5:
		return zerolog.DebugLevel
	case 6:
		return zerolog.TraceLevel
	}
	panic(fmt.Sprintf("invalid log level %d", level))
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.RLock()
	l.Logger.WithLevel(l.printLevel).Msg(string(bytes.TrimSuffix(p, []byte("\n"))))
	l.RUnlock()
	return len(p), nil
}

func (l *Logger) Output(calldepth int, s string) error {
	l.RLock()
	event := l.Logger.Info()
	if calldepth != 2 {
		if l.prefix != "" {
			zerolog.CallerFieldName = "caller_file"
		}
		event.CallerSkipFrame(calldepth)
		event = event.Caller()
	}
	event.Msg(s)
	zerolog.CallerFieldName = "caller"
	l.RUnlock()
	return nil
}

func printLn(e *zerolog.Event, v ...interface{}) {
	strBuf := strBufs.Get().(*strings.Builder)
	for i, val := range v {
		if i > 0 {
			strBuf.WriteString(" ")
		}
		strBuf.WriteString(fmt.Sprint(val))
	}
	e.Msg(strBuf.String())
	strBuf.Reset()
	strBufs.Put(strBuf)
}

type prefixHook struct {
	parent StdCompatLogger
}

func (h prefixHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	if h.parent.Prefix() != "" {
		e.Str("caller", h.parent.Prefix())
	}
}

func Wrap(l zerolog.Logger) *Logger {
	wrapped := &Logger{
		RWMutex:    &sync.RWMutex{},
		printLevel: zerolog.InfoLevel,
	}
	p := prefixHook{wrapped}
	l = l.Hook(p)
	wrapped.Logger = &l
	return wrapped
}
