package crcshard

import (
	"hash/crc64"
	"strconv"
)

type CrcSharder struct {
	shardCount uint
	crcTable   *crc64.Table
}

func New(shardCount byte) *CrcSharder {
	return &CrcSharder{shardCount: uint(shardCount), crcTable: crc64.MakeTable(crc64.ISO)}
}

// GetShardStr: returns number of shard, counts from 1
func (c *CrcSharder) GetShardStr(key string) byte {
	if id, err := strconv.ParseUint(key, 10, 64); err == nil {
		return byte(uint(id) % c.shardCount)
	}
	return byte(uint(crc64.Checksum([]byte(key), c.crcTable))%c.shardCount) + 1
}

// GetShardStr: returns number of shard, counts from 1
func (c *CrcSharder) GetShardUint(key uint) byte {
	return byte(key%c.shardCount) + 1
}
