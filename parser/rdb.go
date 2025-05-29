package parser

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/JanitSri/codecrafters-build-your-own-redis/data"
)

// https://rdb.fnordig.de/file_format.html
func ParseRBDFile(b []byte) map[string]*data.RedisValue {
	pairs := make(map[string]*data.RedisValue)

	// skip the header section -- will always be 'REDIS0011'
	i := 9

	for b[i] != 0xFF {
		if b[i] == 0xFA {
			// parse metadata section
			// contains zero or more "metadata subsections," which each specify a single metadata attribute
			i++
			mki, mk := parseString(b, i)
			i = mki

			mvi, mv := parseString(b, i)
			i = mvi

			log.Printf("metadata key: %s, metadata value: %s", mk, mv)

		} else if b[i] == 0xFE {
			// parse database section
			// contains zero or more "database subsections," which each describe a single database

			// skip the index of the database & 0xFB
			i += 3

			// represents the size of the hash table that stores keys and values
			j := int(b[i])
			// skip size of the hash table that stores the expires of the keys
			i += 2

			for j > 0 {
				if b[i] == 0x00 {
					i = parseType(i, b, pairs, time.Time{})
				} else if b[i] == 0xFC {
					i++
					// expire time expressed in milliseconds, stored as an 8-byte unsigned long
					ms_start := i
					mss := 8
					i += mss
					mse := b[ms_start:i]
					e := binary.LittleEndian.Uint64(mse)

					i = parseType(i, b, pairs, time.Unix(int64(e/1000), int64(e%1000)*int64(time.Millisecond)))
				} else if b[i] == 0xFD {
					i++
					// expire time expressed in seconds, stored as an 4-byte unsigned integer
					ms_start := i
					mss := 4
					i += mss
					mse := b[ms_start:i]
					e := binary.LittleEndian.Uint32(mse)

					i = parseType(i, b, pairs, time.Unix(int64(e), 0))
				}
				j--
			}
		}
	}

	return pairs
}

func parseType(i int, b []byte, pairs map[string]*data.RedisValue, exp time.Time) int {
	switch b[i] {
	case 0x00:
		i++
		ki, key := parseString(b, i)
		i = ki

		vi, val := parseString(b, i)
		i = vi

		rv := data.NewRedisValue(val, exp)
		if !rv.IsExpired() {
			pairs[fmt.Sprintf("%v", key)] = rv
		}

	default:
		log.Fatal(ErrInvalidRDBValTypeError)
	}

	return i
}

func parseString(b []byte, i int) (int, any) {
	var out any
	if b[i] == 0xC0 || b[i] == 0xC2 {
		// integers as a string
		le := b[i] & 0x3F // the last six bits
		i++
		if le == 0 {
			// 0 indicates that an 8 bit integer
			i8 := int8(b[i])
			out = strconv.FormatInt(int64(i8), 10)
			i++
		} else if le == 1 {
			// 1 indicates that a 16 bit integer
			size := 2
			n := b[i : size+i]
			i16 := int16(binary.LittleEndian.Uint16(n))
			out = strconv.FormatInt(int64(i16), 10)
			i += size
		} else if le == 2 {
			// 2 indicates that a 32 bit integer
			size := 4
			n := b[i : size+i]
			i32 := int32(binary.LittleEndian.Uint32(n))
			out = strconv.FormatInt(int64(i32), 10)
			i += size
		}
	} else {
		// length prefixed string
		s := i
		i += int(b[i]) + 1
		out = string(b[s+1 : i])
	}

	return i, out
}
