package crisp

import (
	gost132356510222018 "crisp/pkg/gost_132356510222018"
	gost34122015 "crisp/pkg/gost_34122015"
	xoroshiro256plus "crisp/pkg/xoroshiro256+"
)

type Encoder struct {
	random *xoroshiro256plus.XoroShiroPlus256
	kdf    *gost132356510222018.KDF
	cipher *gost34122015.CtrAcpkm
	seqNum uint32
}
