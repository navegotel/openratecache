package wssearch

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

type HandlerContext struct {
	Settings  Settings
	CacheFile *os.File
	Idx       *ratecache.CacheIndex
	Fhdr      *ratecache.FileHeader
}

func (context *HandlerContext) FindHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(versionInfo)
}

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

func (context *HandlerContext) AddIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
