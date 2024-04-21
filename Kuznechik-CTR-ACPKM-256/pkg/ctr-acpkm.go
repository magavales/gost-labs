package pkg

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"log"
)

type CtrAcpkm struct {
	InitialVector [BlockSize]byte
}

func NewCtrAcpkm() *CtrAcpkm {
	var iv [BlockSize]byte

	_, err := rand.Read(iv[:])
	if err != nil {
		log.Fatalf("%s", err)
		return nil
	}

	return &CtrAcpkm{iv}
}

func (mode *CtrAcpkm) Encrypt(plaintext, key []byte) ([BlockSize]byte, []byte) {
	var (
		gamma      []byte
		ciphertext [BlockSize]byte
		mac        []byte
	)

	gamma = initGamma(mode.InitialVector[:], key)

	X(ciphertext[:], plaintext, gamma)
	mac = createVerificationCode(ciphertext[:], key)

	return ciphertext, mac
}

func (mode *CtrAcpkm) Decrypt(ciphertext, key, mac []byte) [BlockSize]byte {
	var (
		gamma       []byte
		plaintext   [BlockSize]byte
		expectedMac []byte
	)

	expectedMac = createVerificationCode(ciphertext[:], key)
	if !verifyVerificationCode(expectedMac, mac) {
		log.Fatalf("Expected MAC isn't equal to received MAC")
	}

	gamma = initGamma(mode.InitialVector[:], key)

	X(plaintext[:], ciphertext, gamma)

	return plaintext
}

func initGamma(initialVector, key []byte) []byte {
	var (
		gamma []byte
	)

	cipher := NewCipher(key[:])
	encoded := cipher.Encrypt(initialVector[:])

	gamma = append(gamma, encoded[:]...)

	return gamma
}

func createVerificationCode(ciphertext, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(ciphertext)
	ciphertextMac := mac.Sum(nil)

	return ciphertextMac
}

func verifyVerificationCode(expectedMac, mac []byte) bool {
	if string(expectedMac) != string(mac) {
		log.Fatalf("Expected MAC isn't equal to received MAC")
		return false
	}

	return true
}
