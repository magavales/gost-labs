package pkg

import (
	"encoding/binary"
	"math/bits"
)

const (
	Size = 32
)

type XoroShiroPlus256 struct {
	seed [Size]byte
}

func New(seed [Size]byte) *XoroShiroPlus256 {
	return &XoroShiroPlus256{seed: seed}
}

func (x *XoroShiroPlus256) Next() uint64 {
	s0 := binary.BigEndian.Uint64(x.seed[0 : Size/4])
	s1 := binary.BigEndian.Uint64(x.seed[Size/4 : 2*Size/4])
	s2 := binary.BigEndian.Uint64(x.seed[2*Size/4 : 3*Size/4])
	s3 := binary.BigEndian.Uint64(x.seed[3*Size/4:])

	binary.BigEndian.PutUint64(x.seed[0:Size/4], s0^s3^s1)
	binary.BigEndian.PutUint64(x.seed[Size/4:2*Size/4], s1^s2^s0)
	binary.BigEndian.PutUint64(x.seed[2*Size/4:3*Size/4], s2^s0^(s1<<17))
	binary.BigEndian.PutUint64(x.seed[3*Size/4:], bits.RotateLeft64(s3^s1, 45))

	return s0 + s3
}
