package zwrap

import (
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

type testWriter struct {
	t                 *testing.T
	needsPrefix       string
	mustNotHavePrefix string
}

var ErrPrefixMismatch = errors.New("prefix mismatch")

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.t.Helper()
	line := strings.TrimSuffix(string(p), "\n")
	if w.needsPrefix != "" && !strings.Contains(line, w.needsPrefix) {
		w.t.Errorf("expected prefix %q, got %q", w.needsPrefix, line)
		return 0, ErrPrefixMismatch
	}
	if w.mustNotHavePrefix != "" && strings.Contains(line, w.mustNotHavePrefix) {
		w.t.Errorf("unexpected prefix %q, got %q", w.mustNotHavePrefix, line)
		return 0, ErrPrefixMismatch
	}
	w.t.Log(line)
	return len(p), nil
}

type aLogger interface {
	Println(...interface{})
}

type needsLogger struct {
	logger aLogger
}

func (n *needsLogger) SetLogger(logger aLogger) {
	n.logger = logger
}

func (n *needsLogger) DoSomething() {
	n.logger.Println("Hello, world!")
}

func TestWrap(t *testing.T) {
	ExampleWrap()

	writah := &testWriter{t: t}

	zl := zerolog.New(writah).With().Timestamp().Logger()
	wrapped := Wrap(zl)
	myThing := &needsLogger{}
	myThing.SetLogger(wrapped)

	multiLog := func(v ...interface{}) {
		t.Helper()
		wrapped.Print(v...)
		wrapped.Printf("%v", v)
		wrapped.Println(v...)
		wrapped.Error(v...)
		wrapped.Errorf("%v", v)
		wrapped.Errorln(v...)
		wrapped.Debug(v...)
		wrapped.Debugf("%v", v)
		wrapped.Debugln(v...)
		wrapped.Warn(v...)
		wrapped.Warnf("%v", v)
		wrapped.Warnln(v...)
		wrapped.Info(v...)
		wrapped.Infof("%v", v)
		wrapped.Infoln(v...)
		wrapped.Tracef("%v", v)
		wrapped.Trace(v...)
		wrapped.Traceln(v...)
		wrapped.Logf("%v", v)
	}

	t.Run("generic", func(t *testing.T) {
		multiLog("Hello, world!")
	})

	t.Run("prefix", func(t *testing.T) {
		writah.needsPrefix = "prefix: "
		wrapped.SetPrefix("prefix: ")
		multiLog("Hello, world!")
	})

	t.Run("remove prefix", func(t *testing.T) {
		writah.needsPrefix = ""
		writah.mustNotHavePrefix = "prefix: "
		wrapped.SetPrefix("")
		multiLog("Hello, world!")
	})

}

func ExampleWrap() {
	// Create a new zerolog.Logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Demonstrate that we can use the stdlib logger
	myThing := &needsLogger{}
	myThing.SetLogger(log.New(os.Stderr, "stdlog: ", log.LstdFlags))
	myThing.DoSomething()

	// Demonstrate that we can use zerolog when wrapped

	/* Before, does not compile:
	myThing.SetLogger(logger)
	myThing.DoSomething()
	*/

	// The zwrap solution, wrap the logger:
	zl := Wrap(logger)
	myThing.SetLogger(zl)
	myThing.DoSomething()
}
