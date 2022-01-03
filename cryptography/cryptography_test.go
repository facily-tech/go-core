package cryptography

import (
	"testing"

	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// key required 16 digits, nonce required 12 digits
func TestEncryptDecrypt(t *testing.T) {
	telefone := "999999999"
	crip := NewCryptography([]byte("XXXXXfacilyXXXXX"), []byte("XXXfacilyXXX"))
	encrypt, err := crip.Encrypt(telefone)
	assert.Nil(t, err)
	assert.NotEmpty(t, encrypt)

	crip2 := NewCryptography([]byte("XXXXXfacilyXXXXX"), []byte("XXXfacilyXXX"))
	encrypt2, err := crip2.Encrypt(telefone)
	assert.Nil(t, err)
	assert.NotEmpty(t, encrypt2)

	assert.Equal(t, encrypt, encrypt2)

	decrypt, err := crip.Decrypt(encrypt)
	assert.Nil(t, err)

	assert.Equal(t, telefone, decrypt)
}
