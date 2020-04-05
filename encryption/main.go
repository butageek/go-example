package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func createHash(key string) string {
	md5Hash := md5.Sum([]byte(key))
	return hex.EncodeToString(md5Hash[:])
}

func encrypt(data []byte, passphrase string) []byte {
	// generate a new aes cipher using our 32 byte long key
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		panic(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(block)
	// if any error generating new GCM
	// handle them
	if err != nil {
		panic(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	// separate the nonce and the encrypted data
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return decrypt(data, passphrase)
}

func main() {
	fmt.Println("Starting the application...")
	ciphertext := encrypt([]byte("Hello World"), "password")
	fmt.Printf("Encrypted: %x\n", ciphertext)
	plaintext := decrypt(ciphertext, "password")
	fmt.Printf("Decrypted: %s\n", plaintext)
	encryptFile("sample.txt", []byte("Hello World"), "password1")
	fmt.Println(string(decryptFile("sample.txt", "password1")))
}
