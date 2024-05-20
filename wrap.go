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
	mu *sync.RWMutex

	prefix     string
	printLevel zerolog.Level
	forceLevel *zerolog.Level
	noPanic    bool
	noFatal    bool
}

func (l *Logger) Warning(args ...any) {
	l.printLn(l.Logger.Warn(), false, args...)
}

func (l *Logger) Warningln(args ...any) {
	l.printLn(l.Logger.Warn(), false, args...)
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
	l.mu.Lock()
	l.prefix = prefix
	l.mu.Unlock()
}

func (l *Logger) SetPrintLevel(level zerolog.Level) {
	l.mu.Lock()
	l.printLevel = level
	l.mu.Unlock()
}

func (l *Logger) Prefix() string {
	l.mu.RLock()
	p := l.myPrefix()
	l.mu.RUnlock()
	return p
}

func (l *Logger) myPrefix() string {
	return l.prefix
}

func (l *Logger) Println(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.WithLevel(l.printLevel), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.mu.RLock()
	var str string
	switch {
	case len(v) == 0:
		str = format
	case len(v) == 1:
		str = fmt.Sprintf(format, v[0])
	default:
		str = fmt.Sprintf(format, v...)
	}
	l.printLn(l.Logger.WithLevel(l.printLevel), false, str)
	l.mu.RUnlock()
}

func (l *Logger) Print(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.WithLevel(l.printLevel), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Fatal(v ...interface{}) {
	var ok bool
	if _, v, ok = l.checkFatalBypass("", v...); ok {
		l.printLn(l.Logger.Fatal(), true, v...)
		return
	}
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	var ok bool
	if format, v, ok = l.checkFatalBypass(format, v...); ok {
		l.printLn(l.Logger.Fatal(), true, fmt.Sprintf(format, v...))
		return
	}
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) Fatalln(v ...interface{}) {
	var ok bool
	if _, v, ok = l.checkFatalBypass("", v...); ok {
		l.printLn(l.Logger.Fatal(), true, v...)
		return
	}
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Panic(v ...interface{}) {
	var ok bool
	if _, v, ok = l.checkPanicBypass("", v...); ok {
		l.printLn(l.Logger.Panic(), true, v...)
		return
	}
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	var ok bool
	if format, v, ok = l.checkPanicBypass(format, v...); ok {
		l.printLn(l.Logger.Panic(), true, fmt.Sprintf(format, v...))
		return
	}
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) Panicln(v ...interface{}) {
	var ok bool
	if _, v, ok = l.checkPanicBypass("", v...); ok {
		l.printLn(l.Logger.Panic(), true, v...)
		return
	}
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Warn(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Info(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Debug(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Trace(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) Error(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Warn(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Warn(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Info(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Info(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Debug(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Debug(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Trace(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Trace(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Errorln(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Error(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Warnln(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Warn(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Infoln(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Info(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Debugln(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Debug(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Traceln(v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Trace(), false, v...)
	l.mu.RUnlock()
}

func (l *Logger) Verbosef(format string, v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Trace(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Info(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.mu.RLock()
	l.printLn(l.Logger.Warn(), false, fmt.Sprintf(format, v...))
	l.mu.RUnlock()
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	l.SetPrefix(prefix)
	return l
}

func (l *Logger) Logf(format string, v ...interface{}) {
	l.Printf(format, v...)
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.mu.RLock()
	nl := l.Logger.With().Fields(fields).Logger()
	l.Logger = &nl
	l.mu.RUnlock()
	return l
}

// SetLevel is compatibility for ghettovoice/gosip/log.Logger
func (l *Logger) SetLevel(level any) {
	l.mu.Lock()
	nl := l.Logger.Level(castToZlogLevel(level))
	l.Logger = &nl
	l.mu.Unlock()
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.mu.RLock()
	l.Logger.WithLevel(l.printLevel).Msg(string(bytes.TrimSuffix(p, []byte("\n"))))
	l.mu.RUnlock()
	return len(p), nil
}

func (l *Logger) Output(calldepth int, s string) error {
	l.mu.RLock()
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
	l.mu.RUnlock()
	return nil
}

func (l *Logger) transformZEvent(e *zerolog.Event) *zerolog.Event {
	switch *l.forceLevel {
	case zerolog.PanicLevel:
		if !l.noPanic {
			e = l.Logger.Panic()
		}
	case zerolog.FatalLevel:
		if !l.noFatal {
			e = l.Logger.Fatal()
		}
	case zerolog.ErrorLevel:
		e = l.Logger.Error()
	case zerolog.WarnLevel:
		e = l.Logger.Warn()
	case zerolog.InfoLevel:
		e = l.Logger.Info()
	case zerolog.DebugLevel:
		e = l.Logger.Debug()
	case zerolog.TraceLevel:
		e = l.Logger.Trace()
	default:
		panic(fmt.Sprintf("invalid logger config, bad force level %v", l.forceLevel))
	}
	return e
}

func (l *Logger) printLn(e *zerolog.Event, preserve bool, v ...interface{}) {
	if l.forceLevel != nil && !preserve {
		e = l.transformZEvent(e)
	}
	if len(v) == 0 {
		e.Msg("")
		return
	}
	if len(v) == 1 {
		e.Msg(fmt.Sprint(v[0]))
		return
	}
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
	parent *Logger
}

func (h prefixHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	if h.parent.Prefix() != "" {
		e.Str("caller", h.parent.myPrefix())
	}
}

func Wrap(l zerolog.Logger) *Logger {
	wrapped := &Logger{
		mu:         &sync.RWMutex{},
		printLevel: zerolog.InfoLevel,
	}
	p := prefixHook{wrapped}
	l = l.Hook(p)
	wrapped.Logger = &l
	return wrapped
}
