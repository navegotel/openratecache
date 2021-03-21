package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/navegotel/openratecache/pkg/wswrite"
)

func main() {
	clean := flag.Bool("clean", false, "starts with a clean cache")
	flag.Parse()
	configFilename := flag.Args()[0]
	settings, err := wswrite.LoadSettings(configFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Settings loaded from %v", configFilename)

	if *clean == true {
		os.Remove(filepath.Join(settings.CacheDir, settings.CacheFilename))
		os.Remove(filepath.Join(settings.IndexDir, settings.CacheFilename+".idx"))
		log.Printf("Files %v and %v removed from fs",
			filepath.Join(settings.CacheDir, settings.CacheFilename),
			filepath.Join(settings.IndexDir, settings.CacheFilename+".idx"))
	}
	cachefile, idx, err := wswrite.LoadOrCreateCache(settings)
	if err != nil {
		log.Fatal(err)
	}

	context, err := wswrite.NewHandlerContext(settings, cachefile, idx)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/version", context.VersionHandler)
	http.HandleFunc("/import", context.ImportHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), nil))
}
