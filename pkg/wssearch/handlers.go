package wssearch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/exp/mmap"

	"github.com/navegotel/openratecache/pkg/ratecache"
	"github.com/navegotel/openratecache/pkg/wswrite"
)

type HandlerContext struct {
	Settings Settings
	Map      *mmap.ReaderAt
	Idx      *ratecache.CacheIndex
	Fhdr     *ratecache.FileHeader
}

func (context *HandlerContext) FindHandler(w http.ResponseWriter, r *http.Request) {
	var searchRq ratecache.SearchRq
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
	}
	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", 400)
	}
	defer r.Body.Close()
	json.Unmarshal(rqBody, &searchRq)
	validationMsgs, err := searchRq.Validate()
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", 400)
	}
	if len(validationMsgs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(validationMsgs)
		return
	}
	idxResult := context.Idx.Find(&searchRq)
	//fmt.Println(idxResult)
	searchRs := context.Find(idxResult, searchRq)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(searchRs)
}

// AccoListHandler provides an ordered list of all accommodation codes
func (context *HandlerContext) AccoListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		{
			http.Error(w, "Method Not Allowed", 405)
		}
	}
	codeList := context.Idx.GetAccoList()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codeList)
}

// RoomListHandler provides all room rate codes and the
// corresponding occupancies for one accommodation
func (context *HandlerContext) RoomListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		{
			http.Error(w, "Method Not Allowed", 405)
		}
	}
	accoCode := strings.TrimPrefix(r.URL.Path, "/list/rooms/")
	accoCode = strings.Trim(accoCode, "/")
	rooms := context.Idx.GetAccommodation(accoCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

// adds an index entry to the index based on the json data
// received in the body
func (context *HandlerContext) AddIndexHandler(w http.ResponseWriter, r *http.Request) {
	var msg wswrite.NewIdxNotification
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
	}
	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", 400)
	}
	defer r.Body.Close()
	json.Unmarshal(rqBody, &msg)
	context.Idx.AddRoomOccIdx(msg.AccoCode, msg.RoomRateCode, msg.RoomOccIdx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
