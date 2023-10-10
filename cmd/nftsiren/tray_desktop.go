//go:build !android && !ios

package main

import (
	"bytes"
	"image"
	"nftsiren/pkg/log"
	"runtime"
	"sync"

	"fyne.io/systray"
	"github.com/fyne-io/image/ico"
)

var trayRunningMutex sync.Mutex
var trayRunning bool

func systrayIcon(data []byte) []byte {
	if runtime.GOOS == "windows" {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		buf := &bytes.Buffer{}
		err = ico.Encode(buf, img)
		if err != nil {
			panic(err)
		}
		return buf.Bytes()
	}
	return data
}

func TrayStart(iconData []byte) {
	go systray.Run(
		// On ready
		func() {
			log.Debug().Println("Systray started")
			systray.SetIcon(systrayIcon(iconData))
			systray.SetTitle(NAME)
			mShow := systray.AddMenuItem("Show", "Show app")
			mQuit := systray.AddMenuItem("Quit", "Quit from app")
			// go func() {
			trayRunningMutex.Lock()
			trayRunning = true
			trayRunningMutex.Unlock()
			for {
				// log.Debug().Println("Listening for tray events")
				select {
				case <-mShow.ClickedCh:
					// log.Debug().Println("Tray received show event")
					ShowWindowChan <- struct{}{}
				case <-mQuit.ClickedCh:
					// log.Debug().Println("Tray received quit event")
					QuitChan <- struct{}{}
				case <-QuitChan: // This will receive when it is closed
					// log.Debug().Println("Tray received event from QuitChan")
					return
				}
			}
			// }()
		},
		// On exit
		func() {},
	)
}

func TrayRunning() bool {
	trayRunningMutex.Lock()
	defer trayRunningMutex.Unlock()
	return trayRunning
}

func TrayStop() {
	trayRunningMutex.Lock()
	trayRunning = false
	trayRunningMutex.Unlock()
	close(ShowWindowChan)
	close(QuitChan)
	systray.Quit()
	log.Debug().Println("Systray stopped")
}

func isDesktop() bool {
	return true
}

func isMobile() bool {
	return false
}
