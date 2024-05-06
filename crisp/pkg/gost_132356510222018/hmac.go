package pkg

import (
	"github.com/deatil/go-hash/streebog"
)

func Hmac(key []byte, data []byte) []byte {
	b := 128
	k0 := make([]byte, b)
	h := streebog.New512()

	if len(key) > b {
		h.Write(key)
		hash := h.Sum(nil)
		copy(k0, hash[:])
	} else {
		copy(k0, key)
	}

	if len(key) < b {
		padding := make([]byte, b-len(key))
		k0 = append(k0, padding...)
	}

	ipad := make([]byte, b)
	opad := make([]byte, b)
	for i := 0; i < b; i++ {
		ipad[i] = k0[i] ^ 0x36
		opad[i] = k0[i] ^ 0x5c
	}

	h.Write(append(ipad, data...))
	keypad := h.Sum(nil)
	h.Write(append(opad, keypad[:]...))

	return h.Sum(nil)
}
