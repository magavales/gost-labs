package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
	"xoroshift128+/pkg"
	"xoroshift128+/pkg/validation"
)

func main() {
	var (
		seed [32]byte
	)
	startTime := time.Now()

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

	nanosecond := time.Now().Nanosecond()
	binary.LittleEndian.PutUint64(seed[:], uint64(nanosecond))

	generator := pkg.New(seed)
	fmt.Printf("Next [0]: %d\n", generator.Next())
	fmt.Printf("Next [1]: %d\n", generator.Next())
	fmt.Printf("Next [2]: %d\n", generator.Next())
	fmt.Printf("Next [3]: %d\n", generator.Next())
	fmt.Printf("Next [4]: %d\n", generator.Next())
	fmt.Printf("Next [5]: %d\n", generator.Next())
	fmt.Printf("Next [6]: %d\n", generator.Next())
	fmt.Printf("Next [7]: %d\n", generator.Next())

	start := time.Now().UnixMilli()
	pkg.Create1MbFile(generator) // ~ 300 msec
	fmt.Printf("Elapsed [1mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	pkg.Create100MbFile(generator) // ~ 28000 msec (28 sec)
	fmt.Printf("Elapsed [100mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	pkg.Create1000MbFile(generator) // ~ 300000 msec (300 sec)
	fmt.Printf("Elapsed [1000mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	pkg.Create1000Values(generator) // ~ 5 msec
	fmt.Printf("Elapsed [1000val]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	pkg.Create10000Values(generator) // ~ 25 msec
	fmt.Printf("Elapsed [10000val]: %d msec.\n", time.Now().UnixMilli()-start)

	go validation.Validation(channel)
	channel.Path <- os.Args[0]
	channel.Hash <- temp

	temp = []byte{0x23}

	_ = <-channel.Hash
	log.Println("The validation function finished!")

	totalTime := time.Now().Sub(startTime)

	log.Printf("The total working time of the program: %f\n", totalTime.Seconds())

	close(channel.Path)
	close(channel.Hash)
}
