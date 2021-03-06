package ratecache

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTypeDateRangeValue(t *testing.T) {
	firstCheckIn := time.Date(2022, time.November, 15, 0, 0, 0, 0, time.UTC)
	lastCheckIn := time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC)
	testDateRangeRate := DateRangeRate{JSONDate(firstCheckIn), JSONDate(lastCheckIn), 3, 250.00}
	marshalled, _ := json.Marshal(testDateRangeRate)
	newTestDateRangeRate := DateRangeRate{}
	json.Unmarshal(marshalled, &newTestDateRangeRate)
	if testDateRangeRate != testDateRangeRate {
		t.Errorf("Value: %v, expected: %v", testDateRangeRate, newTestDateRangeRate)
	}
}

func TestExplodeRate(t *testing.T) {
	cacheDate := time.Date(2022, time.November, 1, 0, 0, 0, 0, time.UTC)
	firstCheckIn := time.Date(2022, time.November, 2, 0, 0, 0, 0, time.UTC)
	lastCheckIn := time.Date(2022, time.November, 3, 0, 0, 0, 0, time.UTC)
	headerSize := 50
	days := uint16(30)
	dateRangeRate := DateRangeRate{
		FirstCheckIn: JSONDate(firstCheckIn),
		LastCheckIn:  JSONDate(lastCheckIn),
		LengthOfStay: 3,
		Rate:         250.00,
	}
	offset, b := dateRangeRate.ExplodeRate(cacheDate, headerSize, days, 2)
	if offset != 294 {
		t.Errorf("Value %d, expected value 294", offset)
	}
	if len(b) != 2 {
		t.Errorf("Value %d, expected: 2", len(b))
	}
	firstCheckIn = time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC)
	lastCheckIn = time.Date(2022, time.December, 10, 0, 0, 0, 0, time.UTC)
	dateRangeRate.FirstCheckIn = JSONDate(firstCheckIn)
	dateRangeRate.LastCheckIn = JSONDate(lastCheckIn)
	dateRangeRate.LengthOfStay = 1
	offset, b = dateRangeRate.ExplodeRate(cacheDate, headerSize, days, 2)
	if offset != 146 {
		t.Errorf("Value %d, expected value: 24", offset)
	}
	/* need to check this
	if len(b) != 6 {
		t.Errorf("Value %d, expected: 6", len(b))
	}
	*/
}

/*
func TestNegativeSub(t *testing.T) {
	firstCheckIn := time.Date(2022, time.November, 10, 0, 0, 0, 0, time.UTC)
	lastCheckIn := time.Date(2022, time.November, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println(lastCheckIn.Sub(firstCheckIn).Hours() / 24)

	testDate := time.Date(2022, time.November, 1, 0, 0, 0, 0, time.UTC)
	daylapse := 3
	testDate2 := testDate.Add(time.Hour * time.Duration(daylapse) * 24)
	fmt.Println(testDate2)
}
*/
/*
func TestSearchRq(t *testing.T) {
	arr := AccoRoomRate{AccoCode: "ADA00009", RoomRateCodes: make([]string, 0)}
	searchRq := SearchRq{CheckIn: JSONDate(time.Date(2022, time.November, 10, 0, 0, 0, 0, time.UTC)), LengthOfStay: 5, Occupancy: []uint8{12, 25, 25}}
	searchRq.Accommodations = append(searchRq.Accommodations, arr)
	msg, _ := json.Marshal(searchRq)
	fmt.Print(string(msg))
	searchRq.Validate()
}
*/

func TestAges(t *testing.T) {
	ages := []uint8{8, 23, 32}
	jsonStr, _ := json.Marshal(Ages(ages))
	fmt.Println(string(jsonStr))
	var newAges Ages
	json.Unmarshal(jsonStr, &newAges)
	fmt.Println(newAges)
}
