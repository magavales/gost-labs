package pkg

func L(block *[BlockSize]byte) {
	for n := 0; n < BlockSize; n++ {
		var t byte
		for i := 0; i < BlockSize; i++ {
			t ^= GfCache[block[i]][LVec[i]]
		}
		for i := BlockSize - 1; i > 0; i-- {
			block[i] = block[i-1]
		}
		block[0] = t
	}
}

// L-function (inverse)
func LInverse(block *[BlockSize]byte) {
	var t byte
	for n := 0; n < BlockSize; n++ {
		t = block[0]
		copy(block[:], block[1:])
		for i := 0; i < BlockSize-1; i++ {
			t ^= GfCache[block[i]][LVec[i]]
		}
		block[15] = t
	}
}
