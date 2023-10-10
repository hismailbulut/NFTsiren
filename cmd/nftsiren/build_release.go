//go:build !debug

package main

const BUILD_TYPE = "RELEASE"

// Only works in debug builds
func assert(cond bool, format string, args ...any) {}

func isDebug() bool {
	return false
}

func isRelease() bool {
	return true
}
