//go:build android || ios

package config

func GetAutostart() bool {
	// Not supported
	return false
}

func SetAutostart(enabled bool) {
	// Not supported
}

func fixAutostart() {
	// Not supported
}
