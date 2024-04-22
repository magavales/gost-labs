package pkg

import (
	"crypto/md5"
	"io"
	"log"
)

func Validation(channel *Channel) {
	h := md5.New()
	file := <-channel.File
	expectedHash := <-channel.Hash
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	hash := h.Sum(nil)
	log.Printf("%s", hash)

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
