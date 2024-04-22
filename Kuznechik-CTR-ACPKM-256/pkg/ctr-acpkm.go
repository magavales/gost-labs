package pkg

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"log"
)

type CtrAcpkm struct {
	Gamma []byte
}

func NewCtrAcpkm(key []byte) *CtrAcpkm {
	var (
		iv    [BlockSize]byte
		gamma []byte
	)

	_, err := rand.Read(iv[:])
	if err != nil {
		log.Fatalf("%s", err)
		return nil
	}

	gamma = initGamma(iv[:], key)

	return &CtrAcpkm{gamma}
}

func (mode *CtrAcpkm) Encrypt(plaintext, key []byte) ([BlockSize]byte, []byte) {
	var (
		ciphertext [BlockSize]byte
		mac        []byte
	)

	X(ciphertext[:], plaintext, mode.Gamma)
	mac = createVerificationCode(ciphertext[:], key)

	return ciphertext, mac
}

func (mode *CtrAcpkm) Decrypt(ciphertext, key, mac []byte) [BlockSize]byte {
	var (
		plaintext   [BlockSize]byte
		expectedMac []byte
	)

	expectedMac = createVerificationCode(ciphertext[:], key)
	if !verifyVerificationCode(expectedMac, mac) {
		log.Fatalf("Expected MAC isn't equal to received MAC")
	}

	X(plaintext[:], ciphertext, mode.Gamma)

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
