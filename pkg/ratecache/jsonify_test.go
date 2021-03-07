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
	testDateRangeValue := DateRangeValue{JSONDate(firstCheckIn), JSONDate(lastCheckIn), 3, 250.00}
	marshalled, _ := json.Marshal(testDateRangeValue)
	fmt.Println(string(marshalled))
	newTestDateRangeValue := DateRangeValue{}
	json.Unmarshal(marshalled, &newTestDateRangeValue)
	if testDateRangeValue != testDateRangeValue {
		t.Errorf("Value: %v, expected: %v", testDateRangeValue, newTestDateRangeValue)
	}
}
