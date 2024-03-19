package pkg

type RoundKeys struct {
	Keys [10][16]uint8
}

func (rk *RoundKeys) GetRoundKeys(key [32]uint8) {
	var (
		C, z [16]uint8
	)

	x, y := splitKey(key)

	rk.Keys[0] = x
	rk.Keys[1] = y

	for i := 0; i <= 32; i++ {
		for k := range C {
			C[k] = 0
		}
		C[15] = uint8(i)
		C = L(C)

		for k := range z {
			z[k] = Pi[(x[k] ^ C[k])]
		}
		z = L(z)
		for k := range z {
			z[k] = z[k] ^ y[k]
		}

		y = x
		x = z

		if i%8 == 0 {
			rk.Keys[(i >> 2)] = x
			rk.Keys[(i>>2)+1] = y
		}
	}
}

func splitKey(key [32]uint8) ([16]uint8, [16]uint8) {
	var x, y [16]uint8
	for i := 0; i < 16; i++ {
		x[i] = key[i]
		y[i] = key[i+16]
	}

	return x, y
}
