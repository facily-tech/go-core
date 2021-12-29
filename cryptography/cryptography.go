package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type CryptographyI interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type Cryptography struct {
	key   []byte
	nonce []byte
}

func NewCryptography(key []byte, nonce []byte) *Cryptography {

	return &Cryptography{key: key, nonce: nonce}
}

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

func (s *Cryptography) Decrypt(ciphertext string) (string, error) {
	plainText, _ := hex.DecodeString(ciphertext)
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	decrypt, err := aesgcm.Open(nil, s.nonce, []byte(plainText), nil)
	if err != nil {
		return "", err
	}

	return string(decrypt), nil
}
