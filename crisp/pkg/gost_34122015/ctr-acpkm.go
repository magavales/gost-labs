package pkg

import (
	"crisp/pkg/gost_34122015/streebog"
	"crypto/rand"
	"fmt"
	"log"
	"runtime"
)

type CtrAcpkm struct {
	iv [BlockSize]byte
}

func NewCtrAcpkm() *CtrAcpkm {
	var (
		iv [BlockSize]byte
	)

	_, err := rand.Read(iv[:])
	if err != nil {
		log.Fatalf("%s", err)
		return nil
	}

	return &CtrAcpkm{iv}
}

func (mode *CtrAcpkm) Encrypt(plaintext, key []byte) ([BlockSize]byte, []byte) {
	var (
		ciphertext [BlockSize]byte
		mac        []byte
		gamma      []byte
	)

	gamma = initGamma(mode.iv[:], key)

	X(ciphertext[:], plaintext, gamma)
	mac = createVerificationCode(ciphertext[:], key)

	return ciphertext, mac
}

func (mode *CtrAcpkm) Decrypt(ciphertext, key, mac []byte) [BlockSize]byte {
	var (
		plaintext   [BlockSize]byte
		expectedMac []byte
		gamma       []byte
	)

	gamma = initGamma(mode.iv[:], key)

	expectedMac = createVerificationCode(ciphertext[:], key)
	if !verifyVerificationCode(expectedMac, mac) {
		log.Fatalf("Expected MAC isn't equal to received MAC")
	}

	X(plaintext[:], ciphertext, gamma)

	return plaintext
}

func (mode *CtrAcpkm) Clear() {
	for i := 0; i < len(mode.iv); i++ {
		mode.iv[i] = 0x00
	}
	runtime.GC()
	fmt.Printf("Clear mem [CtrAcpkm]: %p\n", &mode)
}

func initGamma(initialVector, key []byte) []byte {
	var (
		gamma []byte
	)

	cipher := NewCipher(key[:])
	defer cipher.Clear()
	encoded := cipher.Encrypt(initialVector[:])

	gamma = append(gamma, encoded[:]...)

	return gamma
}

func createVerificationCode(ciphertext, key []byte) []byte {
	hash := streebog.NewHash()
	hash.Write(ciphertext)
	mac := hash.Sum(nil)
	hash.Reset()

	return mac
}

func verifyVerificationCode(expectedMac, mac []byte) bool {
	if string(expectedMac) != string(mac) {
		log.Fatalf("Expected MAC isn't equal to received MAC")
		return false
	}

	return true
}
