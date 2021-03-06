package ratecache

import (
    "testing"
    "time"
    "os"
    "fmt"
)

func TestRoomOccIdxAddOccItem(t *testing.T){
    roomOccIdx := RoomOccIdx{}
    roomOccIdx.AddOccItem(2, 8, 1)
    roomOccIdx.AddOccItem(9, 14, 1)
    roomOccIdx.AddOccItem(15, 100, 2)
    if roomOccIdx.Total != 4 {
        t.Errorf("Value: %v, expected %v", roomOccIdx.Total, 4)
    }
}

func TestInitIdx(t *testing.T){
    roomOccIdx := RoomOccIdx{}
    roomOccIdx.AddOccItem(2, 8, 1)
    roomOccIdx.AddOccItem(9, 14, 1)
    roomOccIdx.AddOccItem(15, 100, 2)
    
    m := InitIdx()
    _, ok := m["ALC123"]
    if !ok {
        m["ALC123"] = make(map[string][]RoomOccIdx)
        m["ALC123"]["DBLSTBR"] = append(m["ALC123"]["DBLSTBR"], roomOccIdx)
        
    } else {
        m["ALC123"]["DBLSTBR"] = append(m["ALC123"]["DBLSTBR"], roomOccIdx)
    }
    
    roomOccIdx = RoomOccIdx{}
    roomOccIdx.AddOccItem(2, 8, 2)
    roomOccIdx.AddOccItem(15, 100, 2)
    _, ok = m["ALC123"]
    if !ok {
        m["ALC123"] = make(map[string][]RoomOccIdx)
    }
    m["ALC123"]["DBLSTBR"] = append(m["ALC123"]["DBLSTBR"], roomOccIdx)
}

func TestAddRoomOccIdx(t *testing.T) {
    fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 8, 8)
    m := InitIdx()
    
    roomOccIdx := RoomOccIdx{}
    roomOccIdx.AddOccItem(2, 8, 1)
    roomOccIdx.AddOccItem(9, 14, 1)
    roomOccIdx.AddOccItem(15, 100, 2)
    roomOccIdx.Idx = 399
    
    AddRoomOccIdx(m, "ALC123", "DBLSTBR", roomOccIdx)
    
    roomOccIdx = RoomOccIdx{}
    roomOccIdx.AddOccItem(9, 14, 2)
    roomOccIdx.AddOccItem(15, 100, 2)
    roomOccIdx.Idx = 100
    
    AddRoomOccIdx(m, "ALC123", "DBLSTBR", roomOccIdx)
    
    roomOccIdx = RoomOccIdx{}
    roomOccIdx.AddOccItem(18, 100, 1)
    roomOccIdx.Idx = 200
    
    AddRoomOccIdx(m, "ALC123", "SNRBR", roomOccIdx)
    roomOccIdx.Idx = 255
    AddRoomOccIdx(m, "MUC123", "SNRBR", roomOccIdx)
    
    filename := "../../test/data/test.idx"
    SaveIdx(m, fhdr, filename)
    os.Remove(filename)
    
}

func TestReadCacheIdx(t *testing.T){
    fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 24, 24)
    m := InitIdx()
    var roomOccIdx RoomOccIdx
    for i := 0; i < 100000; i++ {
        roomOccIdx = RoomOccIdx{}
        roomOccIdx.AddOccItem(2, 14, 2)
        roomOccIdx.AddOccItem(15, 100, 2)
        roomOccIdx.Idx = uint16(i)
        AddRoomOccIdx(m, fmt.Sprintf("HTL%05d", i), "DBLSTBR", roomOccIdx)
    }
    filename := "../../test/data/test.idx"
    SaveIdx(m, fhdr, filename)
    ReadIdx(fhdr, filename)
}
