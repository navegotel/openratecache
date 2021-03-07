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

// DateRangeValue represents either a rate or an availability
// That is valid for various checkin dates.
type DateRangeValue struct {
	FirstCheckIn JSONDate `json:"firstCheckIn"`
	LastCheckIng JSONDate `json:"lastCheckIn"`
	LengthOfStay uint8    `json:"lengthOfStay"`
	Value        float32  `json:"value"`
}

// DateValue represents a rate or an availability
// for one specific day.
type DateValue struct {
	CheckIn      JSONDate `json:"checkIn"`
	LengthOfStay uint8    `json:"lengthOfStay"`
	Value        float32  `json:"value"`
}

// RoomRates represents partially or completely the
// rates and availabilities for a room.
type RoomRates struct {
	AccoCode       string `json:"accommodationCode"`
	RoomRateCode   string `json:"roomRateCode"`
	Occupancy      []OccupancyItem
	Rates          []DateRangeValue `json:"rates"`
	Availabilities []DateRangeValue `json:"availabilities"`
}

// AccoRoomRate represents one Accommodation and all possible
// RoomRates.
type AccoRoomRate struct {
	AccoCode     string   `json:"accoCode`
	RoomRateCode []string `json:"roomRateCode`
}

// SearchRq transports a set of search parameters
type SearchRq struct {
	CheckIn        JSONDate        `json:"checkIn"`
	LengthOfStay   uint8           `json:"lengthOfStay"`
	Occupancy      []OccupancyItem `json:"occupancy"`
	Accommodations []AccoRoomRate  `json:"accommodations"`
}
