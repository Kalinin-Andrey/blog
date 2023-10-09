package blog

import (
	"encoding"
	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*Blog)(nil)
var _ encoding.BinaryUnmarshaler = (*Blog)(nil)

func (e *Blog) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *Blog) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
