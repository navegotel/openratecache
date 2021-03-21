package wssearch

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

func LoadCache(settings Settings) (*os.File, *ratecache.CacheIndex, *ratecache.FileHeader, error) {
	var f *os.File
	var err error
	var fhdr *ratecache.FileHeader
	idx := ratecache.NewCacheIndex()
	f, err = os.OpenFile(filepath.Join(settings.CacheDir, settings.CacheFilename), os.O_RDWR, 644)
	if err != nil {
		return f, idx, fhdr, err
	}
	hdrBuf := make([]byte, ratecache.FileHeaderSize)
	f.Read(hdrBuf)
	fhdr, err = ratecache.FileHeaderFromByteStr(hdrBuf)
	if err != nil {
		return f, idx, fhdr, err
	}
	err = idx.Load(fhdr, filepath.Join(settings.IndexDir, settings.CacheFilename+".idx"))
	if err != nil {
		return f, idx, fhdr, err
	}
	return f, idx, fhdr, err
}

// HelloHandler for debugging
func (context *HandlerContext) HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello\n")
}
