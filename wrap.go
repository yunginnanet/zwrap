// Package zwrap provides a wrapper for zerolog.Logger that implements the standard library's log.Logger methods,
// as well as other common logging methods as an attempt to provide compatibility with other logging libraries.
package zwrap

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

var strBufs = &sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

// StdCompatLogger is an interface that provides compatibility with the standard library's log.Logger.
type StdCompatLogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Prefix() string
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	SetPrefix(prefix string)
	Output(calldepth int, s string) error
}

// assert that Logger implements StdCompatLogger and that log.Logger implements StdCompatLogger
var _ StdCompatLogger = &Logger{}
var _ StdCompatLogger = &log.Logger{}

// ----------------------------------------------------

type Logger struct {
	*zerolog.Logger
	*sync.RWMutex

	prefix     string
	printLevel zerolog.Level
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
	l.Logger.WithLevel(l.printLevel).Msg(fmt.Sprintln(v...))
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
	for i, v := range v {
		if i > 0 {
			strBuf.WriteString(" ")
		}
		strBuf.WriteString(fmt.Sprint(v))
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
