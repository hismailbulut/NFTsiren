package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"nftsiren/cmd/nftsiren/alerts"
	"nftsiren/cmd/nftsiren/config"
	"nftsiren/pkg/apis/etherscan"
	"nftsiren/pkg/apis/looksrare"
	"nftsiren/pkg/apis/magiceden"
	"nftsiren/pkg/apis/opensea"
	"nftsiren/pkg/bench"
	"nftsiren/pkg/log"
	"nftsiren/pkg/nft"
	"nftsiren/pkg/number"
	"nftsiren/pkg/worker"
)

// Daemon is responsible for fetching data from network and checking for alerts
type Daemon struct {
	// This constantly fetches current ethereum price and gas
	GasTracker *GasTracker
	// This will check ethereum and gas alerts every 10 second
	EthGasChecker *worker.Worker
	ethAlerts     *AlertList // *Alert[alerts.EthereumAlert]
	gasAlerts     *AlertList // *Alert[alerts.GasAlert]
	// Collections has their own workers and they are responsible for checking their alerts
	collectionsMutex sync.RWMutex
	collections      []*Collection
}

func NewDaemon() *Daemon {
	daemon := &Daemon{
		GasTracker:  NewGasTracker(),
		ethAlerts:   new(AlertList),
		gasAlerts:   new(AlertList),
		collections: make([]*Collection, 0),
	}
	daemon.EthGasChecker = worker.New(worker.Settings{
		Name:        "Eth&GasChecker",
		Interval:    time.Second * 10,
		Work:        daemon.CheckEthAndGasAlerts,
		InitialRun:  true,
		PanicHanler: ReportPanic,
	})
	return daemon
}

func (daemon *Daemon) Start() error {
	daemon.ResetApiKeys()
	// Start gas tracker
	daemon.GasTracker.Start()
	// Start checking ethereum and gas alarms
	daemon.EthGasChecker.Start()
	// Load everything from user configuration saved in our servers
	daemon.LoadConfig()
	log.Debug().Println("Daemon started")
	return nil
}

func (daemon *Daemon) Stop() {
	daemon.GasTracker.Stop()
	daemon.EthGasChecker.Stop()
	// Stop every collection worker
	daemon.collectionsMutex.RLock()
	for _, c := range daemon.collections {
		c.Stop()
	}
	daemon.collectionsMutex.RUnlock()
	// Now save everything
	daemon.SaveConfig()
	log.Debug().Println("Daemon stopped")
}

func (daemon *Daemon) ResetApiKeys() {
	// keys := config.GetApiKeys()
	keys, err := config.Load[ApiKeys]("apiKeys")
	if err == nil {
		etherscan.SetApiKey(keys.Etherscan)
		opensea.SetApiKey(keys.Opensea)
		looksrare.SetApiKey(keys.Looksrare)
		magiceden.SetApiKey(keys.Magiceden)
	}
}

// This will be called when ui is going to background
// Currently releases images of all collections
func (daemon *Daemon) ReleaseResources() {
	daemon.collectionsMutex.Lock()
	for _, c := range daemon.collections {
		// log.Debug().Println("Removing image of:", c)
		go c.setImage(nil, nil) // we need to go because it refreshes window
	}
	daemon.collectionsMutex.Unlock()
	// This is a fast way to release all unusued memory at once
	debug.FreeOSMemory()
}

// Reloads all images of collections
func (daemon *Daemon) ReloadResources() {
	daemon.collectionsMutex.Lock()
	for _, c := range daemon.collections {
		// log.Debug().Println("Re-fetching image of:", c)
		c.reFetchImage()
	}
	daemon.collectionsMutex.Unlock()
}

func (daemon *Daemon) CheckEthAndGasAlerts() {
	daemon.ethAlerts.ForEach(func(index int, alert alerts.Alert) {
		alert.Check(number.Number{})
	})
	daemon.gasAlerts.ForEach(func(index int, alert alerts.Alert) {
		alert.Check(number.Number{})
	})
}

func (daemon *Daemon) AddCollection(collection *Collection) bool {
	daemon.collectionsMutex.Lock()
	defer daemon.collectionsMutex.Unlock()
	// Check whether we already have this collection
	for _, c := range daemon.collections {
		if c.String() == collection.String() {
			// We can't add this
			return false
		}
	}
	// Append, start it and return true
	daemon.collections = append(daemon.collections, collection)
	collection.Start()
	return true
}

func (daemon *Daemon) RemoveCollection(collection *Collection) bool {
	daemon.collectionsMutex.Lock()
	defer daemon.collectionsMutex.Unlock()
	// Check whether we already have this collection
	for i, c := range daemon.collections {
		if c.String() == collection.String() {
			assert(c == collection, "found different pointers to the same collection")
			c.Stop()
			daemon.collections = append(daemon.collections[:i], daemon.collections[i+1:]...)
			return true
		}
	}
	// We couldn't find it
	return false
}

func (daemon *Daemon) CollectionCount() int {
	daemon.collectionsMutex.RLock()
	defer daemon.collectionsMutex.RUnlock()
	return len(daemon.collections)
}

