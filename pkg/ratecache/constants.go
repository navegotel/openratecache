package ratecache

// Signature is the signature string for rate files
const Signature = "LOSRATES"

// Version is the format version of the rate file
const Version = 8

// FileHeaderSize is the size of the rate file header in bytes
const FileHeaderSize = 35

// FixBlockHeaderSize is the portion of the block header size
// that does not chane, i.e. without room rate code and acco code
const FixBlockHeaderSize = 24

// FixIdxRecSize is the portion of the record size in the
// index file that does not change, i.e. without room rate code
// and acco code.
const FixIdxRecSize = 26

// RateMask masks the upper 4 bytes of an uint32 which is used
// to transport availability
const RateMask uint32 = 268435455
