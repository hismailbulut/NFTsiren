package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"time"

	"nftsiren/cmd/nftsiren/assets"
	"nftsiren/cmd/nftsiren/cache"
	"nftsiren/cmd/nftsiren/config"
	"nftsiren/pkg/bench"
	"nftsiren/pkg/log"
	"nftsiren/pkg/util"
	"nftsiren/pkg/worker"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

const (
	NAME    = "NFTsiren"
	WEBSITE = "https://www.github.com/hismailbulut/NFTsiren"
)

var (
	// This is injected at compile time
	// making this a constant brokes injection
	// do not change it
	VERSION = ""
)

func main() {
	go run()
	app.Main()
}

var (
	ShowWindowChan    = make(chan struct{})
	RefreshWindowChan = make(chan struct{}, 256) // buffered, because should be callable inside from Layout
	QuitChan          = make(chan struct{})
)

func run() {
	// exit at the end
	var exitCode int
	defer func() { os.Exit(exitCode) }()

	// Capture any panic at the end and generate crash report
	defer RecoverAndReportPanic()

	if isDebug() {
		// print benchmark results at the end
		defer bench.PrintResults(&log.Default)
		// start pprof
		// pprof package imported in ./build_debug.go
		go func() {
			log.Info().Println(http.ListenAndServe("localhost:6060", nil))
		}()
		// log to the local server on mobile
		if isMobile() {
			// connect to the debug server if available and send all log output to it
			conn, err := connectUdpServer("192.168.1.34:19876")
			if err != nil {
				log.Error().Println("Failed to connect logging server:", err)
			} else {
				defer conn.Close()
				log.SetOutput(conn)
			}
		}
	} else {
		log.SetMinLevel(log.LevelNone)
		// log.SetMinLevel(log.LevelError)
	}

	dataDir, err := app.DataDir()
	if err != nil {
		log.Error().Println("Couldn't find data dir:", err)
		// TODO: handle this
	}

	const prefFileName = "preferences.json"

	log.Debug().Println("Preferences path:", filepath.Join(dataDir, NAME, prefFileName))

	// Init preferences
	err = config.Init(dataDir, NAME, prefFileName)
	if err != nil {
		log.Error().Println("Failed to init preferences:", err)
	}

	if isDesktop() {
		// Check if there is existing instance
		quit, clearInstancePort := checkExistingInstance()
		if quit {
			return
		}
		defer clearInstancePort()
	}

	var cacheDir string
	if isMobile() {
		// TODO: we have to find application cache directory in mobile
		// this will work but not a correct way to do it
		cacheDir = filepath.Join(dataDir, "cache")
	} else {
		cacheDir, _ = os.UserCacheDir()
		cacheDir = filepath.Join(cacheDir, NAME)
	}

	log.Debug().Println("Cache dir:", cacheDir)

	// Init cache
	err = cache.Init(cacheDir)
	if err != nil {
		log.Error().Println("Failed to init cache:", err)
	}

	// For debugging
	log.Info().Println("Execpath:", util.ExecPath())
	log.Info().Println("Version:", VERSION)
	log.Info().Println("BuildType:", BUILD_TYPE)

	// Create daemon
	daemon := NewDaemon()
	daemon.Start()
	defer daemon.Stop()

	mainLoop(daemon)
}

func mainLoop(daemon *Daemon) {
	// start systray
	TrayStart(assets.NftsirenLogoData)
	defer TrayStop()

	// safely quit on ctrl-c
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// only first time window visibility
	showWindow := true
	// background service doesn't work on mobile
	if isDesktop() && config.LoadFallback("startInBackground", false) {
		showWindow = false
	}

	// run main loop
	for {
		if showWindow {
			err, exit := windowLoop(ctx, daemon)
			if err != nil {
				log.Error().Println("Window exited with error:", err)
				return
			}
			if exit || isMobile() {
				// Do not continue
				return
			}
			// window is hidden now, release big ui resources (images)
			daemon.ReleaseResources()
		}
		showWindow = true
		// window is hidden, listen for tray events
		log.Debug().Println("Window is hidden now")
	listener:
		for {
			select {
			case <-ShowWindowChan:
				log.Debug().Println("Window is shown")
				// window is shown again, reload resources
				daemon.ReloadResources()
				break listener
			case <-RefreshWindowChan:
			// Consume this here but we don't refresh it because it is closed
			case <-QuitChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}
}

func windowLoop(ctx context.Context, daemon *Daemon) (error, bool) {
	// this is required for windows
	BeginHighresTimer()
	defer EndHighresTimer()

	theme := DefaultTheme()

	const width = 500
	const height = 800
	w := app.NewWindow(
		// Desktop
		app.Title(NAME),
		app.Size(width, height),
		app.MinSize(width, height),
		app.MaxSize(width*1.5, height*1.25),
		// Mobile
		app.NavigationColor(theme.Bg),
		app.StatusColor(theme.Bg),
		app.PortraitOrientation.Option(),
	)

	ui := NewUI(daemon, theme)
	var ops op.Ops

	var fps atomic.Int32
	fpscounter := worker.New(worker.Settings{
		Name:     "FpsCounter",
		Interval: time.Second,
		Work: func() {
			if fps.Load() > 60 {
				log.Debug().Printf("High FPS: %d\n", fps.Load())
			}
			fps.Store(0)
		},
	})
	fpscounter.Start()
	defer fpscounter.Stop()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err, false
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				ui.Layout(gtx)

				end := bench.Begin()
				e.Frame(&ops)
				end("Frame")

				fps.Store(fps.Load() + 1)
			}
		case <-ShowWindowChan:
			log.Debug().Println("Window is already visible")
		case <-RefreshWindowChan:
			w.Invalidate()
		case <-QuitChan:
			return nil, true
		case <-ctx.Done():
			return nil, true
		}
	}
}

func connectUdpServer(addr string) (net.Conn, error) {
	udpServer, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// App should exit if this function returns true
// Always defer the returned function
func checkExistingInstance() (bool, func()) {
	port := config.LoadFallback("port", "")
	if port != "" {
		err := ShowExistingInstance(port)
		if err != nil {
			log.Error().Println("Failed to show existing instace:", err)
		} else {
			log.Info().Println("There is already an open instance")
			return true, func() {}
		}
	}
	// Either there is no server or we failed, create a server
	server, err := NewIpcServer()
	if err != nil {
		log.Error().Println("Failed to create IPC server:", err)
		return false, func() {}
	}
	// Get the port server listens on
	port, err = server.Port()
	if err != nil {
		log.Error().Println("Invalid IPC server port:", err)
		return false, func() {}
	}
	// Set server port in preferences file
	config.Store("port", port)
	config.Save()
	log.Info().Println("IPC server created and listens on port", port)
	// This is closer
	return false, func() {
		config.Store("port", "")
		err := config.Save()
		if err != nil {
			log.Error().Println("Failed to save preferences:", err)
		}
		server.Close()
	}
}

func TODO(msg string) {
	log.Debug().Println("TODO:", msg)
}
