package ratecache

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"sort"
	"sync"
)

type IndexQuery struct {
	AccoCode     string
	RoomRateCode string
	Occupancy    []OccupancyItem
	OccTotal     uint8
}

// AddOccItem adds an Occupancy item, consisting of MinAge, MaxAge and Count to the
// requested occupancy. This method should always be used as setter instead of
// directly appending to the Occupancy attribute.
func (indexQuery *IndexQuery) AddOccItem(MinAge uint8, MaxAge uint8, Count uint8) error {
	if Count == 0 {
		return errors.New("Count cannot be 0")
	}
	newItemP, err := NewOccupancyItem(MinAge, MaxAge, Count)
	if err != nil {
		return err
	}
	indexQuery.Occupancy = append(indexQuery.Occupancy, *newItemP)
	indexQuery.OccTotal += Count
	return nil
}

// RoomOccIdx is one possible occupancy for a room or room rate.
// idx points to the rate block in the cache file.
type RoomOccIdx struct {
	Occupancy []OccupancyItem
	Total     uint8
	Idx       uint32
}

// Match matches a specific group of guests, represented by
// a slice of ages against a room occupancy and returns true
// if the group matches the occupancy requirements.
func (roomOccIdx *RoomOccIdx) Match(guests []uint8) bool {
	if len(guests) != int(roomOccIdx.Total) {
		return false
	}
	counters := make([]uint8, len(roomOccIdx.Occupancy))
	for i, occItem := range roomOccIdx.Occupancy {
		counters[i] = occItem.Count
	}
	for _, guest := range guests {
		fit := false
		for i, occItem := range roomOccIdx.Occupancy {
			if counters[i] > 0 && guest >= occItem.MinAge && guest <= occItem.MaxAge {
				fit = true
				counters[i] -= 1
				break
			}
		}
		if fit == false {
			return false
		}
	}
	for _, count := range counters {
		if count > 0 {
			return false
		}
	}

	return true
}

// ToByteStr returns a byte string representation of RoomOccIdx
// which can be written to the rate cache.
func (roomOccIdx *RoomOccIdx) ToByteStr() *[]byte {
	buf := make([]byte, 28, 28)
	for i, occItm := range roomOccIdx.Occupancy {
		copy(buf[i*3:(i+1)*3], *occItm.ToByteStr())
	}
	binary.BigEndian.PutUint32(buf[24:], roomOccIdx.Idx)
	return &buf
}

