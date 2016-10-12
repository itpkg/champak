package web_test

import (
	"crypto/aes"
	"testing"

	"github.com/itpkg/champak/web"
)

func TestEncryptor(t *testing.T) {
	cip, err := aes.NewCipher([]byte("1234567890123456"))
	if err != nil {
		t.Fatal(err)
	}
	en := web.Encryptor{Cip: cip, Key: []byte("123456")}

	hello := "Hello, Champak!"
	code, err := en.Encrypt([]byte(hello))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(code))
	plain, err := en.Decrypt(code)
	if err != nil {
		t.Fatal(err)
	}
	if string(plain) != hello {
		t.Fatalf("wang %s get %s", hello, string(plain))
	}

	code = en.Sum([]byte(hello))
	t.Log(string(code))
	if !en.Chk([]byte(hello), code) {
		t.Fatalf("check password failed")
	}

}
