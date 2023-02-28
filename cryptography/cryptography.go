/*
Package cryptography was made to encrypt sensitive data.
*/
package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/pkg/errors"
)

// ICryptography is Cryptography interface to use when you wanna make a default setup on its methods.
//
//nolint:lll // go generate line
//go:generate mockgen -source=cryptography.go -destination=cryptography_mock.go -package=cryptography
type ICryptography interface {
	/* Encrypt returns a string with the plainText encrypted
	   if "plainText" was a error then it will return a empty string and a error */
	Encrypt(string) (string, error)
	/* Decrypt returns a string with the ciphertext decrypted
	   if "ciphertext" was a error then it will return a empty string and a error */
	Decrypt(string) (string, error)
}

// Cryptography struct.
type Cryptography struct {
	key   []byte
	nonce []byte
}

// NewCryptography returns a new Cryptography struct.
func NewCryptography(key []byte, nonce []byte) *Cryptography {
	return &Cryptography{key: key, nonce: nonce}
}

// Encrypt will encrypt a plaintext.
func (s *Cryptography) Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", errors.Wrap(err, "can't initialize NewCipher")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "can't initialize NewGCM")
	}
	ciphertext := aesgcm.Seal(nil, s.nonce, []byte(plainText), nil)
	ciphertextStr := hex.EncodeToString(ciphertext)

	return ciphertextStr, nil
}

// Decrypt will decrypt a ciphertext previously encrypted.
func (s *Cryptography) Decrypt(ciphertext string) (string, error) {
	plainText, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", errors.Wrap(err, "error on decodeString to a byte value")
	}
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", errors.Wrap(err, "can't initialize NewCipher")
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "can't initialize NewGCM")
	}
	decrypt, err := aesgcm.Open(nil, s.nonce, plainText, nil)
	if err != nil {
		return "", errors.Wrap(err, "can't decrypt ciphertext")
	}

	return string(decrypt), nil
}
