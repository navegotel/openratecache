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
	fmt.Println(string(marshalled))
	newTestDateRangeRate := DateRangeRate{}
	json.Unmarshal(marshalled, &newTestDateRangeRate)
	if testDateRangeRate != testDateRangeRate {
		t.Errorf("Value: %v, expected: %v", testDateRangeRate, newTestDateRangeRate)
	}
}
