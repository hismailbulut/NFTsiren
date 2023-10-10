//go:build !android && !ios

package config

import (
	"nftsiren/pkg/log"
	"nftsiren/pkg/util"

	"github.com/emersion/go-autostart"
)

func getAutostartApp() *autostart.App {
	return &autostart.App{
		Name:        state.Name.Load(),
		DisplayName: state.Name.Load(),
		Exec:        []string{util.ExecPath()},
	}
}

func GetAutostart() bool {
	return getAutostartApp().IsEnabled()
}

func SetAutostart(enabled bool) {
	appInfo := getAutostartApp()
	if enabled && !appInfo.IsEnabled() {
		err := appInfo.Enable()
		if err != nil {
			log.Error().Println("Failed to enable autostart:", err)
		} else {
			Store("autostartPath", util.ExecPath())
		}
	} else if !enabled && appInfo.IsEnabled() {
		err := appInfo.Disable()
		if err != nil {
			log.Error().Println("Failed to disable autostart:", err)
		} else {
			Delete("autostartPath")
		}
	}
}

func fixAutostart() {
	path, err := Load[string]("autostartPath")
	if err == nil && path != util.ExecPath() {
		log.Warn().Println("Autostart path incorrect, fixing it now")
		SetAutostart(true)
	}
}
