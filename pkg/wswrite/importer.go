package wswrite

import (
	"encoding/binary"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/navegotel/openratecache/pkg/ratecache"
)

type Stats struct {
	RatesImported int
	AvailImported int
	ExecutionTime float64
}

func ImportAriData(context *HandlerContext, data []byte) (Stats, []string, error) {
	execStart := time.Now()
	var roomRates ratecache.RoomRates
	stats := Stats{}
	json.Unmarshal(data, &roomRates)
	msg := roomRates.Validate()
	if len(msg) > 0 {
		stats.ExecutionTime = time.Since(execStart).Seconds()
		return stats, msg, nil
	}
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
			stats.ExecutionTime = time.Since(execStart).Seconds()
			return stats, msg, err
		}
		roomOccIdx := ratecache.RoomOccIdx{Idx: index}
		for _, occupancyItem := range roomRates.Occupancy {
			roomOccIdx.AddOccItem(occupancyItem.MinAge, occupancyItem.MaxAge, occupancyItem.Count)
		}
		context.Idx.AddRoomOccIdx(q.AccoCode, q.RoomRateCode, roomOccIdx)
		context.Idx.Save(context.Fhdr, filepath.Join(context.Settings.IndexDir, context.Settings.CacheFilename+".idx"))
	}
	// Import data into cache
	importRates(context, &stats, index, roomRates.Rates)
	importAvail(context, &stats, index, roomRates.Availabilities)
	stats.ExecutionTime = time.Since(execStart).Seconds()
	return stats, msg, nil
}

func importRates(context *HandlerContext, stats *Stats, index uint32, dateRangeRates []ratecache.DateRangeRate) error {
	hdrSize := context.Fhdr.GetBlockHeaderSize()
	blockPos := context.Fhdr.GetRateBlockStart(index)
	uintBuf := make([]byte, 4)
	for _, dateRangeRate := range dateRangeRates {
		offset, explRange := dateRangeRate.ExplodeRate(context.Fhdr.StartDate, hdrSize, context.Fhdr.Days, context.Settings.DecimalPlaces)
		stats.RatesImported += len(explRange)
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
	hdrSize := context.Fhdr.GetBlockHeaderSize()
	blockPos := context.Fhdr.GetRateBlockStart(index)
	for _, dateRangeAvail := range dateRangeAvails {
		offset, explRange := dateRangeAvail.ExplodeAvail(context.Fhdr.StartDate, hdrSize, context.Fhdr.Days)
		stats.AvailImported += len(explRange)
		bufsize := len(explRange) * 4
		rbuf := make([]byte, bufsize)
		context.CacheFile.ReadAt(rbuf, blockPos+int64(offset))
		for i, avail := range explRange {
			rate := binary.BigEndian.Uint32(rbuf[i*4:(i+1)*4]) & ratecache.RateMask
			value := ratecache.PackRate(rate, avail)
			context.CacheFile.WriteAt(value, blockPos+int64(offset)+int64(i*4))
		}
	}
	return nil
}
