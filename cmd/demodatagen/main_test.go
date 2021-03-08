package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestNewRoom(t *testing.T) {
	firstCheckIn := time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC)
	roomRatesSlice := newRoom()
	for _, roomRates := range roomRatesSlice {
		newRoomRates(&roomRates, firstCheckIn, 3, 30)
		js, _ := json.Marshal(roomRates)
		fmt.Println(string(js))
	}
}