// AppendToIdxFile appends a new index entry to the index file
// without having to re-write the whole index on disk
func (roomOccIdx *RoomOccIdx) AppendToIdxFile(fhdr FileHeader, filename string, accoCode string, roomRateCode string) error {
	blockSize := fhdr.AccoCodeLength + fhdr.RoomRateCodeLength + FixIdxRecSize
	buf := make([]byte, blockSize)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	copy(buf[0:], []byte(accoCode))
	copy(buf[fhdr.AccoCodeLength:], []byte(roomRateCode))
	copy(buf[fhdr.AccoCodeLength+fhdr.RoomRateCodeLength:], *roomOccIdx.ToByteStr())
	f.Write(buf)
	return nil
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

// CacheIndex contains a map with the nested map
// and slice structure for the cache index,
// protected by a mutex.
type CacheIndex struct {
	m map[string]map[string][]RoomOccIdx
	sync.RWMutex
}

// NewCacheIndex returns a pointer to a new CacheIndex
// instance. Returning a pointer is necessary because
// returning a copy of a mutex is not safe.
func NewCacheIndex() *CacheIndex {
	idx := CacheIndex{}
	idx.m = make(map[string]map[string][]RoomOccIdx)
	return &idx
}

// RoomIdx is part of IdxResult and contains a
// room rate code and the corresponding index
// of the room rate block
type RoomIdx struct {
	RoomRateCode string
	Index        uint32
}

// A slice of IdexResult is returned by "Find" as
// a result to the search. This slice is used as
// input to get the actual rate information from
// the rate cache.
type IdxResult struct {
	AccoCode string
	Rooms    []RoomIdx
}

func (idx *CacheIndex) Find(searchRq *SearchRq) []IdxResult {
	var idxResults []IdxResult
	for _, accommodation := range searchRq.Accommodations {
		if len(accommodation.RoomRateCodes) == 0 {
			idx.Lock()
			for key, _ := range idx.m[accommodation.AccoCode] {
				accommodation.RoomRateCodes = append(accommodation.RoomRateCodes, key)
			}
			idx.Unlock()
		}
		idxResult := IdxResult{AccoCode: accommodation.AccoCode}
		for _, room := range accommodation.RoomRateCodes {
			idx.Lock()
			for _, roomOccIdx := range idx.m[accommodation.AccoCode][room] {
				if roomOccIdx.Match(searchRq.Occupancy) {
					roomIdx := RoomIdx{RoomRateCode: room, Index: roomOccIdx.Idx}
					idxResult.Rooms = append(idxResult.Rooms, roomIdx)
					break
				}
			}
			idx.Unlock()
		}
		if len(idxResult.Rooms) > 0 {
			idxResults = append(idxResults, idxResult)
		}
	}
	return idxResults
}

// GetAccoCount returns the number of accommodations in the idx.
func (idx *CacheIndex) GetAccoCount() int {
	return len(idx.m)
}

// Get AccoList returns a slice of all accommodations in the idx.
func (idx *CacheIndex) GetAccoList() []string {
	idx.Lock()
	l := make([]string, len(idx.m))
	i := 0
	for accoCode := range idx.m {
		l[i] = accoCode
		i++
	}
	idx.Unlock()
	sort.Strings(l)
	return l
}

func (idx *CacheIndex) GetAccommodation(AccoCode string) map[string][][]OccupancyItem {
	rooms := make(map[string][][]OccupancyItem)
	idx.Lock()
	entries := idx.m[AccoCode]
	for code, occupancies := range entries {
		for _, occupancy := range occupancies {
			var o []OccupancyItem
			o = append(o, occupancy.Occupancy...)
			rooms[code] = append(rooms[code], o)
		}
	}
	idx.Unlock()
	return rooms
}

//AddRoomOccIdx adds a new RoomOccIdx to the index.
func (idx *CacheIndex) AddRoomOccIdx(accoCode string, roomRateCode string, roomOccIdx RoomOccIdx) error {
	idx.Lock()
	_, ok := idx.m[accoCode]
	if !ok {
		idx.m[accoCode] = make(map[string][]RoomOccIdx)
	}
	idx.m[accoCode][roomRateCode] = append(idx.m[accoCode][roomRateCode], roomOccIdx)
	idx.Unlock()
	return nil
}

// Save saves the whole index to a file.
// Index format is:
// - AccoCode (length as of FileHeader object)
// - RoomCode (length as of FileHeader object)
// - 8 Occupancy items (1 MinAge, 1 MaxAge, 1 Count)
// - Index (uint16)
func (idx *CacheIndex) Save(fhdr *FileHeader, filename string) error {
	blockSize := fhdr.AccoCodeLength + fhdr.RoomRateCodeLength + FixIdxRecSize
	buf := make([]byte, blockSize)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	idx.Lock()
	for accoCode, roomRateMap := range idx.m {
		for roomRateCode, occupancies := range roomRateMap {
			for _, occupancy := range occupancies {
				copy(buf[0:], []byte(accoCode))
				copy(buf[fhdr.AccoCodeLength:], []byte(roomRateCode))
				copy(buf[fhdr.AccoCodeLength+fhdr.RoomRateCodeLength:], *occupancy.ToByteStr())
				//copy(buf[blockSize-4:])
				f.Write(buf)
			}
		}
	}
	idx.Unlock()
	return nil
}

// Load reads the cache index from a file.
func (idx *CacheIndex) Load(fhdr *FileHeader, filename string) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, 644)
	if err != nil {
		return err
	}
	defer f.Close()
	recordSize := int64(fhdr.AccoCodeLength + fhdr.RoomRateCodeLength + FixIdxRecSize)
	statInfo, err := f.Stat()
	if err != nil {
		return err
	}
	fSize := statInfo.Size()
	if fSize%recordSize != 0 {
		return errors.New("Incorrect file size. File may be corrupt")
	}
	buf := make([]byte, recordSize)
	var accoCode string
	var roomRateCode string
	var idxValue uint32
	var roomOccIdx RoomOccIdx
	recordCount := fSize / recordSize
	for i := int64(0); i < recordCount; i++ {
		f.ReadAt(buf, i*recordSize)
		accoCode = string(bytes.Trim(buf[0:fhdr.AccoCodeLength], "\x00"))
		roomRateCode = string(bytes.Trim(buf[fhdr.AccoCodeLength:fhdr.AccoCodeLength+fhdr.RoomRateCodeLength], "\x00"))
		idxValue = binary.BigEndian.Uint32(buf[recordSize-4 : recordSize])
		roomOccIdx = RoomOccIdx{Idx: idxValue}
		for j := 0; j < 8; j++ {
			offStart := int(fhdr.AccoCodeLength+fhdr.RoomRateCodeLength) + j*3
			tripleBuf := buf[offStart : offStart+3]
			if tripleBuf[2] > 0 {
				roomOccIdx.AddOccItem(uint8(tripleBuf[0]), uint8(tripleBuf[1]), uint8(tripleBuf[2]))
			}
		}
		idx.AddRoomOccIdx(accoCode, roomRateCode, roomOccIdx)
	}
	return nil
}

