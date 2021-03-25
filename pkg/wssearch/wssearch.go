package wssearch

import (
	"path/filepath"

	"golang.org/x/exp/mmap"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

func LoadCache(settings Settings) (*mmap.ReaderAt, *ratecache.CacheIndex, *ratecache.FileHeader, error) {
	var mp *mmap.ReaderAt
	var err error
	var fhdr *ratecache.FileHeader
	idx := ratecache.NewCacheIndex()
	mp, err = mmap.Open(filepath.Join(settings.CacheDir, settings.CacheFilename))
	if err != nil {
		return mp, idx, fhdr, err
	}
	hdrBuf := make([]byte, ratecache.FileHeaderSize)
	mp.ReadAt(hdrBuf, 0)
	fhdr, err = ratecache.FileHeaderFromByteStr(hdrBuf)
	if err != nil {
		return mp, idx, fhdr, err
	}
	err = idx.Load(fhdr, filepath.Join(settings.IndexDir, settings.CacheFilename+".idx"))
	if err != nil {
		return mp, idx, fhdr, err
	}
	return mp, idx, fhdr, err
}
