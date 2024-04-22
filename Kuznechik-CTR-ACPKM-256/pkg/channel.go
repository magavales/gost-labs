package pkg

import "os"

type Channel struct {
	File chan *os.File
	Hash chan []byte
}

func NewChannel() *Channel {
	return &Channel{
		File: make(chan *os.File),
		Hash: make(chan []byte),
	}
}
