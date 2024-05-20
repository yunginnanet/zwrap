package zwrap

import (
	"strings"
	"sync"
)

var strBufs = &sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}
