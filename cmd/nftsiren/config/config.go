package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"nftsiren/pkg/bench"
	"nftsiren/pkg/mutex"
	"nftsiren/pkg/util"
)

var state struct {
	Name   mutex.Value[string] // name of the folder, most probably application name
	Path   mutex.Value[string] // path to the preferences file
	Config *mutex.Map[string, json.RawMessage]
}

// configPath is the root of the system configuration paths
// folderName should be application name, it is also used for generating shortcut for autostart
// and filename will be the filename of the preferences file (eg "preferences.json")
// so the resulting absolute path is configPath/folderName/fileName
func Init(configPath, folderName, fileName string) error {
	configFolder := filepath.Join(configPath, folderName)
	err := os.MkdirAll(configFolder, 0666)
	if err != nil {
		return err
	}
	// all done, initialize
	state.Name.Store(folderName)
	state.Path.Store(filepath.Join(configFolder, fileName))
	state.Config = mutex.NewMap[string, json.RawMessage]()
	// load the file content if it exists
	if !util.FileExists(state.Path.Load()) {
		// No error, we just don't have a saved preference
		return nil
	}
	prefFile, err := os.Open(state.Path.Load())
	if err != nil {
		return err
	}
	defer prefFile.Close()
	err = json.NewDecoder(prefFile).Decode(&state.Config)
	if err != nil {
		return err
	}
	// fix autostart path if user enabled and executable moved
	fixAutostart()
	return nil
}

// Reports whether the key is found in the config
func Has(key string) bool {
	return state.Config.Has(key)
}

// Stores the data to the disk
func Store[T any](key string, value T) error {
	defer bench.Begin()()
	valueJson, err := json.Marshal(value)
	if err != nil {
		return err
	}
	state.Config.Store(key, valueJson)
	return Save()
}

type ErrorNotFound struct{}

func (err *ErrorNotFound) Error() string {
	return "value not found"
}

// Reports whether the returned error from load is not found error
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*ErrorNotFound)
	return ok
}

// Returned error also indicates whether it is found or not
// If the key is not found and there is no error it returns ErrorNotFound
// You can also ignore the error
func Load[T any](key string) (T, error) {
	defer bench.Begin()()
	value, ok := state.Config.Load(key)
	if !ok {
		return *new(T), &ErrorNotFound{}
	}
	var ret T
	err := json.Unmarshal(value, &ret)
	if err != nil {
		return *new(T), err
	}
	return ret, nil
}

// Tries to load given key and returns if it is found
// returns fallback if it is not found or there is an error
func LoadFallback[T any](key string, fallback T) T {
	value, err := Load[T](key)
	if err != nil {
		return fallback
	}
	return value
}

func Delete(key string) {
	state.Config.Delete(key)
}

// You do not need to call this explicitly because Store calls it everytime
func Save() error {
	defer bench.Begin()()
	prefFile, err := os.Create(state.Path.Load())
	if err != nil {
		return err
	}
	defer prefFile.Close()
	encoder := json.NewEncoder(prefFile)
	encoder.SetIndent("", "\t")
	return encoder.Encode(state.Config)
}
