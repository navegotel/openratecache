package main

import (
	"flag"
	"log"

	"github.com/navegotel/openratecache/pkg/wssearch"
)

func main() {
	flag.Parse()
	configFilename := flag.Args()[0]
	settings, err := wssearch.LoadSettings(configFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Settings loaded from %v", configFilename)
	cachefile, idx, fhdr, err := wssearch.LoadCache(settings)
	if err != nil {
		log.Fatal(err)
	}
	context := wssearch.HandlerContext{Settings: settings, CacheFile: cachefile, Fhdr: fhdr, Idx: idx}
}
