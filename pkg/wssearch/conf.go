package wssearch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Settings struct {
	Port          int    `json:"port"`
	CacheDir      string `json:"cacheDir"`
	IndexDir      string `json:"indexDir"`
	CacheFilename string `json:"cacheFilename"`
	DecimalPlaces uint8  `json:"decimalPlaces"`
}

func LoadSettings(filename string) (Settings, error) {
	s := Settings{}
	f, err := os.Open(filename)
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(buf, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// CreateInitialSettings creates a settings file with everything
// necessary for starting the web service. Everything else is taken
// from the cache header file.
func CreateInitialSettings(filename string) error {
	s := Settings{
		Port:          2507,
		CacheDir:      "/mnt/ratecache",
		IndexDir:      "/opt/ratecache",
		CacheFilename: "demo.bin",
		DecimalPlaces: 2,
	}
	jstr, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(filename), jstr, 0644)
	return err
}
