package app

import (
	"Kuznechik-CTR-ACPKM-256/pkg"
	"bytes"
	"log"
	"os"
	"sync"
	"time"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
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

	cipherMode := pkg.NewCtrAcpkm(nil, key[:])
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

func (a *App) Test() {
	startTime := time.Now()
	var (
		wg  sync.WaitGroup
		key = [32]byte{
			0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
			0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
			0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
			0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		}

		plainText = [64]byte{
			0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
			0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
			0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
			0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xee, 0xff, 0x0a,
			0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
			0x99, 0xaa, 0xbb, 0xcc, 0xee, 0xff, 0x0a, 0x00,
			0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
			0xaa, 0xbb, 0xcc, 0xee, 0xff, 0x0a, 0x00, 0x11,
		}
		vector = []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0xd0, 0xb0, 0x62, 0x63, 0x65, 0x66}
	)

	startAlgorithmTime := time.Now()

	channel := pkg.NewChannel()
	defer close(channel.Path)
	defer close(channel.Hash)

	go pkg.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- nil

	temp := <-channel.Hash
	log.Println("The validation function finished!")

	cipherMode := pkg.NewCtrAcpkm(vector, key[:])
	defer cipherMode.Clear()

	for i := 0; i < len(plainText); i += 16 {
		wg.Add(1)
		go func(part int) {
			defer wg.Done()
			encrypted, mac := cipherMode.Encrypt(plainText[part:part+16], key[:])
			decrypt := cipherMode.Decrypt(encrypted[:], key[:], mac)

			if !bytes.Equal(plainText[part:part+16], decrypt[:]) {
				log.Fatalln("incorrect decrypt")
			}
		}(i)
	}
	wg.Wait()

	go pkg.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- temp

	temp = <-channel.Hash
	log.Println("The validation function finished!")

	totalTime := time.Now().Sub(startTime)
	totalAlgorithmTime := time.Now().Sub(startAlgorithmTime)

	log.Printf("The total working time of the program: %f\n", totalTime.Seconds())
	log.Printf("The total working time of the algotihm CTR-ACPKM: %f\n", totalAlgorithmTime.Seconds())
}
