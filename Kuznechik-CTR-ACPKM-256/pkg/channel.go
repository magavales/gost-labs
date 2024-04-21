package pkg

import "os"

type Channel struct {
	File chan *os.File
	Hash chan []byte
}

func NewChannel() *Channel {
	return &Channel{}
}
