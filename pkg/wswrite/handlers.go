package wswrite

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

// HandlerContext contains information that needs to be shared between all handlers.
type HandlerContext struct {
	Settings  Settings
	CacheFile *os.File
	Idx       *ratecache.CacheIndex
	Fhdr      *ratecache.FileHeader
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
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
	}
	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", 400)
	}
	err = ImportAriData(context, rqBody)
	if err != nil {
		http.Error(w, "Bad Request", 400)
	}
}

// HelloHandler for debugging
func (context *HandlerContext) HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello\n")
}
