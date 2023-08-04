package zwrap

import (
	"log"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

type testWriter struct {
	t *testing.T
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.t.Log(strings.TrimSuffix(string(p), "\n"))
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

func TestWrapper(t *testing.T) {
	// Create a new zerolog.Logger
	logger := zerolog.New(&testWriter{t: t}).With().Timestamp().Logger()

	// Demonstrate that we can use the stdlib logger
	myThing := &needsLogger{}
	myThing.SetLogger(log.New(&testWriter{t: t}, "stdlog: ", log.LstdFlags))
	myThing.DoSomething()

	// Demonstrate that we can use zerolog when wrapped

	/* Before, does not compile:
	myThing.SetLogger(logger)
	myThing.DoSomething()
	*/

	// The zwrap solution, wrap the logger:
	zl := Wrap(&logger)
	myThing.SetLogger(zl)
	myThing.DoSomething()
}
