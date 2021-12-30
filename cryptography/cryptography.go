package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

// ICryptography is Cryptography interface to use when you wanna make a default setup on its methods.
type ICryptography interface {
	/* Encrypt returns a string with the plainText encrypted
	   if "plainText" was a error then it will return a empty string and a error */
	Encrypt(string) (string, error)
	/* Decrypt returns a string with the ciphertext decrypted
	   if "ciphertext" was a error then it will return a empty string and a error */
	Decrypt(string) (string, error)
}

type Cryptography struct {
	key   []byte
	nonce []byte
}

// NewCryptography returns a new Cryptography struct
func NewCryptography(key []byte, nonce []byte) *Cryptography {

	return &Cryptography{key: key, nonce: nonce}
}

// Encrypt will encrypt a plaintext.
func (s *Cryptography) Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ciphertext := aesgcm.Seal(nil, s.nonce, []byte(plainText), nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt will decrypt a ciphertext previously encrypted .
func (s *Cryptography) Decrypt(ciphertext string) (string, error) {
	plainText, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	decrypt, err := aesgcm.Open(nil, s.nonce, plainText, nil)
	if err != nil {
		return "", err
	}

	return string(decrypt), nil
}
