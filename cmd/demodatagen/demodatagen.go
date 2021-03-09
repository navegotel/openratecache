// Package demodatagen generates a json file with fake ARI data
// which can be imported into the rate cache.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

var airports = []string{"FRA", "MUC", "CGN", "DUS", "STR", "HHN", "LEJ", "BER", "BRE", "NRN",
	"ABC", "ALC", "LEI", "BCN", "OVD", "BIO", "GRO", "MAD", "AGP", "PNA", "SDR", "SVQ", "VLC", "VGO", "ZAZ", "PMI", "TFS", "TFN", "MAH",
	"CDG", "ORY", "NCE", "LYS", "MRS", "TLS", "BOD", "NTE", "BVA", "RUN", "LIL", "MPL", "AJA", "BIA", "SXB", "BES", "BIQ", "FSC", "TLN", "EGC",
	"PSR", "AOT", "BRI", "BDS", "FOG", "CRV", "SUF", "REG", "NAP", "QSR", "BLQ", "FRL", "PMF", "RAN", "RMI", "RMI", "TRS", "FCO", "CIA", "ALL", "BGY", "VBS", "AOI", "CUF", "TRN", "AHO", "CAG", "OLB", "TTB", "CTA", "LMP", "PMO",
	"LIS", "OPO", "FAO", "FNC", "PDL",
	"ATH", "HER", "SKG", "RHO", "CFU", "CHQ", "KGS", "JTR", "ZTH", "JMK",
	"VIE", "SZG", "INN", "GRZ", "LNZ", "KLU",
	"BRU", "CRL", "OST", "ANR", "LGG",
	"AMS", "EIN", "RTM", "MST", "GRQ",
	"ZRH", "GVA", "BSL", "BRN",
	"ARN", "GOT", "BMA", "NYO", "MMX", "VST",
	"CPH", "BLL", "AAL", "AAR", "FAE",
	"OSL", "BGO", "TRD", "SVG", "TOS", "TRF", "BOO", "AES", "KRS", "HAU",
	"BUD", "DEB",
	"IST", "AYT", "SAW", "ESB", "ADB", "ADA", "DLM", "BJV", "TZX",
	"CAI", "HRG", "SSH", "HBE"}

var roomTypes = []string{"SGLSTAO", "SGLSTBR", "DBLBDAO", "DBLBDBR", "DBLSTAO", "DBLSTBR", "DBLSTHB",
	"DBLDXAO", "DBLDXBR", "DBLDXHB", "DBLFRAO", "DBLFRHB", "DBLFRAO", "JUSUIAO", "JUSUIBR", "JUSUIHB"}

func newAccoCode(i int) string {
	airport := airports[rand.Intn(len(airports))]
	code := fmt.Sprintf("%v%05d", airport, i)
	return code
}

func newRoomRateCode() string {
	roomType := roomTypes[rand.Intn(len(roomTypes))]
	roomRateCode := fmt.Sprintf("%v%3d", roomType, rand.Intn(999))
	return roomRateCode
}

func newRoom(accoCode string) []ratecache.RoomRates {
	var roomRates ratecache.RoomRates
	roomRateCode := newRoomRateCode()
	roomRatesSlice := make([]ratecache.RoomRates, 0)
	switch prefix := roomRateCode[:3]; prefix {
	case "SGL":
		roomRates.AccoCode = accoCode
		roomRates.RoomRateCode = roomRateCode
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 18, MaxAge: 100, Count: 1})
		roomRatesSlice = append(roomRatesSlice, roomRates)
	case "DBL":
		roomRates.AccoCode = accoCode
		roomRates.RoomRateCode = roomRateCode
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 3, MaxAge: 12, Count: 1})
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 13, MaxAge: 16, Count: 1})
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 17, MaxAge: 100, Count: 2})
		roomRatesSlice = append(roomRatesSlice, roomRates)
		roomRates.Occupancy = roomRates.Occupancy[:0]
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 13, MaxAge: 16, Count: 1})
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 17, MaxAge: 100, Count: 2})
		roomRatesSlice = append(roomRatesSlice, roomRates)
		roomRates.Occupancy = roomRates.Occupancy[:0]
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 3, MaxAge: 12, Count: 2})
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 17, MaxAge: 100, Count: 2})
		roomRatesSlice = append(roomRatesSlice, roomRates)
		roomRates.Occupancy = roomRates.Occupancy[:0]
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 17, MaxAge: 100, Count: 2})
		roomRatesSlice = append(roomRatesSlice, roomRates)
	case "JUS":
		roomRates.AccoCode = accoCode
		roomRates.RoomRateCode = roomRateCode
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 13, MaxAge: 16, Count: 1})
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 17, MaxAge: 100, Count: 3})
		roomRatesSlice = append(roomRatesSlice, roomRates)
		roomRates.Occupancy = roomRates.Occupancy[:0]
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 13, MaxAge: 16, Count: 2})
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 17, MaxAge: 100, Count: 2})
		roomRatesSlice = append(roomRatesSlice, roomRates)
		roomRates.Occupancy = roomRates.Occupancy[:0]
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 17, MaxAge: 100, Count: 3})
		roomRatesSlice = append(roomRatesSlice, roomRates)
		roomRates.Occupancy = roomRates.Occupancy[:0]
		roomRates.Occupancy = append(roomRates.Occupancy, ratecache.OccupancyItem{MinAge: 17, MaxAge: 100, Count: 4})
		roomRatesSlice = append(roomRatesSlice, roomRates)
	}
	return roomRatesSlice
}

