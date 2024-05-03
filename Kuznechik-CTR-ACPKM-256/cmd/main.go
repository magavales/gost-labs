package main

import (
	"Kuznechik-CTR-ACPKM-256/pkg"
	"bytes"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	startTime := time.Now()
	var (
		wg  sync.WaitGroup
		key = [32]byte{
			0x80, 0x94, 0xA8, 0xBC, 0xC0, 0xD4, 0xE8, 0xFC,
			0x81, 0x95, 0xA9, 0xBD, 0xC1, 0xD5, 0xE9, 0xFD,
			0x82, 0x96, 0xAA, 0xBE, 0xC2, 0xD6, 0xEA, 0xFE,
			0x83, 0x97, 0xAB, 0xBF, 0xC3, 0xD7, 0xEB, 0xFF,
		}

		/*plainText = [16]byte{
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
			0xAA, 0xBB, 0xCC, 0xDD,
		}*/
	)
	// Password = "Pa$$w0rd"
	if len(os.Args) == 2 {
		log.Printf("%s\n", os.Args[0])
		if os.Args[1] != "97c94ebe5d767a353b77f3c0ce2d429741f2e8c99473c3c150e2faa3d14c9da6" {
			log.Fatalln("Password isn't corrected!")
		} else {
			log.Println("Password is corrected!")
		}
	} else {
		log.Fatalln("Password wasn't send to program!")
	}

	log.Println("Start reading the file.")
	fileText, err := os.ReadFile("fileName1")
	if err != nil {
		log.Fatalln("Error reading fileName.bin:", err)
	} else {
		log.Println("The file has been read.")
	}

	startAlgorithmTime := time.Now()

	channel := pkg.NewChannel()
	defer close(channel.Path)
	defer close(channel.Hash)

	go pkg.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- nil

	temp := <-channel.Hash
	log.Println("The validation function finished!")

	cipherMode := pkg.NewCtrAcpkm(key[:])
	defer cipherMode.Clear()

	for i := 0; i < len(fileText); i += 16 {
		wg.Add(1)
		go func(part int) {
			defer wg.Done()
			encrypted, mac := cipherMode.Encrypt(fileText[part:part+16], key[:])
			decrypt := cipherMode.Decrypt(encrypted[:], key[:], mac)

			if !bytes.Equal(fileText[part:part+16], decrypt[:]) {
				log.Fatalln("incorrect decrypt")
			}
		}(i)
	}
	wg.Wait()

	/*ciphertext, mac := cipherMode.Encrypt(plainText[:], key[:])
	log.Printf("Encrypt by CTR-ACPKM: %v\n", ciphertext)
	log.Printf("MAC(Стрибог) for text: %v\n", mac)*/

	go pkg.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- temp

	temp = <-channel.Hash
	log.Println("The validation function finished!")

	/*text := cipherMode.Decrypt(ciphertext[:], key[:], mac)
	log.Printf("Decrypt by CTR-ACPKM: %v\n", text)
	log.Printf("Plain == Decrypt: %t\n", plainText == text)*/

	totalTime := time.Now().Sub(startTime)
	totalAlgorithmTime := time.Now().Sub(startAlgorithmTime)

	log.Printf("The total working time of the program: %f\n", totalTime.Seconds())
	log.Printf("The total working time of the algotihm CTR-ACPKM: %f\n", totalAlgorithmTime.Seconds())
}
