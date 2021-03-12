package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/navegotel/openratecache/pkg/ratecache"
	"github.com/navegotel/openratecache/pkg/wswrite"
)

func importHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
	}
	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", 400)
	}
	var roomRates ratecache.RoomRates
	json.Unmarshal(rqBody, &roomRates)
	fmt.Println(len(roomRates.Rates))
	fmt.Println(len(roomRates.Availabilities))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello\n")
}

func main() {
	var configFilename string
	if len(os.Args) == 2 {
		configFilename = os.Args[1]
	} else {
		configFilename = "/etc/openratecache.conf"
	}
	settings, err := wswrite.LoadSettings(configFilename)
	if err != nil {
		log.Fatal("Could not read settings file.")
	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/import", importHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", settings.Port), nil))
}
