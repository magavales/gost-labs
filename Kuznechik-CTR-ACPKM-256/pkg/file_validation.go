package pkg

import (
	"crypto/md5"
	"io"
	"log"
	"os"
)

func Validation(channel *Channel) {
	path := <-channel.Path
	expectedHash := <-channel.Hash

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	hash := h.Sum(nil)
	log.Printf("%s", hash)
	h = nil

	if expectedHash == nil {
		channel.Hash <- hash
		log.Println("Program gets the hash in the buffer.")
	} else {
		if string(expectedHash) == string(hash) {
			channel.Hash <- hash
			log.Println("File is correct.")
		} else {
			log.Fatalf("The integrity of the file has been violated!")
		}
	}
}
