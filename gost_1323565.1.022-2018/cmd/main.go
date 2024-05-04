package main

import (
	"gost_1323565.1.022-2018/pkg"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	startTime := time.Now()
	var (
		key = []byte{0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79,
			0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f, 0x80, 0x81, 0x82, 0x83}
		T = []byte{
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		}
		P = []byte{0x6b, 0x75, 0x7a, 0x6e, 0x65, 0x63, 0x68, 0x69, 0x6b, 0x0d, 0x0a}
		U = []byte{0x41, 0x6c, 0x69, 0x63, 0x65, 0x0d, 0x0a, 0x0d, 0x0a}
		A = []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	)

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

	channel := pkg.NewChannel()

	go pkg.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- nil

	temp := <-channel.Hash
	log.Println("The validation function finished!")

	kdf := pkg.NewKDF(key, T, P, U, A)
	defer kdf.Clear()
	clearData(&key)
	clearData(&T)
	clearData(&P)
	clearData(&U)
	clearData(&A)
	result := kdf.Generate()

	go pkg.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- temp

	temp = []byte{0x23}

	_ = <-channel.Hash
	log.Println("The validation function finished!")

	totalTime := time.Now().Sub(startTime)

	log.Printf("Result: %s", string(result))
	log.Printf("The total working time of the program: %f\n", totalTime.Seconds())

	close(channel.Path)
	close(channel.Hash)
}

func clearData(data *[]byte) {
	for i := 0; i < len(*data); i++ {
		(*data)[i] = byte(rand.Intn(10000))
	}
}
