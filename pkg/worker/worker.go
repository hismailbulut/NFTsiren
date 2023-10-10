package worker

import (
	"fmt"
	"runtime/debug"
	"time"

	"nftsiren/pkg/log"
)

type Settings struct {
	Name              string        // Used only for logging
	Interval          time.Duration // Must be set
	Work              func()        // Must be non-nil
	InitialRun        bool          //
	PanicHanler       func(error)   //
	RestartAfterPanic bool          //
	MaxPanics         int           // 0 means unlimited
	EnableLogging     bool
}

type Worker struct {
	Settings
	stop chan struct{}
}

func New(settings Settings) *Worker {
	if settings.Interval <= 0 {
		panic("Worker interval must bigger than zero")
	}
	if settings.Work == nil {
		panic("Worker func can not be nil")
	}
	if settings.EnableLogging {
		log.Debug().Println("Worker", settings.Name, "created")
	}
	return &Worker{Settings: settings, stop: make(chan struct{})}
}

func (w *Worker) start(panicCount int) {
	defer w.recover(panicCount)
	if w.EnableLogging {
		log.Debug().Println("Worker", w.Name, "started")
	}
	if w.InitialRun {
		begin := time.Now()
		w.Work()
		if w.EnableLogging {
			log.Debug().Println("Worker", w.Name, "initially worked in", time.Since(begin))
		}
	}
	prev := time.Now()
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.stop:
			if w.EnableLogging {
				log.Debug().Println("Worker", w.Name, "stopped")
			}
			return
		case <-ticker.C:
			begin := time.Now()
			w.Work()
			if w.EnableLogging {
				log.Debug().Println("Worker", w.Name, "worked in", time.Since(begin), "after", time.Since(prev))
			}
			prev = begin
		}
	}
}

func (w *Worker) recover(panicCount int) {
	rec := recover()
	if rec != nil {
		panicCount++
		err := fmt.Errorf("%v", rec)
		log.Error().Println("Worker", w.Name, "panicked:", err)
		log.Error().Println("Stack Trace:\n", string(debug.Stack()))
		if w.PanicHanler != nil {
			w.PanicHanler(err)
		}
		// restart this worker
		if w.RestartAfterPanic && (w.MaxPanics <= 0 || panicCount < w.MaxPanics) {
			go w.start(panicCount)
		}
	}
}

func (w *Worker) Start() {
	go w.start(0)
}

func (w *Worker) Stop() {
	close(w.stop)
}