func (idx *CacheIndex) LoadFromCache(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, 644)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, FileHeaderSize)
	f.Read(buf)
	fhdr, err := FileHeaderFromByteStr(buf)
	if err != nil {
		return err
	}
	blockHeaderSize := fhdr.GetBlockHeaderSize()
	hdrbuf := make([]byte, blockHeaderSize)
	var accoCode string
	var roomRateCode string
	var roomOccIdx RoomOccIdx
	for i := uint32(0); i < fhdr.RateBlockCount; i++ {
		f.ReadAt(hdrbuf, fhdr.GetRateBlockStart(i))
		accoCode = string(bytes.Trim(hdrbuf[0:fhdr.AccoCodeLength], "\x00"))
		roomRateCode = string(bytes.Trim(hdrbuf[fhdr.AccoCodeLength:fhdr.AccoCodeLength+fhdr.RoomRateCodeLength], "\x00"))
		roomOccIdx = RoomOccIdx{Idx: i}
		for j := int(fhdr.AccoCodeLength + fhdr.RoomRateCodeLength); j < blockHeaderSize; j += 3 {
			if uint8(hdrbuf[j+2]) > 0 {
				roomOccIdx.AddOccItem(uint8(hdrbuf[j]), uint8(hdrbuf[j+1]), uint8(hdrbuf[j+2]))
			}
		}
		idx.AddRoomOccIdx(accoCode, roomRateCode, roomOccIdx)
	}
	return nil
}

func cmpOccupancy(occ1 []OccupancyItem, occ2 []OccupancyItem) bool {
	if len(occ1) != len(occ2) {
		return false
	}
	for _, oi1 := range occ1 {
		match := false
		for _, oi2 := range occ2 {
			if oi1.MinAge == oi2.MinAge && oi1.MaxAge == oi2.MaxAge && oi1.Count == oi2.Count {
				match = true
				break
			}
		}
		if match == false {
			return false
		}
	}
	return true
}

func (idx *CacheIndex) Get(q IndexQuery) (uint32, bool) {
	idx.Lock()
	occupancies := idx.m[q.AccoCode][q.RoomRateCode]
	for _, occupancy := range occupancies {
		if occupancy.Total == q.OccTotal {
			if cmpOccupancy(q.Occupancy, occupancy.Occupancy) == true {
				index := occupancy.Idx
				idx.Unlock()
				return index, true
			}
		}
	}
	idx.Unlock()
	return 0, false
}
