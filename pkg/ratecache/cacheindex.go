//Package index provides an index for ratecache files, tools for
//creating indexes from a cache files and saving indexes to disk

package ratecache

import (
	"errors"
	"os"

	//"fmt"
	"bytes"
	"encoding/binary"
)

// RoomOccIdx is one possible occupancy for a room or room rate.
// idx points to the rate block in the cache file.
type RoomOccIdx struct {
	Occupancy []OccupancyItem
	Total     uint8
	Idx       uint16
}

// ToByteStr returns a byte string representation of RoomOccIdx
// which can be written to the rate cache.
func (roomOccIdx *RoomOccIdx) ToByteStr() *[]byte {
	buf := make([]byte, 26, 26)
	for i, occItm := range roomOccIdx.Occupancy {
		copy(buf[i*3:(i+1)*3], *occItm.ToByteStr())
	}
	binary.BigEndian.PutUint16(buf[24:], roomOccIdx.Idx)
	return &buf
}

// AddOccItem adds one occupuncy item to the occupancy.
func (roomOccIdx *RoomOccIdx) AddOccItem(MinAge uint8, MaxAge uint8, Count uint8) error {
	if Count == 0 {
		return errors.New("Count cannot be 0")
	}
	newItemP, err := NewOccupancyItem(MinAge, MaxAge, Count)
	if err != nil {
		return err
	}
	roomOccIdx.Occupancy = append(roomOccIdx.Occupancy, *newItemP)
	roomOccIdx.Total += Count
	return nil
}

// InitIdx returns a nested map as used as cache index.
func InitIdx() map[string]map[string][]RoomOccIdx {
	m := make(map[string]map[string][]RoomOccIdx)
	return m
}

// AddRoomOccIdx adds the object to index map. Room is added if not exists,
// accommodation is added if not exists.
func AddRoomOccIdx(m map[string]map[string][]RoomOccIdx, AccoCode string, RoomRateCode string, roomOccIdx RoomOccIdx) {
	_, ok := m[AccoCode]
	if !ok {
		m[AccoCode] = make(map[string][]RoomOccIdx)
	}
	m[AccoCode][RoomRateCode] = append(m[AccoCode][RoomRateCode], roomOccIdx)
}

// SaveIdx saves index to file.
// Index format is:
// - AccoCode (length as of FileHeader object)
// - RoomCode (length as of FileHeader object)
// - 8 Occupancy items (1 MinAge, 1 MaxAge, 1 Count)
// - Index (uint16)
func SaveIdx(m map[string]map[string][]RoomOccIdx, fhdr *FileHeader, filename string) error {
	buf := make([]byte, fhdr.AccoCodeLength+fhdr.RoomRateCodeLength+FixIdxRecSize)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for accoCode, roomRateMap := range m {
		for roomRateCode, occupancies := range roomRateMap {
			for _, occupancy := range occupancies {
				copy(buf[0:], []byte(accoCode))
				copy(buf[fhdr.AccoCodeLength:], []byte(roomRateCode))
				copy(buf[fhdr.AccoCodeLength+fhdr.RoomRateCodeLength:], *occupancy.ToByteStr())
				f.Write(buf)
			}
		}
	}
	return nil
}

// ReadIdx reads a binary file from disk and builds a nested map
// as index for the rate cache.
func ReadIdx(fhdr *FileHeader, filename string) (map[string]map[string][]RoomOccIdx, error) {
	m := InitIdx()
	f, err := os.OpenFile(filename, os.O_RDONLY, 644)
	if err != nil {
		return m, err
	}
	defer f.Close()
	recordSize := int64(fhdr.AccoCodeLength + fhdr.RoomRateCodeLength + FixIdxRecSize)
	statInfo, err := f.Stat()
	if err != nil {
		return m, err
	}
	fSize := statInfo.Size()
	if fSize%recordSize != 0 {
		return m, errors.New("Incorrect file size. File may be corrupt")
	}
	buf := make([]byte, recordSize)
	var accoCode string
	var roomRateCode string
	var idx uint16
	var roomOccIdx RoomOccIdx
	recordCount := fSize / recordSize
	for i := int64(0); i < recordCount; i++ {
		f.ReadAt(buf, i*recordSize)
		accoCode = string(bytes.Trim(buf[0:fhdr.AccoCodeLength], "\x00"))
		roomRateCode = string(bytes.Trim(buf[fhdr.AccoCodeLength:fhdr.AccoCodeLength+fhdr.RoomRateCodeLength], "\x00"))
		idx = binary.BigEndian.Uint16(buf[recordSize-2 : recordSize])
		roomOccIdx = RoomOccIdx{Idx: idx}
		for j := int(fhdr.AccoCodeLength + fhdr.RoomRateCodeLength); j < int(recordSize-2); j += 3 {
			if uint8(buf[j+2]) > 0 {
				roomOccIdx.AddOccItem(uint8(buf[j]), uint8(buf[j+1]), uint8(buf[j+2]))
			}
			AddRoomOccIdx(m, accoCode, roomRateCode, roomOccIdx)
		}
	}
	return m, nil
}
