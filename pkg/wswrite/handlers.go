package wswrite

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

// HandlerContext contains information that needs to be shared between all handlers.
type HandlerContext struct {
	Settings  Settings
	CacheFile *os.File
	Idx       *ratecache.CacheIndex
	Fhdr      *ratecache.FileHeader
}

type ImportInfo struct {
	Errors []string `json:"errors"`
	Stats  Stats    `json:"stats"`
}

// NewHandlerContext creates a new handler context
func NewHandlerContext(settings Settings, cacheFile *os.File, idx *ratecache.CacheIndex) (*HandlerContext, error) {
	context := HandlerContext{Settings: settings, CacheFile: cacheFile, Idx: idx}
	buf := make([]byte, ratecache.FileHeaderSize)
	cacheFile.Read(buf)
	fhdr, err := ratecache.FileHeaderFromByteStr(buf)
	if err != nil {
		return &context, err
	}
	context.Fhdr = fhdr
	return &context, nil
}

// ImportHandler imports data into the rate cache.
func (context *HandlerContext) ImportHandler(w http.ResponseWriter, r *http.Request) {
	var importInfo ImportInfo
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
	}
	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", 400)
	}
	importInfo.Stats, importInfo.Errors, err = ImportAriData(context, rqBody)
	if err != nil {
		http.Error(w, "Bad Request", 400)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(importInfo)

}

type VersionInfo struct {
	Release            string    `json:"release"`
	FormatVersion      byte      `json:"formatVersion"`
	CacheDate          time.Time `json:"cacheDate"`
	AccommodationCount int       `json:"accommodationCount"`
	RateBlockCount     uint32    `json:"rateBlockCount"`
	RateCount          uint64    `json:"rateCount"`
}

// VersionHandler for basic cache information
func (context *HandlerContext) VersionHandler(w http.ResponseWriter, r *http.Request) {
	versionInfo := VersionInfo{Release: ratecache.Release,
		FormatVersion:      ratecache.Version,
		CacheDate:          context.Fhdr.StartDate,
		AccommodationCount: context.Idx.GetAccoCount(),
		RateBlockCount:     context.Fhdr.RateBlockCount,
		RateCount:          uint64(context.Fhdr.Days) * uint64(context.Fhdr.MaxLos) * uint64(context.Fhdr.RateBlockCount),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(versionInfo)
}
