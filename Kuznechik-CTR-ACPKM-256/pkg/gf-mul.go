package pkg

func GF_mul(x, y uint8) uint8 {
	var z uint8 = 0

	for y != 0 {
		if y&1 == 1 {
			z = z ^ x
		}
		if x&0x80 != 0 {
			x = (x << 1) ^ 0xC3
		} else {
			x = x << 1
		}
		y = y >> 1
	}
	return z
}
