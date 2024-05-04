package pkg

import (
	"Kuznechik-CTR-ACPKM-256/pkg/streebog"
	"crypto/rand"
	"log"
)

type CtrAcpkm struct {
	Gamma []byte
}

func NewCtrAcpkm(vector, key []byte) *CtrAcpkm {
	var (
		iv    [BlockSize]byte
		gamma []byte
	)

	if vector == nil {
		_, err := rand.Read(iv[:])
		if err != nil {
			log.Fatalf("%s", err)
			return nil
		}

		gamma = initGamma(iv[:], key)
	} else {
		gamma = initGamma(vector, key)
	}

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

func (mode *CtrAcpkm) Clear() {
	mode.Gamma = nil
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
