package zwrap

func (l *Logger) checkPanicBypass(fmt string, v ...interface{}) (string, []interface{}, bool) {
	if !l.noPanic {
		return fmt, v, true
	}

	if fmt != "" {
		fmt = "[PANIC BYPASSED] " + fmt
		return fmt, v, false
	}

	nv := make([]interface{}, len(v)+1)
	nv[0] = "[PANIC BYPASSED]"
	copy(nv[1:], v)
	v = nil
	return fmt, nv, false
}

func (l *Logger) checkFatalBypass(fmt string, v ...interface{}) (string, []interface{}, bool) {
	if !l.noFatal {
		return fmt, v, true
	}

	if fmt != "" {
		fmt = "[FATAL BYPASSED] " + fmt
		return fmt, v, false
	}

	nv := make([]interface{}, len(v)+1)
	nv[0] = "[FATAL BYPASSED]"
	copy(nv[1:], v)
	v = nil
	return fmt, nv, false
}

func (l *Logger) NoPanics(b bool) {
	l.mu.Lock()
	l.noPanic = b
	l.mu.Unlock()
}

func (l *Logger) NoFatals(b bool) {
	l.mu.Lock()
	l.noFatal = b
	l.mu.Unlock()
}

func (l *Logger) WithNoPanics() *Logger {
	l.NoPanics(true)
	return l
}

func (l *Logger) WithNoFatals() *Logger {
	l.NoFatals(true)
	return l
}

func (l *Logger) ForceLevel(level any) {
	l.mu.Lock()
	nl := castToZlogLevel(level)
	l.forceLevel = &nl
	l.printLevel = nl
	nll := l.Logger.Level(nl)
	l.Logger = &nll
	l.mu.Unlock()
}

func (l *Logger) WithForceLevel(level any) *Logger {
	l.ForceLevel(level)
	return l
}
