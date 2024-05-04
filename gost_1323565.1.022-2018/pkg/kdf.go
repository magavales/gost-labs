package pkg

import (
	"fmt"
	"math/rand"
	"strconv"
)

const KeySize = 32

type KDF struct {
	Key  []byte
	Size int
	T    []byte
	P    []byte
	U    []byte
	A    []byte
}

func NewKDF(key, T, P, U, A []byte) *KDF {
	return &KDF{
		Key:  key,
		Size: KeySize,
		T:    T,
		P:    P,
		U:    U,
		A:    A,
	}
}

func (k *KDF) Generate() []byte {
	k1 := firstRound(k.T, k.Key)
	k2 := secondRound(k1, k.P, k.U, k.A)

	return k2
}

func (k *KDF) Clear() {
	k.Key = []byte(fmt.Sprintf(strconv.Itoa(rand.Intn(10000000))))
	k.Size = rand.Intn(1000000)
	k.T = []byte(fmt.Sprintf(strconv.Itoa(rand.Intn(100000))))
	k.P = []byte(fmt.Sprintf(strconv.Itoa(rand.Intn(10000))))
	k.U = []byte(fmt.Sprintf(strconv.Itoa(rand.Intn(10000000))))
	k.A = []byte(fmt.Sprintf(strconv.Itoa(rand.Intn(10000000))))
}

func firstRound(T, key []byte) []byte {
	mac := Hmac(key, T)

	return mac[:32]
}

func secondRound(k1, P, U, A []byte) []byte {
	C := uint64(1)
	zi := [64]byte{0x80, 0x94, 0xA8, 0xBC, 0xC0, 0xD4, 0xE8, 0xFC,
		0x81, 0x95, 0xA9, 0xBD, 0xC1, 0xD5, 0xE9, 0xFD,
		0x82, 0x96, 0xAA, 0xBE, 0xC2, 0xD6, 0xEA, 0xFE,
		0x83, 0x97, 0xAB, 0xBF, 0xC3, 0xD7, 0xEB, 0xFF}

	retString := make([]byte, 64)

	for i := uint64(0); i < 1; i++ {
		formatStr := format(zi[:], A, setL(), P, U, C)
		mac := Hmac(k1, formatStr)
		copy(zi[:], mac[:])
		copy(retString[i*32:], zi[:])
		C++
	}
	result := make([]byte, KeySize)
	copy(result, retString)

	return result
}

func format(zi, L, P, U, A []byte, ci uint64) []byte {
	retLen := 1 + 32 + 64 + len(L) + len(P) + len(A) + len(U)
	retString := make([]byte, retLen)
	pos := 0
	retString[pos] = byte(0xFC)
	pos++
	C := make([]byte, 32)
	for i := 0; i < 32; i++ {
		C[i] = byte((ci >> (uint(i) * 8)) & 0xff)
	}
	for i := 0; i < 32; i++ {
		retString[pos] = C[i]
		pos++
	}
	for i := 0; i < 64; i++ {
		retString[pos] = zi[i]
		pos++
	}
	for i := 0; i < len(L); i++ {
		retString[pos] = L[i]
		pos++
	}
	for i := 0; i < len(P); i++ {
		retString[pos] = P[i]
		pos++
	}
	for i := 0; i < len(U); i++ {
		retString[pos] = U[i]
		pos++
	}
	for i := 0; i < len(A); i++ {
		retString[pos] = A[i]
		pos++
	}

	return retString
}

func setL() []byte {
	var L [8]byte
	for i := 0; i < 8; i++ {
		L[i] = byte((KeySize >> (i * 8)) & 0xff)
	}
	return L[:]
}
