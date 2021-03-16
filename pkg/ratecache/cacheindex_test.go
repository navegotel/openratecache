package ratecache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRoomOccIdxAddOccItem(t *testing.T) {
	roomOccIdx := RoomOccIdx{}
	roomOccIdx.AddOccItem(2, 8, 1)
	roomOccIdx.AddOccItem(9, 14, 1)
	roomOccIdx.AddOccItem(15, 100, 2)
	if roomOccIdx.Total != 4 {
		t.Errorf("Value: %v, expected %v", roomOccIdx.Total, 4)
	}
}

func TestCacheIndex(t *testing.T) {
	idx := NewCacheIndex()
	roomOccIdx := RoomOccIdx{}
	roomOccIdx.AddOccItem(2, 8, 1)
	roomOccIdx.AddOccItem(9, 14, 1)
	roomOccIdx.AddOccItem(15, 100, 2)
	idx.AddRoomOccIdx("ALC001", "DBL001", roomOccIdx)
	if len(idx.m) != 1 {
		t.Errorf("Value: %d, expected 1", len(idx.m))
	}
	roomOccIdx = RoomOccIdx{}
	roomOccIdx.AddOccItem(9, 14, 1)
	roomOccIdx.AddOccItem(15, 100, 2)
	idx.AddRoomOccIdx("ALC001", "DBL001", roomOccIdx)
	if len(idx.m) != 1 {
		t.Errorf("Value: %d, expected 1", len(idx.m))
	}
	if len(idx.m["ALC001"]["DBL001"]) != 2 {
		t.Errorf("Value: %d, expected 1", len(idx.m["ALC001"]["DBL001"]))
	}
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	idxFilename := "../../test/data/test.idx"
	idx.Save(fhdr, idxFilename)
	idx2 := NewCacheIndex()
	idx2.Load(fhdr, idxFilename)
	if len(idx2.m) != 1 {
		t.Errorf("Value: %d, expected 1", len(idx.m))
	}

}

func TestLoadFromCache(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(2, 13, 1)
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	byteStr := CreateRateBlock(fhdr, rbhdr)
	filename, _ := InitRateFile(fhdr, testfolder, "../../test/data/test_load_from_cache.bin", 0)
	f, _ := os.OpenFile(filepath.Join(testfolder, filename), os.O_RDWR, 644)
	AddRateBlockToFile(f, byteStr)
	rbhdr, _ = NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	byteStr = CreateRateBlock(fhdr, rbhdr)
	AddRateBlockToFile(f, byteStr)
	f.Close()
	idx2 := NewCacheIndex()
	idx2.LoadFromCache(filepath.Join(testfolder, filename))

}

func TestGetOrCreate(t *testing.T) {
	idx := NewCacheIndex()
	roomOccIdx := RoomOccIdx{Idx: 0}
	roomOccIdx.AddOccItem(16, 100, 2)
	idx.AddRoomOccIdx("ALC01", "DBL01", roomOccIdx)
	idx.AddRoomOccIdx("ALC01", "DBL02", roomOccIdx)
	roomOccIdx = RoomOccIdx{Idx: 1}
	roomOccIdx.AddOccItem(3, 15, 1)
	roomOccIdx.AddOccItem(16, 100, 1)
	idx.AddRoomOccIdx("ALC001", "DBL01", roomOccIdx)
	roomOccIdx = RoomOccIdx{Idx: 2}
	roomOccIdx.AddOccItem(3, 15, 1)
	roomOccIdx.AddOccItem(16, 100, 2)
	idx.AddRoomOccIdx("ALC01", "DBL01", roomOccIdx)

	q := IndexQuery{AccoCode: "ALC01", RoomRateCode: "DBL01"}
	q.AddOccItem(16, 100, 2)
	idx.GetOrCreate(q)
}
