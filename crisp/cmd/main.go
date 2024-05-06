package main

import (
	"bytes"
	"crisp/pkg/crisp"
	validation "crisp/pkg/validation"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	BlockSize = 16
	KeySize   = 32
)

func main() {
	var (
		randomSeed [32]byte
		plainText  = [BlockSize * 4]byte{
			0xA0, 0xB0, 0xC0, 0xD0, 0xE0, 0xF0, 0xA0, 0xB0,
			0xC0, 0xD0, 0xE0, 0xF0, 0xA0, 0xB0, 0xC0, 0xD0,

			0xA1, 0xB1, 0xC1, 0xD1, 0xE1, 0xF1, 0xA1, 0xB1,
			0xC1, 0xD1, 0xE1, 0xF1, 0xA1, 0xB1, 0xC1, 0xD1,

			0xA2, 0xB2, 0xC2, 0xD2, 0xE2, 0xF2, 0xA2, 0xB2,
			0xC2, 0xD2, 0xE2, 0xF2, 0xA2, 0xB2, 0xC2, 0xD2,

			0xA3, 0xB3, 0xC3, 0xD3, 0xE3, 0xF3, 0xA3, 0xB3,
			0xC3, 0xD3, 0xC3, 0xF3, 0xA3, 0xB3, 0xC3, 0xD3,
		}

		decodeText []byte

		key = [KeySize]byte{
			0x80, 0x94, 0xA8, 0xBC, 0xC0, 0xD4, 0xE8, 0xFC,
			0x81, 0x95, 0xA9, 0xBD, 0xC1, 0xD5, 0xE9, 0xFD,
			0x82, 0x96, 0xAA, 0xBE, 0xC2, 0xD6, 0xEA, 0xFE,
			0x83, 0x97, 0xAB, 0xBF, 0xC3, 0xD7, 0xEB, 0xFF,
		}
	)

	fmt.Printf("Plain: %s\n", hex.EncodeToString(plainText[:]))

	// Password = "Pa$$w0rd"
	if len(os.Args) == 2 {
		if os.Args[1] != "97c94ebe5d767a353b77f3c0ce2d429741f2e8c99473c3c150e2faa3d14c9da6" {
			log.Fatalln("Password isn't corrected!")
		} else {
			log.Println("Password is corrected!")
		}
	} else {
		log.Fatalln("Password didn't send to program!")
	}

	channel := validation.NewChannel()

	go validation.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- nil

	temp := <-channel.Hash
	log.Println("The validation function finished!")

	binary.BigEndian.PutUint16(randomSeed[:], uint16(time.Now().Nanosecond()))
	newCrisp := crisp.NewCrisp(key[:], randomSeed)
	defer newCrisp.Clear()

	b, err := os.ReadFile("xorshift_1mb.bin")
	if err != nil {
		log.Fatalln("failed to open file: " + err.Error())
	}

	start := time.Now().Second()
	for i := 0; i < len(b); i += BlockSize {
		message := newCrisp.EncodeNextBlock(b[i : i+BlockSize])
		decoded := newCrisp.DecodeNextBlock(message.Digits)

		if !bytes.Equal(b[i:i+BlockSize], decoded) {
			log.Fatalln("incorrect decode")
		}

		decodeText = append(decodeText, decoded...)

		if i%(4096*4) == 0 {
			fmt.Printf("Complete %d/%d [%.2f%%]\n", i, len(b), float64(i)/float64(len(b))*100)
		}
	}

	go validation.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- temp

	temp = []byte{0x23}

	_ = <-channel.Hash
	log.Println("The validation function finished!")

	fmt.Printf("Elapsed: %d msec.\n", time.Now().Second()-start)

}
