package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/tcolgate/gostikkit/evpkdf"
)

func main() {
	b, _ := m()
	md := hex.EncodeToString(b)
	fmt.Printf("md5: %s\n", md)

	fmt.Printf("md5 base64: %s\n", base64.StdEncoding.EncodeToString(b))

	ciphertext := encryptCBC("scnace", md)

	fmt.Printf("ciphertext: %s\n", base64.StdEncoding.EncodeToString(ciphertext))

	// fmt.Printf("iv: %s\n", base64.StdEncoding.EncodeToString(iv))

}

func m() ([]byte, error) {
	hasher := md5.New()
	if _, err := hasher.Write([]byte("hi")); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

var opensslmagic = []byte{0x53, 0x61, 0x6c, 0x74, 0x65, 0x64, 0x5f, 0x5f}

func addSalt(ciphertext, salt []byte) []byte {
	if len(salt) == 0 {
		return ciphertext
	}
	return append(append(opensslmagic, salt...), ciphertext...)
}

func encryptCBC(plaintext string, passphares string) []byte {

	salt := genChars(8)

	keylen := 32
	key := make([]byte, keylen)
	ivlen := aes.BlockSize
	iv := make([]byte, ivlen)

	keymat := evpkdf.New(md5.New, []byte(passphares), salt, keylen+ivlen, 1)
	keymatbuf := bytes.NewReader(keymat)

	n, err := keymatbuf.Read(key)
	if n != keylen || err != nil {
		panic("keymaterial was short reading key")
	}

	n, err = keymatbuf.Read(iv)
	if n != ivlen || err != nil {
		panic("keymaterial was short reading iv")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	padded, _ := zeroPadding([]byte(plaintext), block.BlockSize())
	ciphertext := make([]byte, len(padded))

	if _, err := crand.Read(iv); err != nil {
		panic(err)
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, padded)

	fmt.Printf("cipher: %s\n", base64.StdEncoding.EncodeToString(ciphertext))
	fmt.Printf("iv: %v\n", base64.StdEncoding.EncodeToString(iv))
	fmt.Printf("key: %v\n", base64.StdEncoding.EncodeToString(key))

	return addSalt(ciphertext, salt)
}

func zeroPadding(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := uint8(1)
	for ((len(data) + int(padlen)) % blocklen) != 0 {
		padlen++
	}

	if int(padlen) > blocklen {
		panic(fmt.Sprintf("generated invalid padding length %v for block length %v", padlen, blocklen))
	}
	pad := bytes.Repeat([]byte{byte(0)}, int(padlen))
	return append(data, pad...), nil
}

func pkcs7Pad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := uint8(1)
	for ((len(data) + int(padlen)) % blocklen) != 0 {
		padlen++
	}

	if int(padlen) > blocklen {
		panic(fmt.Sprintf("generated invalid padding length %v for block length %v", padlen, blocklen))
	}
	pad := bytes.Repeat([]byte{byte(padlen)}, int(padlen))
	return append(data, pad...), nil
}

var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func genChars(n int) []byte {
	out := make([]byte, n)
	rs := make([]byte, n)
	crand.Read(rs)
	for i := 0; i < n; i++ {
		out[i] = chars[uint(rs[i])%uint(len(chars))]
	}
	return out
}
