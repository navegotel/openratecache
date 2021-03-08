package ratecache

import (
	"fmt"
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
	fmt.Println(t)
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

// DateRangeAvail represents the number of available
// rooms for a range of check-in dates.
type DateRangeAvail struct {
	FirstCheckIn JSONDate `json:"firstCheckIn"`
	LastCheckIn  JSONDate `json:"lastCheckIn"`
	LengthOfStay uint8    `json:"lengthOfStay"`
	Available    uint8    `json:"available"`
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
	Rates          []DateRangeRate `json:"rates"`
	Availabilities []DateRangeRate `json:"availabilities"`
}

// AddRate adds a RateRangeValue to RoomRates.Rates.
func (roomRates *RoomRates) AddRate(FirstCheckIn time.Time, LastCheckIn time.Time, LengthOfStay uint8, Rate float32) error {
	drv := DateRangeRate{FirstCheckIn: JSONDate(FirstCheckIn), LastCheckIn: JSONDate(LastCheckIn), LengthOfStay: LengthOfStay, Rate: Rate}
	roomRates.Rates = append(roomRates.Rates, drv)
	return nil
}

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
