# zwrap
[![GoDoc](https://godoc.org/git.tcp.direct/kayos/zwrap?status.svg)](https://godoc.org/git.tcp.direct/kayos/zwrap)
[![Go Report Card](https://goreportcard.com/badge/git.tcp.direct/kayos/zwrap)](https://goreportcard.com/report/git.tcp.direct/kayos/zwrap)


[zwrap](https://git.tcp.direct/kayos/zwrap) is a simple compatibility wrapper around [zerolog](https://github.com/rs/zerolog) that allows for a package to use many different logging libraries _(including stdlib's `log` package)_ without having to import them directly.

## Usage

```go
package main

import (
	"os"
	"log"

	"git.tcp.direct/kayos/zwrap"
	"github.com/rs/zerolog"
)

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

func main() {
	// Create a new zerolog.Logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	
	// Demonstrate that we can use the stdlib logger
	myThing := &needsLogger{}
	myThing.SetLogger(log.New(os.Stdout, "", log.LstdFlags))
	myThing.DoSomething()

	// Demonstrate that we can use zerolog when wrapped

	/* Before, does not compile:
	myThing.SetLogger(logger)
	myThing.DoSomething()
	*/
	
	// The zwrap solution, wrap the logger:
	zl := zwrap.Wrap(logger)
	myThing.SetLogger(zl)
	myThing.DoSomething()
}

```
