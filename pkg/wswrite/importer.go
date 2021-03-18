package wswrite

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

type Stats struct {
	RatesDelivered int
	RatesImported  int
	AvailDelivered int
	AvailImported  int
	ExecutionTime  time.Duration
}

func ImportAriData(context *HandlerContext, data []byte) error {
	var roomRates ratecache.RoomRates
	stats := Stats{}
	json.Unmarshal(data, &roomRates)
	//Get index first
	q := ratecache.IndexQuery{AccoCode: roomRates.AccoCode, RoomRateCode: roomRates.RoomRateCode}
	for _, occupancyItem := range roomRates.Occupancy {
		q.AddOccItem(occupancyItem.MinAge, occupancyItem.MaxAge, occupancyItem.Count)
	}
	index, found := context.Idx.Get(q)
	if found == false {
		rbhdr, _ := ratecache.NewRateBlockHeader(roomRates.AccoCode, roomRates.RoomRateCode)
		byteStr := ratecache.CreateRateBlock(context.Fhdr, rbhdr)
		var err error
		index, err = ratecache.AddRateBlockToFile(context.CacheFile, byteStr)
		if err != nil {
			return err
		}
		roomOccIdx := ratecache.RoomOccIdx{Idx: index}
		for _, occupancyItem := range roomRates.Occupancy {
			roomOccIdx.AddOccItem(occupancyItem.MinAge, occupancyItem.MaxAge, occupancyItem.Count)
		}
		context.Idx.AddRoomOccIdx(q.AccoCode, q.RoomRateCode, roomOccIdx)
	}
	// Import data into cache
	importRates(context, &stats, index, roomRates.Rates)
	importAvail(context, &stats, index, roomRates.Availabilities)

	return nil
}

func importRates(context *HandlerContext, stats *Stats, index uint32, dateRangeRates []ratecache.DateRangeRate) error {
	hdrSize := context.Fhdr.GetBlockHeaderSize()
	blockPos := context.Fhdr.GetRateBlockStart(index)
	uintBuf := make([]byte, 4)
	for _, dateRangeRate := range dateRangeRates {
		offset, explRange := dateRangeRate.ExplodeRate(context.Fhdr.StartDate, hdrSize, context.Fhdr.Days)
		bufsize := len(explRange) * 4
		rbuf := make([]byte, bufsize)
		context.CacheFile.ReadAt(rbuf, blockPos+int64(offset))
		for i, rate := range explRange {
			avail := binary.BigEndian.Uint32(rbuf[i*4:(i+1)*4]) & ratecache.AvailMask
			rate := rate | avail
			binary.BigEndian.PutUint32(uintBuf, rate)
			context.CacheFile.WriteAt(uintBuf, (blockPos + int64(offset) + int64(4*i)))
		}

	}
	return nil
}

func importAvail(context *HandlerContext, stats *Stats, index uint32, dateRangeAvails []ratecache.DateRangeAvail) error {
	for _, dateRangeAvail := range dateRangeAvails {
		fmt.Println(dateRangeAvail)
	}
	return nil
}
