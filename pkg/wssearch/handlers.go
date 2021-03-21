package wssearch

import (
	"os"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

type HandlerContext struct {
	Settings  Settings
	CacheFile *os.File
	Idx       *ratecache.CacheIndex
	Fhdr      *ratecache.FileHeader
}
