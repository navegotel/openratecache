package wswrite

import (
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
	// Import data into cache...
	importRates(context, &stats, 0, roomRates.Rates)
	importAvail(context, &stats, 0, roomRates.Availabilities)

	return nil
}

func importRates(context *HandlerContext, stats *Stats, index int, dateRangeRates []ratecache.DateRangeRate) error {
	//TODO get Index and create new block if necessary
	hdrSize := context.Fhdr.GetBlockHeaderSize()
	for _, dateRangeRate := range dateRangeRates {
		offset, explRange := dateRangeRate.ExplodeRate(context.Fhdr.StartDate, hdrSize, context.Fhdr.Days)
		fmt.Println(offset)
		fmt.Println(explRange)

	}
	return nil
}

func importAvail(context *HandlerContext, stats *Stats, index int, dateRangeAvails []ratecache.DateRangeAvail) error {
	for dateRangeAvail := range dateRangeAvails {
		fmt.Println(dateRangeAvail)
	}
	return nil
}