func (daemon *Daemon) CollectionAtIndex(index int) *Collection {
	daemon.collectionsMutex.RLock()
	defer daemon.collectionsMutex.RUnlock()
	return daemon.collections[index]
}

func (daemon *Daemon) AddEthAlert(alert *Alert[alerts.EthereumAlert]) bool {
	return daemon.ethAlerts.Add(alert)
}

func (daemon *Daemon) RemoveEthAlert(alert *Alert[alerts.EthereumAlert]) bool {
	return daemon.ethAlerts.Remove(alert)
}

func (daemon *Daemon) AddGasAlert(alert *Alert[alerts.GasAlert]) bool {
	return daemon.gasAlerts.Add(alert)
}

func (daemon *Daemon) RemoveGasAlert(alert *Alert[alerts.GasAlert]) bool {
	return daemon.gasAlerts.Remove(alert)
}

type CollectionSaveInfo struct {
	Market nft.Marketplace          `json:"market"`
	Symbol string                   `json:"symbol"`
	Alerts []alerts.CollectionAlert `json:"alerts"`
}

func (daemon *Daemon) LoadConfig() {
	defer bench.Begin()()
	// Load ethereum alerts
	ethAlerts := config.LoadFallback[[]alerts.EthereumAlert]("ethAlerts", nil)
	for _, params := range ethAlerts {
		alert := NewEthAlert(params, daemon)
		if !daemon.AddEthAlert(alert) {
			log.Error().Println("Already in the list:", alert)
		}
	}
	// Load gas alerts
	gasAlerts := config.LoadFallback[[]alerts.GasAlert]("gasAlerts", nil)
	for _, params := range gasAlerts {
		alert := NewGasAlert(params, daemon)
		if !daemon.AddGasAlert(alert) {
			log.Error().Println("Already in the list:", alert)
		}
	}
	// Load collections
	collections := config.LoadFallback[[]CollectionSaveInfo]("collections", nil)
	for _, info := range collections {
		// log.Debug().Println("Loading user collection:", c.Marketplace, c.Slug)
		collection := NewCollection(daemon, info.Market, info.Symbol)
		if !daemon.AddCollection(collection) {
			log.Error().Println("Already in the list:", collection)
		} else {
			// Load alerts
			for _, params := range info.Alerts {
				alert := NewCollectionAlert(params, collection)
				if !collection.AddAlert(alert) {
					log.Error().Println("Already in the list:", alert)
				}
			}
		}
	}
}

// Uploads user configuration to the server
func (daemon *Daemon) SaveConfig() {
	defer bench.Begin()()
	// Save ethereum alerts
	ethAlerts := make([]alerts.EthereumAlert, daemon.ethAlerts.Len())
	daemon.ethAlerts.ForEach(func(index int, alert alerts.Alert) {
		ethAlerts[index] = alert.(*Alert[alerts.EthereumAlert]).Handle()
	})
	config.Store("ethAlerts", ethAlerts)
	// Save gas alerts
	gasAlerts := make([]alerts.GasAlert, daemon.gasAlerts.Len())
	daemon.gasAlerts.ForEach(func(index int, alert alerts.Alert) {
		gasAlerts[index] = alert.(*Alert[alerts.GasAlert]).Handle()
	})
	config.Store("gasAlerts", gasAlerts)
	// Save collections
	daemon.collectionsMutex.RLock()
	collectionInfos := make([]CollectionSaveInfo, len(daemon.collections))
	for i, collection := range daemon.collections {
		collectionInfos[i].Market = collection.Market.Load()
		collectionInfos[i].Symbol = collection.Symbol.Load()
		// Collection alerts
		colAlerts := make([]alerts.CollectionAlert, collection.alerts.Len())
		collection.alerts.ForEach(func(index int, alert alerts.Alert) {
			colAlerts[index] = alert.(*Alert[alerts.CollectionAlert]).Handle()
		})
		collectionInfos[i].Alerts = colAlerts
	}
	daemon.collectionsMutex.RUnlock()
	config.Store("collections", collectionInfos)
	// Done, now save preferences
	err := config.Save()
	if err != nil {
		log.Error().Println("Failed to save preferences:", err)
	}
}

func ReportPanic(err error) {
	log.Error().Println("Panic:", err)
	log.Error().Println("Trace:", string(debug.Stack()))
	// Generate crash report
	if isRelease() {
		DumpCrashReport(".", "nftsiren", err)
	}
}

func RecoverAndReportPanic() {
	if rec := recover(); rec != nil {
		err := fmt.Errorf("%v", rec)
		ReportPanic(err)
	}
}

func OpenLogFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_SYNC|os.O_WRONLY, 0666)
}

func DumpCrashReport(path string, name string, err error) error {
	file, ferr := OpenLogFile(path)
	if ferr != nil {
		return ferr
	}
	defer file.Close()
	// Write crash report to file
	fmt.Fprintf(file, "%s crash report %s\n", name, time.Now())
	fmt.Fprintf(file, "Error: %v\n", err)
	fmt.Fprintf(file, "Trace: \n%s\n", debug.Stack())
	return nil
}
