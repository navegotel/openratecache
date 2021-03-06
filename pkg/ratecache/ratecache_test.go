package ratecache

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
	//"io/ioutil"
)

const testfolder = "../../test/data"

func TestStrToTime(t *testing.T) {
	tm, _ := StrToTime("20221125")
	ctm := time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC)
	if tm != ctm {
		t.Errorf("Value: %v, expected %v", tm, ctm)
	}
}

func TestTimeToStr(t *testing.T) {
	tm := time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC)
	timeStr := TimeToStr(tm)
	timeStrC := "20221125"
	if timeStr != timeStrC {
		t.Errorf("Value: %v, expected: %v", timeStr, timeStrC)
	}
}

func TestCreateFileHeader(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	if fhdr.Supplier != "TEST" {
		t.Errorf("Value: %v, expected TEST", fhdr.Supplier)
	}
}

func TestFileHeaderToByteStr(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	byteStr := fhdr.ToByteStr()
	if len(byteStr) != 35 {
		t.Errorf("Value: %v, expected: 35", len(byteStr))
	}
}

func TestGetBlockHeaderSize(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	if fhdr.GetBlockHeaderSize() != 120 {
		t.Errorf("Value: %v, expected: 120", fhdr.GetBlockHeaderSize())
	}
}

func TestCreateRoomRateHeader(t *testing.T) {
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	if rbhdr.accoCode != "ALC123" {
		t.Errorf("Value: %v, expected: ALC123", rbhdr.accoCode)
	}
	if rbhdr.roomRateCode != "DBLSTDBRBAR" {
		t.Errorf("Value: %v, expected: DBLSTDBRBAR", rbhdr.roomRateCode)
	}
}

func TestAddOccupancyItem(t *testing.T) {
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(2, 13, 1)
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	if len(rbhdr.occupancy) != 3 {
		t.Errorf("Value: %v, expected: 3", len(rbhdr.occupancy))
	}
	if rbhdr.occupancy[0].MaxAge != 13 {
		t.Errorf("Value: %v, expected: 13", rbhdr.occupancy[0].MaxAge)
	}
}

func TestRateBlockHeaderToByteStr(t *testing.T) {
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(2, 13, 1)
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	byteStr := rbhdr.ToByteStr(32, 64)
	if len(byteStr) != 120 {
		t.Errorf("Value: %v, expected: 120", len(byteStr))
	}
}

func TestRateBlockHeaderFromByteStr(t *testing.T) {
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(2, 13, 1)
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	byteStr := rbhdr.ToByteStr(32, 64)
	rbhdr2, _ := RateBlockHeaderFromByteStr(byteStr, 32, 64)
	if rbhdr2.accoCode != "ALC123" {
		fmt.Println([]byte(rbhdr2.accoCode))
		t.Errorf("Value: %v, expected: ALC123", rbhdr.accoCode)
	}
}

func TestFileHeaderFromByteStr(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	byteStr := fhdr.ToByteStr()
	//err := ioutil.WriteFile("cachefile", byteStr, 0644)
	fhdr2, err := FileHeaderFromByteStr(byteStr)
	if err != nil {
		t.Error(err)
	} else {
		if fhdr2.Currency != "EUR" {
			t.Errorf("Value: %v, expected: EUR", fhdr2.Currency)
		}
	}
}

func TestCreateRateBlock(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(2, 13, 1)
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	byteStr := CreateRateBlock(fhdr, rbhdr)
	expectedLen := 32 + 64 + 24 + 14*400*4
	if len(byteStr) != expectedLen {
		t.Errorf("Value: %v, expected: %v", len(byteStr), expectedLen)
	}
	if fhdr.GetRateBlockSize() != expectedLen {
		t.Errorf("Value: %v, expected: %v", fhdr.GetRateBlockSize(), expectedLen)
	}
}

func TestInitRateFile(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	filename, err := InitRateFile(fhdr, testfolder, 100)
	if err != nil {
		t.Error(err)
	}
	if filename == "" {
		t.Error("Got empty string as filename")
	}
	os.Remove(filename)
}

