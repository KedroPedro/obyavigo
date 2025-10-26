package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type AESMail struct {
	block cipher.Block
}

func NewAESMail(b cipher.Block) *AESMail {
	return &AESMail{
		block: b,
	}
}

func (m *AESMail) EncryptMail(mail string) (string, error) {
	ciphertext := make([]byte, aes.BlockSize+len(mail))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(m.block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(mail))

	return hex.EncodeToString(ciphertext), nil

}

func (m *AESMail) DecryptMail(encrypted string) (string, error) {

	ciphertext, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(m.block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
