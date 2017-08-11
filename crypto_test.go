package main 

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"io/ioutil"
	"os"
)

var o *OpenSSL = NewOpenSSL()

func TestDecryptString(t *testing.T){
	// the content has been generated with:
	// echo -n "it works perfectly" | openssl enc -aes-128-cbc -pass pass:essai -salt -base64
	result, err := o.DecryptString(`essai`, `U2FsdGVkX1+byhH87dJhyoozSyunGV1EQn8qi2hP74kbKHBleEiYXa3dAYy2LkmU`)
	assert.Nil(t, err, "decryption failed")
	assert.Equal(t, "it works perfectly", string(result), "content match.")
}

func TestDecryptFile(t *testing.T){	

	result, err := o.DecryptFile(`essai`, "resources/test.txt.enc")
	assert.Nil(t, err, "decryption failed")
	assert.Equal(t, "it works perfectly", string(result), "content match.")
}

func TestEncryptString(t *testing.T){
	result, err := o.EncryptString(`essai`, `it works perfectly`)
	assert.Nil(t, err, "encryption ok.")
	t.Log(string(result))
}

func TestEncryptFile(t *testing.T){	
	content := []byte("régénère-toi, !!")
	tmpfile, err := ioutil.TempFile("", "ossl-test")
	if err != nil {
		t.Error(err)
	}

	t.Logf("tmpfile is: %s\n", tmpfile.Name())
	defer os.Remove(tmpfile.Name()) // clean up

	if err := o.EncryptFile("essai", string(content), tmpfile.Name()); err != nil {
		t.Error(err)
	}

	result, err := o.DecryptFile(`essai`, tmpfile.Name())
	assert.Nil(t, err, "decryption failed")
	assert.Equal(t, content, result, "content match.")
}
