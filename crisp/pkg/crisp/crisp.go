package crisp

import (
	gost132356510222018 "crisp/pkg/gost_132356510222018"
	gost34122015 "crisp/pkg/gost_34122015"
	xoroshiro256plus "crisp/pkg/xoroshiro256+"
	"encoding/binary"
	"log"
	"runtime"
)

const (
	BlockSize  = 16
	KeySize    = 32
	PacketSize = 56
)

type Crisp struct {
	Encoder    Encoder
	Decoder    Decoder
	randomSeed [32]byte
}

func NewCrisp(key []byte, randomSeed [32]byte) *Crisp {
	if len(key) != KeySize {
		log.Fatalln("Key size should be 32 bytes")
	}

	cipher := gost34122015.NewCtrAcpkm()
	kdf := gost132356510222018.NewKDF(key[:])
	return &Crisp{
		Decoder: Decoder{
			random: xoroshiro256plus.New(randomSeed),
			kdf:    kdf,
			cipher: cipher,
			seqNum: 0,
		},
		Encoder: Encoder{
			random: xoroshiro256plus.New(randomSeed),
			kdf:    kdf,
			cipher: cipher,
			seqNum: 0,
		},
		randomSeed: randomSeed,
	}
}

func (c *Crisp) Encode(plainText []byte) []Message {
	var res []Message

	c.Reset()
	for i := 0; i < len(plainText); i += BlockSize {
		message := c.EncodeNextBlock(plainText[i : i+BlockSize])
		res = append(res, message)
	}

	return res
}

func (c *Crisp) EncodeNextBlock(plainText []byte) Message {
	var (
		seqNum [4]byte
		seed   [8]byte
	)

	if len(plainText) != BlockSize {
		log.Fatalln("Block size should be 16 bytes")
	}
	e := c.Encoder
	binary.BigEndian.PutUint32(seqNum[:], e.seqNum)
	binary.BigEndian.PutUint64(seed[:], e.random.Next())

	e.kdf.Key = seed[:]
	key := e.kdf.GenerateKey()

	block := plainText[:]

	ciphertext, mac := e.cipher.Encrypt(block, key)

	var message []byte
	message = append(message, ExternalKeyIdFlagWithVersion...)
	message = append(message, CS...)
	message = append(message, KeyId...)
	message = append(message, seqNum[:]...)
	message = append(message, ciphertext[:]...)
	message = append(message, mac...)

	e.seqNum += 1 // complete current iteration and prepare next
	return Message{
		ExternalKeyIdFlagWithVersion: ExternalKeyIdFlagWithVersion,
		CS:                           CS,
		KeyId:                        KeyId,
		SeqNum:                       seqNum[:],
		Payload:                      ciphertext[:],
		ICV:                          mac[:],
		Digits:                       message,
	}
}

func (c *Crisp) Decode(cipherText [][]byte) [][]byte {
	for i, b := range cipherText {
		if len(b) != PacketSize {
			log.Fatalf("Block size of block [%d] should be 56 bytes", i)
		}
	}

	var res [][]byte
	for _, b := range cipherText {
		decoded := c.DecodeNextBlock(b)
		res = append(res, decoded)
	}

	return res
}

func (c *Crisp) DecodeNextBlock(cipherText []byte) []byte {
	if len(cipherText) != PacketSize {
		log.Fatalln("Block size should be equal 56 bytes")
	}
	d := c.Decoder

	payload := cipherText[8:24]
	mac := cipherText[24:56]

	key := d.kdf.GenerateKey()
	decrypt := d.cipher.Decrypt(payload, key, mac)

	return decrypt[:]
}

func (c *Crisp) Reset() {
	c.Encoder.seqNum = 0
	c.Encoder.random = xoroshiro256plus.New(c.randomSeed)
	c.Decoder.seqNum = 0
	c.Decoder.random = xoroshiro256plus.New(c.randomSeed)
}

func (c *Crisp) Clear() {
	for i := 0; i < len(c.randomSeed); i++ {
		c.randomSeed[i] = 0x00
	}
	c.Decoder.kdf.Clear()
	c.Decoder.cipher.Clear()
	c.Decoder.seqNum = 0
	c.Encoder.kdf.Clear()
	c.Encoder.cipher.Clear()
	c.Encoder.seqNum = 0
	runtime.GC()
	log.Printf("Clear mem [Crisp]: %p\n", &c)
}
