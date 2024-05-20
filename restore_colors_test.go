package zwrap

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
)

type colorTester struct {
	t    *testing.T
	last []byte
}

func (c *colorTester) Write(p []byte) (n int, err error) {
	if len(c.last) == 0 {
		c.last = make([]byte, len(p))
		copy(c.last, p)
		return len(p), nil
	}
	if bytes.Equal(c.last, p) {
		c.t.Errorf("\ncaught second output: %s\nwhich is the same as the first: %s", string(p), string(c.last))
	}
	return len(p), nil
}

func TestLegacyColorizer(t *testing.T) {
	tw := &colorTester{t: t}
	zlc := zerolog.NewConsoleWriter()
	zlc.Out = tw
	zlc.NoColor = false
	zl := Wrap(zerolog.New(zlc))
	zl.Trace("yeet")
	if len(tw.last) == 0 {
		t.Fatalf("test writer busted")
	}
	zlc = zerolog.NewConsoleWriter()
	zlc.FormatLevel = LogLevelFmt(false)
	zlc.NoColor = false
	zlc.Out = tw
	zl = Wrap(zerolog.New(zlc))
	zl.Trace("yeet")
}
