package ratecache

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// JSONDate represents date in ISO 8601 date format: YYYY-MM-DD.
type JSONDate time.Time

// MarshalJSON returns date as ISO 8601 date.
func (jd JSONDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(jd).Format("2006-01-02"))
	return []byte(stamp), nil
}

// UnmarshalJSON takes an ISO 8601 date string and returns a JSONDate object.
func (jd *JSONDate) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02", strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	*jd = JSONDate(t)
	return nil
}

// DateRangeRate represents a rate
// that is valid for various checkin dates.
type DateRangeRate struct {
	FirstCheckIn JSONDate `json:"firstCheckIn"`
	LastCheckIn  JSONDate `json:"lastCheckIn"`
	LengthOfStay uint8    `json:"lengthOfStay"`
	Rate         float32  `json:"rate"`
}

// ExplodeRate returns the exploded rates as a uint32 slice and the offset
// for the first rate in the room rate block. Check-in dates that are beyond
// the valid scope of the cache, i.e. the configured check-in dates in the
// future, will be cut off.
func (drr DateRangeRate) ExplodeRate(cacheDate time.Time, hdrSize int, days uint16, DecimalPlaces uint8) (int, []uint32) {
	lastCheckIn := time.Time(drr.LastCheckIn)
	firstCheckIn := time.Time(drr.FirstCheckIn)
	maxCheckIn := cacheDate.AddDate(0, 0, int(days))
	var rates []uint32
	if firstCheckIn.Before(cacheDate) {
		firstCheckIn = cacheDate
	}
	if lastCheckIn.Before(cacheDate) {
		return 0, rates
	}
	if firstCheckIn.After(maxCheckIn) {
		return 0, rates
	}
	if lastCheckIn.After(maxCheckIn) {
		lastCheckIn = maxCheckIn
	}
	length := int(lastCheckIn.Sub(firstCheckIn).Hours()/24 + 1)
	rates = make([]uint32, length)
	for i := 0; i < length; i++ {
		rates[i] = uint32(drr.Rate * float32(math.Pow10(int(DecimalPlaces))))
	}
	losBlockOffset := int(hdrSize) + (int(drr.LengthOfStay-1) * int(days) * 4)
	dayOffset := int(firstCheckIn.Sub(cacheDate).Hours()/24) * 4
	offset := losBlockOffset + dayOffset
	return offset, rates
}

// DateRangeAvail represents the number of available
// rooms for a range of check-in dates.
type DateRangeAvail struct {
	FirstCheckIn JSONDate `json:"firstCheckIn"`
	LastCheckIn  JSONDate `json:"lastCheckIn"`
	LengthOfStay uint8    `json:"lengthOfStay"`
	Available    uint8    `json:"available"`
}

// ExplodeAvail is similar to ExplodeRates but returns
// a slice with the availabilities instead of rates
func (dra *DateRangeAvail) ExplodeAvail(cacheDate time.Time, hdrSize int, days uint16) (int, []uint8) {
	if dra.Available > 15 {
		dra.Available = 15
	}
	lastCheckIn := time.Time(dra.LastCheckIn)
	firstCheckIn := time.Time(dra.FirstCheckIn)
	var avails []uint8
	maxCheckIn := cacheDate.AddDate(0, 0, int(days))
	// handle checkIn dates outside of cache scope
	if firstCheckIn.Before(cacheDate) {
		firstCheckIn = cacheDate
	}
	if lastCheckIn.Before(cacheDate) {
		return 0, avails
	}
	if firstCheckIn.After(maxCheckIn) {
		return 0, avails
	}
	if lastCheckIn.After(maxCheckIn) {
		lastCheckIn = maxCheckIn
	}
	length := int(lastCheckIn.Sub(firstCheckIn).Hours()/24 + 1)
	avails = make([]uint8, length)
	for i := 0; i < length; i++ {
		avails[i] = dra.Available
	}
	//calculate offset inside block
	losBlockOffset := int(hdrSize) + (int(dra.LengthOfStay-1) * int(days) * 4)
	dayOffset := int(firstCheckIn.Sub(cacheDate).Hours()/24) * 4
	offset := losBlockOffset + dayOffset
	return offset, avails
}

// DateRate represents a rate or an availability
// for one specific day.
type DateRate struct {
	CheckIn      JSONDate `json:"checkIn"`
	LengthOfStay uint8    `json:"lengthOfStay"`
	Rate         float32  `json:"rate"`
}

// RoomRates represents partially or completely the
// rates and availabilities for a room.
type RoomRates struct {
	AccoCode       string `json:"accommodationCode"`
	RoomRateCode   string `json:"roomRateCode"`
	Occupancy      []OccupancyItem
	Rates          []DateRangeRate  `json:"rates"`
	Availabilities []DateRangeAvail `json:"availabilities"`
}

func (roomRates *RoomRates) Validate() []string {
	var msg []string
	if len(roomRates.AccoCode) == 0 {
		msg = append(msg, "Missing AccoCode")
	}
	if len(roomRates.RoomRateCode) == 0 {
		msg = append(msg, "Missing RoomRateCode")
	}
	if len(roomRates.Occupancy) == 0 {
		msg = append(msg, "No Occupancy Specified")
	}
	return msg
}

// AddRate adds a DateRangeRate to RoomRates.Rates.
func (roomRates *RoomRates) AddRate(FirstCheckIn time.Time, LastCheckIn time.Time, LengthOfStay uint8, Rate float32) error {
	drr := DateRangeRate{FirstCheckIn: JSONDate(FirstCheckIn), LastCheckIn: JSONDate(LastCheckIn), LengthOfStay: LengthOfStay, Rate: Rate}
	roomRates.Rates = append(roomRates.Rates, drr)
	return nil
}

// AddAvail adds a DateRangeAvail to RoomRates.Rates.
func (roomRates *RoomRates) AddAvail(FirstCheckIn time.Time, LastCheckIn time.Time, LengthOfStay uint8, Available uint8) error {
	dra := DateRangeAvail{FirstCheckIn: JSONDate(FirstCheckIn), LastCheckIn: JSONDate(LastCheckIn), LengthOfStay: LengthOfStay, Available: Available}
	roomRates.Availabilities = append(roomRates.Availabilities, dra)
	return nil
}

//////////////////////////////////
// Request and Response formats //
//////////////////////////////////

// AccoRoomRate represents one Accommodation and all possible.
// RoomRates.
type AccoRoomRate struct {
	AccoCode     string   `json:"accoCode"`
	RoomRateCode []string `json:"roomRateCode"`
}

// SearchRq transports a set of search parameters.
type SearchRq struct {
	CheckIn        JSONDate        `json:"checkIn"`
	LengthOfStay   uint8           `json:"lengthOfStay"`
	Occupancy      []OccupancyItem `json:"occupancy"`
	Accommodations []AccoRoomRate  `json:"accommodations"`
}

//SearchRsOption groups accommodation with rate info
//for one specific combination of check-in and los.
type SearchRsOption struct {
	AccoCode     string  `json:"accoCode"`
	RoomRateCode string  `json:"roomRateCode"`
	Rate         float64 `json:"rate"`
	Availability uint8   `json:"availability"`
}

// SearchRs transports a search result.
type SearchRs struct {
	CheckIn      JSONDate         `json:"checkIn"`
	LengthOfStay uint8            `json:"lengthOfStay"`
	Options      []SearchRsOption `json:"options"`
}
