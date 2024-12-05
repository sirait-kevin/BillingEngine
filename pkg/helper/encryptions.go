package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var encryptionKey = []byte("a very very very very secret key") // 32 bytes

// Encrypt encrypts plain text string into cipher text
func Encrypt(plainText string) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plainText))
	return ciphertext, nil
}

// Decrypt decrypts cipher text into plain text string
func Decrypt(cipherText []byte) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", err
	}
	iv := cipherText[:aes.BlockSize]
	ciphertext := cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// HashPassword hashes plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// GenerateSignature where key = is secret key)
func GenerateSignature(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	io.WriteString(h, data)

	return fmt.Sprintf("%x", string(h.Sum(nil)))
}
