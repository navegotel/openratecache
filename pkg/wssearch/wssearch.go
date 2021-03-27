package wssearch

import (
	"log"
	"math"
	"path/filepath"
	"time"

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

func (context HandlerContext) Find(idxResults []ratecache.IdxResult, searchRq ratecache.SearchRq) ratecache.SearchRs {
	searchRs := ratecache.SearchRs{CheckIn: searchRq.CheckIn, LengthOfStay: searchRq.LengthOfStay}
	for _, idxResult := range idxResults {
		accoOption := ratecache.SearchRsAccoOption{AccoCode: idxResult.AccoCode}
		for _, room := range idxResult.Rooms {
			roomOption := ratecache.SearchRsRoomOption{RoomRateCode: room.RoomRateCode}
			rate, avail, err := context.Fhdr.GetRateInfoFromMap(*context.Map, room.Index, time.Time(searchRq.CheckIn), searchRq.LengthOfStay)
			if err != nil {
				log.Print(err)
			}
			if avail > 0 && rate > 0 {
				roomOption.Rate = float64(rate) / math.Pow10(int(context.Settings.DecimalPlaces))
				roomOption.Availability = avail
				accoOption.Rooms = append(accoOption.Rooms, roomOption)
			}
		}
		searchRs.Options = append(searchRs.Options, accoOption)
	}
	return searchRs
}
