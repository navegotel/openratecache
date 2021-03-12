package wswrite

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Settings contains config settings for ws write.
type Settings struct {
	Port                     int    `json:"port"`
	CacheDir                 string `json:"cacheDir"`
	CacheFilename            string `json:"cacheFilename"`
	Supplier                 string `json:"supplier"`
	Currency                 string `json:"currency"`
	MaxLos                   uint8  `json:"maxLos"`
	Days                     uint16 `json:"days"`
	InitialRateBlockCapacity uint16 `json:"initialRateBlockCapacity"`
}

// LoadSettings loads settings for ws write from a json file.
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

// CreateInitialSettings generates an initial confif file for ws write.
func CreateInitialSettings(filename string) error {
	s := Settings{Port: 2511, CacheDir: "/mnt/ratecache", CacheFilename: "demo.bin", Supplier: "DEMO", Currency: "EUR", MaxLos: 14, Days: 360, InitialRateBlockCapacity: 40000}
	jstr, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(filename), jstr, 0644)
	return err
}
