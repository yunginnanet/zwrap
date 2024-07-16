package zwrap

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

type GRPCCompatLogger interface {
	Info(args ...any)
	Infoln(args ...any)
	Infof(format string, args ...any)
	Warning(args ...any)
	Warningln(args ...any)
	Warningf(format string, args ...any)
	Error(args ...any)
	Errorln(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalln(args ...any)
	Fatalf(format string, args ...any)
	V(l int) bool
}

type ZWrapLogger interface {
	StdCompatLogger
	GRPCCompatLogger
}

// assert that Logger implements StdCompatLogger and GRPCCompatLogger.
var (
	_ StdCompatLogger  = &Logger{}
	_ GRPCCompatLogger = &Logger{}
	_ ZWrapLogger      = &Logger{}
)
