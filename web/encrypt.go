package web

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
)

//Encryptor encryptor
type Encryptor struct {
	Cip cipher.Block `inject:"aes.cip"`
	Key []byte       `inject:"hmac.key"`
}

//Encrypt encrypt buffer
func (p *Encryptor) Encrypt(buf []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(p.Cip, iv)
	val := make([]byte, len(buf))
	cfb.XORKeyStream(val, buf)

	return append(val, iv...), nil
}

//Decrypt decrypt buffer
func (p *Encryptor) Decrypt(buf []byte) ([]byte, error) {
	bln := len(buf)
	cln := bln - aes.BlockSize
	ct := buf[0:cln]
	iv := buf[cln:bln]

	cfb := cipher.NewCFBDecrypter(p.Cip, iv)
	val := make([]byte, cln)
	cfb.XORKeyStream(val, ct)
	return val, nil
}

//Sum sum a password
func (p *Encryptor) Sum(plain []byte) []byte {
	mac := hmac.New(sha512.New, p.Key)
	return mac.Sum(plain)
}

//Chk check password
func (p *Encryptor) Chk(plain, code []byte) bool {
	return hmac.Equal(p.Sum(plain), code)
}
