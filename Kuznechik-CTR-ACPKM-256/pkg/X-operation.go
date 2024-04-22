package pkg

import "log"

func X(dst, src1, src2 []byte) {
	if len(dst) != BlockSize {
		log.Fatalln("dst is not 16 bytes")
	}
	if len(src1) != BlockSize {
		log.Fatalln("src1 is not 16 bytes")
	}
	if len(src2) != BlockSize {
		log.Fatalln("src2 is not 16 bytes")
	}
	for i := 0; i < BlockSize; i++ {
		dst[i] = src1[i] ^ src2[i]
	}
}
