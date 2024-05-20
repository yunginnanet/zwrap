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
	mustLevel         *zerolog.Level
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

	if w.mustLevel != nil {
		if !strings.Contains(line, w.mustLevel.String()) {
			w.t.Errorf("expected level %q, got %q", w.mustLevel.String(), line)
			return 0, ErrPrefixMismatch
		}
		lvl := strings.Split(line, `"level":"`)[1]
		lvl = strings.Split(lvl, `"`)[0]
		if lvl != w.mustLevel.String() {
			w.t.Errorf("expected level %q, got %q", w.mustLevel.String(), lvl)
			return 0, ErrPrefixMismatch
		}
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
	n.logger.Println("yeet")
}

type leveled struct {
	name        string
	test        func(*Logger, *testing.T)
	shouldPanic bool
	panicked    bool
	t           *testing.T
}

const (
	expected   = "expected"
	unexpected = "unexpected"
)

func (l *leveled) expS() string {
	if l.shouldPanic {
		return expected
	}
	return unexpected
}

func (l *leveled) Run() {
	l.t.Run(l.name, func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				l.t.Logf("caught panic (%s): %v", l.expS(), r)
				l.panicked = true
			}
		}()
		zl := zerolog.New(os.Stderr).With().Timestamp().Logger()
		wrapped := Wrap(zl)
		l.test(wrapped, t)
	})
	if (l.shouldPanic && !l.panicked) || (!l.shouldPanic && l.panicked) {
		l.t.Errorf("%s panic during test: %s", l.expS(), l.name)
	}
}

func TestWrap(t *testing.T) {
	ExampleWrap()

	writah := &testWriter{t: t}

	zl := zerolog.New(writah).With().Timestamp().Logger()
	wzl := Wrap(zl)
	myThing := &needsLogger{}
	myThing.SetLogger(wzl)

	multiLog := func(wrapped *Logger, v ...interface{}) {
		t.Helper()
		toggleFatals := wrapped.noFatal == false
		togglePanics := wrapped.noPanic == false
		wrapped.Print(v...)
		wrapped.Printf("f: %v", v...)
		wrapped.Println(v...)
		wrapped.Error(v...)
		wrapped.Errorf("f: %v", v...)
		wrapped.Errorln(v...)
		wrapped.Debug(v...)
		wrapped.Debugf("f: %v", v...)
		wrapped.Debugln(v...)
		wrapped.Warn(v...)
		wrapped.Warnf("f: %v", v...)
		wrapped.Warnln(v...)
		wrapped.Info(v...)
		wrapped.Infof("f: %v", v...)
		wrapped.Infoln(v...)
		wrapped.Tracef("f: %v", v...)
		wrapped.Trace(v...)
		wrapped.Traceln(v...)
		wrapped.Logf("f: %v", v...)
		wrapped.Warning(v...)
		wrapped.Warningf("f: %v", v...)
		wrapped.Println("")
		wrapped.Println()
		if toggleFatals {
			wrapped.NoFatals(true)
		}
		wrapped.Fatal(v...)
		wrapped.Fatalf("f: %v", v...)
		wrapped.Fatalln(v...)
		if toggleFatals {
			wrapped.NoFatals(false)
		}
		if togglePanics {
			wrapped.NoPanics(true)
		}
		wrapped.Panic(v...)
		wrapped.Panicf("f: %v", v...)
		wrapped.Panicln(v...)
		if togglePanics {
			wrapped.NoPanics(false)
		}
	}

	if wzl.V(0) {
		t.Error("V(0) should always return false")
	}

	t.Run("generic", func(t *testing.T) {
		multiLog(wzl, "yeet")
	})

	t.Run("prefix", func(t *testing.T) {
		writah.needsPrefix = "prefix: "
		wzl.SetPrefix("prefix: ")
		multiLog(wzl, "yeet")
		writah.needsPrefix = "prefix2: "
		wzl = wzl.WithPrefix("prefix2: ")
		multiLog(wzl, "yeet")
	})

	t.Run("remove prefix", func(t *testing.T) {
		writah.needsPrefix = ""
		writah.mustNotHavePrefix = "prefix: "
		wzl.SetPrefix("")
		multiLog(wzl, "yeet")
	})

	forceLevelTests := []leveled{
		{
			name: "trace",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.ForceLevel(zerolog.TraceLevel)
				multiLog(wrapped, "yeet")
			},
			shouldPanic: false,
			t:           t,
		},
		{
			name: "trace_with_panic",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.ForceLevel(zerolog.TraceLevel)
				multiLog(wrapped, "yeet")
				wrapped.Panic("yeet")
			},
			shouldPanic: true,
			t:           t,
		},
		{
			name: "debug",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.ForceLevel(zerolog.DebugLevel)
				multiLog(wrapped, "yeet")
			},
			shouldPanic: false,
			t:           t,
		},
		{
			name: "info",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.ForceLevel(zerolog.InfoLevel)
				multiLog(wrapped, "yeet")
			},
			shouldPanic: false,
			t:           t,
		},
		{
			name: "warn",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.ForceLevel(zerolog.WarnLevel)
				multiLog(wrapped, "yeet")
			},
			shouldPanic: false,
			t:           t,
		},
		{
			name: "error",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.ForceLevel(zerolog.ErrorLevel)
				multiLog(wrapped, "yeet")
			},
			shouldPanic: false,
			t:           t,
		},
		{
			name: "fatal",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.ForceLevel(zerolog.FatalLevel)
				wrapped.NoFatals(true)
				multiLog(wrapped, "yeet")
			},
			shouldPanic: false,
			t:           t,
		},
		{
			name: "panic",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.ForceLevel(zerolog.PanicLevel)
				multiLog(wrapped, "yeet")
			},
			shouldPanic: true,
			t:           t,
		},
		{
			name: "panic_with_panic",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.Panic("yeet")
			},
			shouldPanic: true,
			t:           t,
		},
	}

	for _, test := range forceLevelTests {
		test.name = "force_level_" + test.name
		test.Run()
	}

	panicAndFatalBypassTests := []leveled{
		{
			name: "no_panic",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.NoPanics(true)
				wrapped.Panic("yeet!!")
			},
			shouldPanic: false,
			t:           t,
		},
		{
			name: "no_fatal",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.NoFatals(true)
				wrapped.Fatal("yeet")
			},
			shouldPanic: false, // I guess the test should fail if it os.Exits anyway..? :^)
			t:           t,
		},
		{
			name: "no_panic_no_fatal_force_panic",
			test: func(wrapped *Logger, t *testing.T) {
				wrapped.NoPanics(true)
				wrapped.NoFatals(true)
				wrapped.ForceLevel(zerolog.PanicLevel)
				multiLog(wrapped, "yeet!!")
				wrapped.Panic("yeet!!")
			},
			shouldPanic: false,
			t:           t,
		},
	}

	for _, test := range panicAndFatalBypassTests {
		test.name = "bypass_" + test.name
		test.Run()
	}

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
