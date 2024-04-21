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

	if expectedHash == nil || string(expectedHash) == string(hash) {
		channel.Hash <- hash
		log.Println("File is correct.")
	} else {
		log.Fatalf("The integrity of the file has been violated!")
	}
}
