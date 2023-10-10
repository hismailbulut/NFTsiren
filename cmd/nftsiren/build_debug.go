//go:build debug

package main

import (
	"fmt"
	_ "net/http/pprof"
)

const BUILD_TYPE = "DEBUG"

// Only works in debug builds
func assert(cond bool, format string, args ...any) {
	if !cond {
		panic(fmt.Errorf(format, args...))
	}
}

func isDebug() bool {
	return true
}

func isRelease() bool {
	return false
}
