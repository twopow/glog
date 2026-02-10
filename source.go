package glog

import (
	"runtime"
	"strings"
)

var pkgPrefix = "github.com/twopow/glog."

// outsideCaller returns the PC of the caller outside the given package prefix
func outsideCaller() uintptr {
	const depth = 16

	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:]) // skip Callers + this func
	frames := runtime.CallersFrames(pcs[:n])

	for {
		f, more := frames.Next()
		if !strings.Contains(f.Function, pkgPrefix) {
			return f.PC
		}
		if !more {
			break
		}
	}
	return 0
}
