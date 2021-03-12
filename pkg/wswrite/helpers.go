package wswrite

import (
	"os"
	"time"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

func LoadOrCreateCache(settings Settings) (*os.File, error) {
	var f *os.File
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fhdr, _ := ratecache.NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
		ratecache.InitRateFile()
	}
}
