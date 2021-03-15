package wswrite

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

// GetToday returns a new Time object wit just Year, Month and Day set.
func GetToday() time.Time {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return t
}

// LoadOrCreateCache is a convenience function that will return a
// a file pointer with read/write access to a rate cache file.
// If no file exists a new rate cache file will be created.
func LoadOrCreateCache(settings Settings) (*os.File, *ratecache.CacheIndex, error) {
	var f *os.File
	var err error
	idx := ratecache.NewCacheIndex()

	fhdr, err := ratecache.NewFileHeader(settings.Supplier, GetToday(), settings.Currency, settings.MaxLos, settings.Days, settings.AccoCodeLength, settings.RoomRateCodeLength)
	if err != nil {
		return f, idx, errors.New("Cannot create file header object")
	}
	_, err = os.Stat(filepath.Join(settings.CacheDir, settings.CacheFilename))
	if os.IsNotExist(err) {
		ratecache.InitRateFile(fhdr, settings.CacheDir, settings.CacheFilename, settings.InitialRateBlockCapacity)
		idx.Save(fhdr, filepath.Join(settings.IndexDir, settings.CacheFilename+".idx"))
	} else {
		err = idx.Load(fhdr, filepath.Join(settings.IndexDir, settings.CacheFilename+".idx"))
		if err != nil {
			return f, idx, err
		}
	}
	f, err = os.OpenFile(filepath.Join(settings.CacheDir, settings.CacheFilename), os.O_RDWR, 644)
	return f, idx, err
}
