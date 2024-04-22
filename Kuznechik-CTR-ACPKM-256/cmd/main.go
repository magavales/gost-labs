package main

import (
	"Kuznechik-CTR-ACPKM-256/pkg"
	"fmt"
	"log"
	"os"
)

func main() {
	var (
		key = [32]byte{
			0x80, 0x94, 0xA8, 0xBC, 0xC0, 0xD4, 0xE8, 0xFC,
			0x81, 0x95, 0xA9, 0xBD, 0xC1, 0xD5, 0xE9, 0xFD,
			0x82, 0x96, 0xAA, 0xBE, 0xC2, 0xD6, 0xEA, 0xFE,
			0x83, 0x97, 0xAB, 0xBF, 0xC3, 0xD7, 0xEB, 0xFF,
		}

		plainText = [16]byte{
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
			0xAA, 0xBB, 0xCC, 0xDD,
		}
	)
	file, err := os.Open("gost.exe")
	if err != nil {
		log.Fatalf("%s", err)
	}

	channel := pkg.NewChannel()

	go pkg.Validation(channel)
	channel.File <- file
	channel.Hash <- nil

	temp := <-channel.Hash
	log.Println("The validation function finished!")

	fmt.Printf("Plain: %v\n", plainText)

	cipherMode := pkg.NewCtrAcpkm(key[:])
	ciphertext, mac := cipherMode.Encrypt(plainText[:], key[:])
	fmt.Printf("Encrypt by CTR-ACPKM: %v\n", ciphertext)
	fmt.Printf("MAC(HMAC-SHA-256) for text: %v\n", mac)

	channel.File <- file
	channel.Hash <- temp

	temp = <-channel.Hash
	log.Println("The validation function finished!")

	text := cipherMode.Decrypt(ciphertext[:], key[:], mac)
	fmt.Printf("Decrypt by CTR-ACPKM: %v\n", text)
	fmt.Printf("Plain == Decrypt: %t\n", plainText == text)

	/*c := pkg.NewCipher(key[:])

	encrypted := c.Encrypt(plainText[:])
	fmt.Printf("Encrypt: %v\n", *encrypted)

	decrypt := c.Decrypt(encrypted[:])
	fmt.Printf("Decrypt: %v\n", *decrypt)

	fmt.Printf("Plain == Decrypt: %t", plainText == *decrypt)*/

	_ = file.Close()
	close(channel.File)
	close(channel.Hash)
}
