package pkg

func L(block [16]uint8) [16]uint8 {
	var (
		i, j int
		x    uint8
	)

	for j = 0; j < 16; j++ {
		x = block[15]
		for i = 14; i >= 0; i-- {
			block[i+1] = block[i]
			x = x ^ GF_mul(block[i], LVector[i])
		}
		block[0] = x
	}

	return block
}

func LInv(block [16]uint8) [16]uint8 {
	var (
		i, j int
		x    uint8
	)
	for j = 0; j < 16; j++ {
		x = block[0]
		for i = 0; i < 15; i++ { // Just process in reverse sequence.
			block[i] = block[i+1]
			x = x ^ GF_mul(block[i], LVector[i])
		}
		block[15] = x
	}

	return block
}
