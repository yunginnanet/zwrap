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

/*
StdCompatLogger is an interface that provides compatibility with the standard library's log.Logger.

# Original methods

_Note: not all methods are implemented._

func Fatal(v ...interface{})
func Fatalf(format string, v ...interface{})
func Fatalln(v ...interface{})
func Flags() int
func Output(calldepth int, s string) error
func Panic(v ...interface{})
func Panicf(format string, v ...interface{})
func Panicln(v ...interface{})
func Prefix() string
func Print(v ...interface{})
func Printf(format string, v ...interface{})
func Println(v ...interface{})
func SetFlags(flag int)
func SetOutput(w io.Writer)
func SetPrefix(prefix string)
func Writer() io.Writer
*/
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
}

// assert that Logger implements StdCompatLogger and that log.Logger implements StdCompatLogger
var _ StdCompatLogger = &Logger{}
var _ StdCompatLogger = &log.Logger{}

// ----------------------------------------------------

type Logger struct {
	*zerolog.Logger
	prefix string
	*sync.RWMutex
}

func (l *Logger) SetPrefix(prefix string) {
	l.Lock()
	l.prefix = prefix
	l.Logger = nil
	nl := l.Logger.With().Str("caller", prefix).Logger()
	l.Logger = &nl
	l.Unlock()
}

func (l *Logger) Prefix() string {
	l.RLock()
	defer l.RUnlock()
	return l.prefix
}

func (l *Logger) Println(v ...interface{}) {
	l.RLock()
	l.Logger.Info().Msg(fmt.Sprintln(v...))
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

func Wrap(l *zerolog.Logger) *Logger {
	return &Logger{
		Logger:  l,
		RWMutex: &sync.RWMutex{},
	}
}
