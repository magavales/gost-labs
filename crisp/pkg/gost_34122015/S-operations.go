package pkg

func S(block *[BlockSize]byte) {
	for i := 0; i < BlockSize; i++ { // substitute byte by S-Box
		block[i] = SBox[int(block[i])]
	}
}

// S-substitute (inverse)
func SInverse(block *[BlockSize]byte) {
	for i := 0; i < BlockSize; i++ { // substitute byte by inverse S-Box
		block[i] = SBoxInv[int(block[i])]
	}
}
