//go:build android || ios

package main

func TrayStart([]byte) {}

func TrayRunning() bool { return false }

func TrayStop() {}

func isDesktop() bool {
	return false
}

func isMobile() bool {
	return true
}
