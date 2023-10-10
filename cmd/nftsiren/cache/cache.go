package cache

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"nftsiren/pkg/bench"
	"nftsiren/pkg/images"
	"nftsiren/pkg/mutex"
	"nftsiren/pkg/util"
)

// 7 days
const fileValidityTime = time.Hour * 24 * 7

var cacheDir mutex.Value[string]

func Init(dir string) error {
	// Create cacheDir if not exist
	err := os.MkdirAll(dir, 0666)
	if err != nil {
		return err
	}
	cacheDir.Store(dir)
	return nil
}

func Save(uri string, data []byte) error {
	defer bench.Begin()()
	if !initialized() {
		panic("cache: not initialized")
	}
	return os.WriteFile(getUriPath(uri), data, 0666)
}

// There will be no error if the uri can not be found
// so check both of them
func Load(uri string) ([]byte, error) {
	defer bench.Begin()()
	if !initialized() {
		panic("cache: not initialized")
	}
	// find file name
	path := getUriPath(uri)
	if util.FileExists(path) {
		// Stat the file and check it's creation date
		// we don't allow files created earlier than fileValidityTime
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if time.Since(info.ModTime()) < fileValidityTime {
			return os.ReadFile(path)
		}
	}
	// not found
	return nil, nil
}

func SaveImage(uri string, data image.Image) error {
	defer bench.Begin()()
	var b bytes.Buffer
	// jpeg is generally faster than png
	err := jpeg.Encode(&b, data, nil)
	if err != nil {
		return err
	}
	return Save(uri, b.Bytes())
}

func LoadImage(uri string) (*image.RGBA, error) {
	defer bench.Begin()()
	data, err := Load(uri)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return images.Parse(data)
}

func getUriPath(uri string) string {
	hash := sha256.Sum256([]byte(uri))
	hashHex := hex.EncodeToString(hash[:])
	return filepath.Join(cacheDir.Load(), hashHex)
}

func initialized() bool {
	return cacheDir.Load() != ""
}
