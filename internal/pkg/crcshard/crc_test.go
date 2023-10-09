package crcshard

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func BenchmarkCrcSharder_GetShard(b *testing.B) {
	uuids := make([]string, 10000000)
	for i := range uuids {
		uuids[i] = uuid.NewV4().String()
	}
	sharder := New(10)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sharder.GetShardStr(uuids[n%10000000])
	}
}
