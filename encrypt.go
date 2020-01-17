package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

const (
	AESKey = "wrdvpnisthebest!"
	AESIV  = "wrdvpnisthebest!"
)

func Encrypt(text string) string {
	result := make([]byte, len(text))

	block, err := aes.NewCipher([]byte(AESKey))
	if err != nil {
		panic(err)
	}

	enc := cipher.NewCFBEncrypter(block, []byte(AESIV))
	enc.XORKeyStream(result, []byte(text))
	return hex.EncodeToString([]byte(AESIV)) + hex.EncodeToString(result)
}

func Decrypt(text string) string {
	iv, _ := hex.DecodeString(text[0:32])
	textB, _ := hex.DecodeString(text[32:])
	result := make([]byte, len(textB))

	block, err := aes.NewCipher([]byte(AESKey))
	if err != nil {
		panic(err)
	}

	dec := cipher.NewCFBDecrypter(block, iv)
	dec.XORKeyStream(result, textB)
	return string(result)
}
