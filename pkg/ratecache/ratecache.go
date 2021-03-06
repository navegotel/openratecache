//Package ratecache provides functions to efficiently store rate and
//availability data in a buffer or a file.

package ratecache

import (
    "os"
    "time"
    "errors"
    "fmt"
    "bytes"
    "encoding/binary"
    "path/filepath"
)

//Convert the Time object into a date string as used 
//in the rate file header.
func TimeToStr(t time.Time) string {
    return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

//Take a date string in rate header format and returns a time.Time
//object.
func StrToTime(s string) (time.Time, error) {
    t, err := time.Parse("20060102", s)
    return t, err
}

//Packs rate and availability into a single uint32.
func PackRate(rate uint32, avail uint8) []byte {
    r := (rate & RateMask) | (uint32(avail) << 28)
    buf := make([]byte, 4)
    binary.BigEndian.PutUint32(buf, r)
    return buf
}

//Unpacks values for rate and availability and returns them separately.
func UnpackRate(buf []byte) (uint32, uint8){
    r := binary.BigEndian.Uint32(buf)
    rate := r & RateMask
    avail := uint8(r >> 28)
    return rate, avail
}

//Represent a guest type identified by a date range that
//can occupy a room. One occupancy is made up of one or more
//OccupancyItems.
type OccupancyItem struct {
    MinAge uint8
    MaxAge uint8
    Count uint8
}

func NewOccupancyItem(MinAge uint8, MaxAge uint8, Count uint8) (*OccupancyItem, error) {
    if MinAge > MaxAge {
        err := errors.New("MinAge cannot be greater than MaxAge")
        return nil, err
    }
    item := OccupancyItem{MinAge: MinAge, MaxAge: MaxAge, Count: Count}
    return &item, nil
}

func (item *OccupancyItem) ToByteStr() *[]byte {
    byteStr := make([]byte, 3)
    byteStr[0] = item.MinAge
    byteStr[1] = item.MaxAge
    byteStr[2] = item.Count
    return &byteStr
}

//Represent header information from a rate cache file which can be
//read from and written to a byte string.
type FileHeader struct {
    Signature string
    Supplier string
    Version uint8
    StartDate time.Time
    Currency string
    MaxLos uint8
    Days uint16
    AccoCodeLength uint8
    RoomRateCodeLength uint8
    RateBlockCount uint16
}

//file header constructor
func NewFileHeader(Supplier string, StartDate time.Time, Currency string, MaxLos uint8, Days uint16, AccoCodeLength uint8, RoomRateCodeLength uint8)(*FileHeader, error){
    if len(Currency) != 3 {
        return nil, errors.New("currency must be a 3 letter ISO currency code")
    }
    if len(Supplier) > 8 {
        return nil, errors.New("supplier code must not be longer than 8 bytes")
    }
    fhdr := FileHeader{Signature: Signature, Version: Version, Supplier: Supplier, StartDate: StartDate, Currency: Currency, MaxLos: MaxLos, Days: Days, AccoCodeLength: AccoCodeLength, RoomRateCodeLength: RoomRateCodeLength, RateBlockCount: 0}
    return &fhdr, nil
}

//Parse a file header into a FileHeader object.
func FileHeaderFromByteStr (byteStr []byte)(*FileHeader, error){
    if string(byteStr[:8]) != Signature {
        return nil, errors.New("byteStr is not in rate cache format")
    }
    if int(byteStr[16]) != Version {
        return nil, errors.New(fmt.Sprintf("Wrong version. expected version is %d, got %d", Version, int(byteStr[16])))  
    }
    fhdr := FileHeader{Signature: Signature, Version: Version}
    fhdr.Supplier = string(bytes.Trim(byteStr[8:16], "\x00"))
    //fmt.Println(string(byteStr[17:25]))
    t, err := StrToTime(string(byteStr[17:25]))
    if err != nil {
        return nil, err
    }
    fhdr.StartDate = t
    fhdr.Currency = string(byteStr[25:28])
    fhdr.MaxLos = uint8(byteStr[28])
    fhdr.Days = binary.BigEndian.Uint16(byteStr[29:31])
    fhdr.AccoCodeLength = uint8(byteStr[31])
    fhdr.RoomRateCodeLength = uint8(byteStr[32])
    fhdr.RateBlockCount = binary.BigEndian.Uint16(byteStr[33:35])
    return &fhdr, nil
}

//Calculate rate block header size.
func (fhdr *FileHeader) GetBlockHeaderSize() int {
    blockHeaderSize := int(fhdr.AccoCodeLength) + int(fhdr.RoomRateCodeLength) + int(FixBlockHeaderSize)
    return blockHeaderSize
}

//Return the total size of a rate block including the header
func (fhdr *FileHeader) GetRateBlockSize() int {
    blockSize := fhdr.GetBlockHeaderSize() + int(fhdr.Days) * int(fhdr.MaxLos) * 4
    return blockSize
}

//Return offset of rate block from its index. First block has index 0.
func (fhdr *FileHeader) GetRateBlockStart(index uint16) int64 {
    return int64(FileHeaderSize) + int64(fhdr.GetRateBlockSize()) * int64(index)
}

//Get position of rate in rate cache
func (fhdr *FileHeader) GetRatePos(idx uint16, date time.Time, los uint8) (int64, error){
    if fhdr.RateBlockCount == 0 {
        return 0, errors.New("Rate cache is empty")
    }
    if (fhdr.RateBlockCount - 1) < idx {
        return 0, errors.New("Index too big, not enough rate blocks.")
    }
    losBlockOffset := int64(fhdr.GetBlockHeaderSize()) + (int64(los - 1) * int64(fhdr.Days) * 4)
    dayOffset := int64(date.Sub(fhdr.StartDate ).Hours()/24) * 4
    losStart := losBlockOffset + dayOffset
    rateBlockStart := fhdr.GetRateBlockStart(idx)
    //fmt.Printf(" los: %d\n BlockHeaderSize: %d\n losBlockOffset: %d\n dayOffset: %d\n rateBlockStart: %v\n losStart: %v\n ", los, fhdr.GetBlockHeaderSize(), losBlockOffset, dayOffset, rateBlockStart, losStart)
    return rateBlockStart + losStart, nil
}

//Write one rate/avail to rate cache.
func (fhdr *FileHeader) SetRateInfo(f *os.File, idx uint16, date time.Time, los uint8, rate uint32, avail uint8) error {
    val := PackRate(rate, avail)
    ratePos, err := fhdr.GetRatePos(idx, date, los)
    if err != nil {
        return err   
    }
    _, err = f.WriteAt(val, ratePos)
    f.Sync()
    return err
}

//Get one rate/avail from rate cache.
func (fhdr *FileHeader) GetRateInfo(f *os.File, idx uint16, date time.Time, los uint8) (uint32, uint8, error) {
    ratePos, err := fhdr.GetRatePos(idx, date, los)
    if err != nil {
        return 0, 0, err   
    }
    buf := make([]byte, 4)
    f.ReadAt(buf, ratePos)
    rate, avail := UnpackRate(buf)
    return rate, avail, nil
}

//Create a file header as byte string from object.
func (fhdr *FileHeader) ToByteStr() []byte {
    byteStr := []byte(Signature)
    byteStr = append(byteStr, fhdr.Supplier...)
    for i:=len(fhdr.Supplier); i < 8; i++ {
        byteStr = append(byteStr, byte(0))
    }
    byteStr = append(byteStr, fhdr.Version)
    datestr := TimeToStr(fhdr.StartDate)
    byteStr = append(byteStr, datestr...)
    byteStr = append(byteStr, fhdr.Currency...)
    byteStr = append(byteStr, fhdr.MaxLos)
    daysStr := make([]byte, 2)
    binary.BigEndian.PutUint16(daysStr, fhdr.Days)
    byteStr = append(byteStr, daysStr...)
    byteStr = append(byteStr, fhdr.AccoCodeLength)
    byteStr = append(byteStr, fhdr.RoomRateCodeLength)
    countStr := make([]byte, 2)
    binary.BigEndian.PutUint16(countStr, fhdr.RateBlockCount)
    byteStr = append(byteStr, countStr...)
    return byteStr
}

//Represent rate block header information which can be read 
//from and written to byte string.
type RateBlockHeader struct {
    accoCode string
    roomRateCode string
    occupancy []*OccupancyItem
}

//Rate block header constructor
func NewRateBlockHeader(accoCode string, roomRateCode string)(*RateBlockHeader, error){
    rbhdr := RateBlockHeader{accoCode: accoCode, roomRateCode: roomRateCode}
    return &rbhdr, nil
}

//Create rate block header object from byte string.
func RateBlockHeaderFromByteStr(byteStr []byte, AccoCodeLength uint8, RoomRateCodeLength uint8)(*RateBlockHeader, error){
    if int(AccoCodeLength) + int(RoomRateCodeLength) > len(byteStr) {
        return nil, errors.New("byteStr is not long enough or accoCodeLength/roomRateCodeLength are wrong")
    }
    rbhdr := RateBlockHeader{}
    bytes.Trim(byteStr[8:16], "\x00")
    rbhdr.accoCode = string(bytes.Trim(byteStr[:AccoCodeLength], "\x00"))
    rbhdr.roomRateCode = string(bytes.Trim(byteStr[AccoCodeLength:RoomRateCodeLength], "\x00"))
    offset := int(AccoCodeLength + RoomRateCodeLength)
    for i := 0; i < 24; i += 3 {
        item := byteStr[offset + i: offset + i + 3]
        if item[2] > 0 {
            rbhdr.AddOccupancyItem(uint8(item[0]), uint8(item[1]), uint8(item[2]))
        }
    }
    return &rbhdr, nil
}

//Add occupancy item to the rate block header.
func (rbhdr *RateBlockHeader)AddOccupancyItem(MinAge uint8, MaxAge uint8, Count uint8) error {
    if len(rbhdr.occupancy) == 8 {
        err := errors.New("Cannot add mor than 8 occupancy items")
        return err
    }
    item, err := NewOccupancyItem(MinAge, MaxAge, Count)
    if err != nil {
        return err
    }
    rbhdr.occupancy = append(rbhdr.occupancy, item)
    return nil
}

//Create rate block header as byte string from object
func (rbhdr *RateBlockHeader) ToByteStr(AccoCodeLength uint8, RoomRateCodeLength uint8) []byte {
    byteStr := []byte(rbhdr.accoCode)
    for i:=uint8(len(rbhdr.accoCode)); i < AccoCodeLength; i++ {
        byteStr = append(byteStr, byte(0))
    }
    byteStr = append(byteStr, rbhdr.roomRateCode...)
    for i:=uint8(len(rbhdr.roomRateCode)); i < RoomRateCodeLength; i++ {
        byteStr = append(byteStr, byte(0))
    }
    for _, v := range(rbhdr.occupancy) {
        byteStr = append(byteStr, *v.ToByteStr()...)
    }
    padding := make([]byte, (8 - len(rbhdr.occupancy)) * 3)
    byteStr = append(byteStr, padding...)
    return byteStr
}

//Create empty rate block.
func CreateRateBlock(fhdr *FileHeader, rbhdr *RateBlockHeader) []byte {
    byteStr := rbhdr.ToByteStr(fhdr.AccoCodeLength, fhdr.RoomRateCodeLength)
    bsLength := 4 * int(fhdr.MaxLos) * int(fhdr.Days)
    for i := 0; i < bsLength; i++ {
        byteStr = append(byteStr, byte(0))
    }
    return byteStr
}

//Create a new rate file on disc and resets the rateBlockCount of
//file header object. folder is the folder, where the file 
//is going to live.
func InitRateFile(fhdr *FileHeader, folder string, blockCount int) (string, error) {
    fhdr.RateBlockCount = 0
    byteStr := fhdr.ToByteStr()
    rateBlockSize := fhdr.GetRateBlockSize()
    emptyBlock := make([]byte, rateBlockSize)
    t := time.Now()
    filename := fhdr.Supplier + t.Format(".20060102150405")+ ".bin" 
    f, err := os.OpenFile(filepath.Join(folder, filename), os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_WRONLY, 0644)
    defer f.Close()
    if err != nil {
        return "", err
    }
    count, err := f.Write(byteStr)
    if err != nil {
        return "", err
    }
    if count != FileHeaderSize {
        return "", errors.New("Could not write complete file header")
    }
    for i := 0; i < blockCount; i++ {
        count, err := f.Write(emptyBlock)
        if err != nil {
            return "", err
        }
        if count != rateBlockSize {
        return "", errors.New("Could not write block padding")
    }
    }
    return filename, nil
}

//Add rate block to cache file. The block position is determined by
//the RateBlockCount. Method will return the index, not the count!
//Do not forget to update the rateBlockCount on your FileHeader object!
func AddRateBlockToFile(f *os.File, byteStr []byte) (uint16, error) {
    var err error
    buf := make([]byte, 2)
    f.ReadAt(buf, 33)
    RateBlockCount := binary.BigEndian.Uint16(buf)
    blockSize := len(byteStr)
    _, err = f.WriteAt(byteStr, int64(FileHeaderSize + int(RateBlockCount) * blockSize))
    if err != nil {
        return 0, err
    }
    f.Sync()
    RateBlockCount++
    binary.BigEndian.PutUint16(buf, RateBlockCount)
    _, err = f.WriteAt(buf, 33)
    if err != nil {
        return 0, err
    }
    return RateBlockCount - 1, nil
}

