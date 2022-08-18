package gocrc

import (
	"hash/crc32"
	"hash/crc64"
)

func CheckSum32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func CheckSum64ECMA(data []byte) uint64 {
	return crc64.Checksum(data, crc64.MakeTable(crc64.ECMA))
}

func CheckSum64ISO(data []byte) uint64 {
	return crc64.Checksum(data, crc64.MakeTable(crc64.ISO))
}