func TestAppendRateBlockToFile(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(2, 13, 1)
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	byteStr := CreateRateBlock(fhdr, rbhdr)
	filename, _ := InitRateFile(fhdr, testfolder, 0)
	f, _ := os.OpenFile(filepath.Join(testfolder, filename), os.O_RDWR, 644)
	defer f.Close()
	defer os.Remove(filename)
	idx, _ := AddRateBlockToFile(f, byteStr)
	if idx != 0 {
		t.Errorf("Value: %v, expected: %v", idx, 0)
	}
	idx, _ = AddRateBlockToFile(f, byteStr)
	if idx != 1 {
		t.Errorf("Value: %v, expected: %v", idx, 1)
	}
}

func TestGetRateBlockStart(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(2, 13, 1)
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	byteStr := CreateRateBlock(fhdr, rbhdr)
	filename, _ := InitRateFile(fhdr, testfolder, 0)
	f, _ := os.OpenFile(filename, os.O_RDWR, 644)
	defer f.Close()
	defer os.Remove(filename)
	AddRateBlockToFile(f, byteStr)
	if fhdr.GetRateBlockStart(0) != FileHeaderSize {
		t.Errorf("Value: %v, expected: %v", fhdr.GetRateBlockStart(0), FileHeaderSize)
	}
	if fhdr.GetRateBlockStart(1) != 22555 {
		t.Errorf("Value: %v, expected: %v", fhdr.GetRateBlockStart(1), 22555)
	}
}

func TestPackRate(t *testing.T) {
	packedRate := PackRate(45000, 12)
	rate, avail := UnpackRate(packedRate)
	if rate != 45000 {
		t.Errorf("Value: %v, expected: %v", rate, 45000)
	}
	if avail != 12 {
		t.Errorf("Value: %v, expected: %v", avail, 12)
	}
}

/*
func TestGetRatePos(t *testing.T){
    fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
    fhdr.rateBlockCount += 1
    checkIn := time.Date(2022, time.December, 10, 0, 0, 0, 0, time.UTC)
    los:= uint8(4)
    ratePos, err := fhdr.GetRatePos(1, checkIn, los )
    if err == nil {
        t.Error("Expected error because idx is too big")
    }

    fhdr.rateBlockCount += 1
    ratePos, err = fhdr.GetRatePos(1, checkIn, los )
    if ratePos != 27895 {
        t.Errorf("Value: %v, expected: %v", ratePos, 27895)
    }
}
*/

func TestSetGetRateInfo(t *testing.T) {
	fhdr, _ := NewFileHeader("TEST", time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC), "EUR", 14, 400, 32, 64)
	rbhdr, _ := NewRateBlockHeader("ALC123", "DBLSTDBRBAR")
	rbhdr.AddOccupancyItem(2, 13, 1)
	rbhdr.AddOccupancyItem(14, 17, 1)
	rbhdr.AddOccupancyItem(18, 100, 2)
	byteStr := CreateRateBlock(fhdr, rbhdr)
	filename, _ := InitRateFile(fhdr, testfolder, 0)
	f, _ := os.OpenFile(filepath.Join(testfolder, filename), os.O_RDWR, 644)
	defer f.Close()
	defer os.Remove(filename)
	AddRateBlockToFile(f, byteStr)
	fhdr.RateBlockCount++
	checkIn := time.Date(2022, time.December, 10, 0, 0, 0, 0, time.UTC)
	err := fhdr.SetRateInfo(f, 0, checkIn, 4, 25500, 4)
	if err != nil {
		t.Error(err)
	}
	checkIn = time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC)
	err = fhdr.SetRateInfo(f, 0, checkIn, 1, 16581375, 15)
	if err != nil {
		t.Error(err)
	}
	checkIn = time.Date(2022, time.November, 26, 0, 0, 0, 0, time.UTC)
	err = fhdr.SetRateInfo(f, 0, checkIn, 1, 16581375, 15)
	if err != nil {
		t.Error(err)
	}
	checkIn = time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC)
	err = fhdr.SetRateInfo(f, 0, checkIn, 2, 16581375, 15)
	if err != nil {
		t.Error(err)
	}
	rate, avail, _ := fhdr.GetRateInfo(f, 0, checkIn, 2)
	if rate != 16581375 {
		t.Errorf("Value: %v, expected: %v", rate, 16581375)
	}
	if avail != 15 {
		t.Errorf("Value: %v, expected: %v", avail, 15)
	}
}
