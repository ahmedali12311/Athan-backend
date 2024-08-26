package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"app/config"
)

func Encrypt(value string) (*string, error) {
	byteValue := []byte(value)
	block, err := aes.NewCipher(config.CryptoKey)
	if err != nil {
		return nil, err
	}
	ciphered := make([]byte, aes.BlockSize+len(value))
	iv := ciphered[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphered[aes.BlockSize:], byteValue)
	encrypted := base64.URLEncoding.EncodeToString(ciphered)
	return &encrypted, nil
}

func Decrypt(value string) (*string, error) {
	ciphered, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(config.CryptoKey)
	if err != nil {
		return nil, err
	}
	if len(ciphered) < aes.BlockSize {
		return nil, errors.New("decrypting text is too short")
	}
	iv := ciphered[:aes.BlockSize]
	ciphered = ciphered[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphered, ciphered)
	decrypted := string(ciphered)
	return &decrypted, nil
}
