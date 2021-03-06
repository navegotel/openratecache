package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

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
	mp, idx, fhdr, err := wssearch.LoadCache(settings)
	if err != nil {
		log.Fatal(err)
	}
	context := wssearch.HandlerContext{Settings: settings, Map: mp, Fhdr: fhdr, Idx: idx}

	http.HandleFunc("/list/accommodation", context.AccoListHandler)
	http.HandleFunc("/list/rooms/", context.RoomListHandler)
	http.HandleFunc("/find", context.FindHandler)
	http.HandleFunc("/addindex", context.AddIndexHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), nil))
}
