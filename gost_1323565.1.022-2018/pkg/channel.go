package pkg

type Channel struct {
	Path chan string
	Hash chan []byte
}

func NewChannel() *Channel {
	return &Channel{
		Path: make(chan string),
		Hash: make(chan []byte),
	}
}