func newRoomRates(roomRates *ratecache.RoomRates, startDate time.Time, maxLos int, days int) {
	i := 0
	dateSpan := 0
	firstCheckIn := time.Time(startDate)
	lastCheckIn := firstCheckIn.AddDate(0, 0, rand.Intn(12))
	for los := 1; los <= maxLos; los++ {
		i = 0
		firstCheckIn = time.Time(startDate)
		lastCheckIn = firstCheckIn.AddDate(0, 0, rand.Intn(12))
		for i <= days {
			dateSpan = rand.Intn(12)
			roomRates.AddRate(firstCheckIn, lastCheckIn, uint8(los), float32(2500+rand.Intn(6500))/100*float32(los))
			firstCheckIn = lastCheckIn.AddDate(0, 0, 1)
			lastCheckIn = firstCheckIn.AddDate(0, 0, dateSpan)
			i += dateSpan
		}

	}
}

func newRoomAvail(roomRates *ratecache.RoomRates, startDate time.Time, maxLos int, days int) {
	i := 0
	dateSpan := 0
	firstCheckIn := time.Time(startDate)
	lastCheckIn := firstCheckIn.AddDate(0, 0, rand.Intn(12))
	for los := 1; los <= maxLos; los++ {
		i = 0
		firstCheckIn = time.Time(startDate)
		lastCheckIn = firstCheckIn.AddDate(0, 0, rand.Intn(12))
		for i <= days {
			dateSpan = rand.Intn(12)
			roomRates.AddAvail(firstCheckIn, lastCheckIn, uint8(los), uint8(rand.Intn(10)))
			firstCheckIn = lastCheckIn.AddDate(0, 0, 1)
			lastCheckIn = firstCheckIn.AddDate(0, 0, dateSpan)
			i += dateSpan
		}

	}
}

func generateDemoData(accoCount int, folder string, maxLos int, days int) {
	var accoCode string
	var roomRates []ratecache.RoomRates
	firstCheckIn := time.Now()
	dirExists, _ := ratecache.DirExists(folder)
	if dirExists == false {
		os.MkdirAll(folder, os.ModePerm)
	} else {
		dirIsEmpty, _ := ratecache.DirIsEmpty(folder)
		if dirIsEmpty == false {
			log.Fatal("Specified output folder is not empty. Exiting.")
		}

	}
	for i := 0; i <= accoCount; i++ {
		accoCode = newAccoCode(i)
		for j := 0; j < 5; j++ {
			roomRates = newRoom(accoCode)
			for k, roomRate := range roomRates {
				newRoomRates(&roomRate, firstCheckIn, maxLos, days)
				newRoomAvail(&roomRate, firstCheckIn, maxLos, days)
				jstr, err := json.Marshal(roomRate)
				if err != nil {
					log.Fatal("Could not generate json")
				}
				filename := fmt.Sprintf("%v%05d%d%d.json", accoCode, i, j, k)
				err = ioutil.WriteFile(filepath.Join(folder, filename), jstr, 0644)
				if err != nil {
					log.Fatal("could not write to file")
				}
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	accoCount := flag.Int("n", 1000, "Number of accommodations to be generated")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	defaultFolder := filepath.Join(usr.HomeDir, "demodata")
	folder := flag.String("o", defaultFolder, "The folder to which the demo data is going to be saved")
	maxLos := flag.Int("l", 14, "MaxLos, the maximum length of stay for which rates are stored in the cache")
	days := flag.Int("d", 360, "Days the number of check-in dates in the future for which rates are stored in the cache")
	flag.Parse()
	generateDemoData(*accoCount, *folder, *maxLos, *days)
}
