package crypto

// This code is taken from https://github.com/Luzifer/go-openssl and adapted to work
// with aes-128-cbc instead of a 256 bytes key.
// There is also this blogpost https://dequeue.blogspot.ch/2014/11/decrypting-something-encrypted-with.html 
// that was very interesting.
//
// Those methods are compatible with the openssl commandline. To encrypt, use:
// 
//     openssl enc -aes-128-cbc -pass pass:PASSWORD -salt -base64
//
// To decrypt, put the content into a file and use:
//
//     openssl enc -aes-128-cbc -d -a -pass pass:PASSWORD -in INPUT.txt
//

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
)

// OpenSSL is a helper to generate OpenSSL compatible encryption
// with autmatic IV derivation and storage. As long as the key is known all
// data can also get decrypted using OpenSSL CLI.
// Code from http://dequeue.blogspot.de/2014/11/decrypting-something-encrypted-with.html


// OpenSSL salt is always this string + 8 bytes of actual salt
const openSSLSaltHeader =  "Salted__"

type openSSLCreds struct {
	key []byte
	iv  []byte
}


func DecryptFile(passphrase, filepath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return DecryptString(passphrase, string(dat))
}

// DecryptString decrypts a string that was encrypted using OpenSSL and AES-128-CBC
func DecryptString(passphrase, encryptedBase64String string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedBase64String)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("Data is too short")
	}
	saltHeader := data[:aes.BlockSize]
	if string(saltHeader[:8]) != openSSLSaltHeader {
		return nil, fmt.Errorf("Does not appear to have been encrypted with OpenSSL, salt header missing.")
	}
	salt := saltHeader[8:]
	creds, err := extractOpenSSLCreds([]byte(passphrase), salt)
	if err != nil {
		return nil, err
	}
	return decrypt(creds.key, creds.iv, data)
}

func decrypt(key, iv, data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("bad blocksize(%v), aes.BlockSize = %v\n", len(data), aes.BlockSize)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCDecrypter(c, iv)
	cbc.CryptBlocks(data[aes.BlockSize:], data[aes.BlockSize:])
	out, err := pkcs7Unpad(data[aes.BlockSize:], aes.BlockSize)
	if out == nil {
		return nil, err
	}
	return out, nil
}

func EncryptFile(passphrase, plaintextString string, filepath string) error {
	dat, err := EncryptString(passphrase, plaintextString)
	if err != nil {
		return err
	}
	// with openssl commandline, if in encrypted file after base64 line not new line you get error "error reading input file"
	return ioutil.WriteFile(filepath, []byte(string(dat) + "\n"), 0644)
}

// EncryptString encrypts a string in a manner compatible to OpenSSL encryption
// functions using AES-256-CBC as encryption algorithm
func EncryptString(passphrase, plaintextString string) ([]byte, error) {
	salt := make([]byte, 8) // Generate an 8 byte salt
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	data := make([]byte, len(plaintextString)+aes.BlockSize)
	copy(data[0:], openSSLSaltHeader)
	copy(data[8:], salt)
	copy(data[aes.BlockSize:], plaintextString)

	creds, err := extractOpenSSLCreds([]byte(passphrase), salt)
	if err != nil {
		return nil, err
	}

	enc, err := encrypt(creds.key, creds.iv, data)
	if err != nil {
		return nil, err
	}

	return []byte(base64.StdEncoding.EncodeToString(enc)), nil
}

func encrypt(key, iv, data []byte) ([]byte, error) {
	padded, err := pkcs7Pad(data, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCEncrypter(c, iv)
	cbc.CryptBlocks(padded[aes.BlockSize:], padded[aes.BlockSize:])

	return padded, nil
}

// openSSLEvpBytesToKey follows the OpenSSL (undocumented?) convention for extracting the key and IV from passphrase.
// It uses the EVP_BytesToKey() method which is basically:
// D_i = HASH^count(D_(i-1) || password || salt) where || denotes concatentaion, until there are sufficient bytes available
// 32 bytes since we're expecting to handle AES-128, 16 bytes for a key and 16 bytes for the IV
func extractOpenSSLCreds(password, salt []byte) (openSSLCreds, error) {
	m := make([]byte, 32)
	prev := []byte{}
	for i := 0; i < 3; i++ {
		prev = hash(prev, password, salt)
		copy(m[i*16:], prev)
	}
	return openSSLCreds{key: m[:16], iv: m[16:]}, nil
}

func hash(prev, password, salt []byte) []byte {
	a := make([]byte, len(prev)+len(password)+len(salt))
	copy(a, prev)
	copy(a[len(prev):], password)
	copy(a[len(prev)+len(password):], salt)
	return md5sum(a)
}

func md5sum(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

// pkcs7Pad appends padding.
func pkcs7Pad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(data, pad...), nil
}

// pkcs7Unpad returns slice of the original data without padding.
func pkcs7Unpad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > blocklen || padlen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padlen], nil
}
