package web

import (
	"math/rand"
	"os"
	"os/exec"
	"syscall"
)

//Shell exec shell command
func Shell(cmd string, args ...string) error {
	bin, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}
	return syscall.Exec(bin, append([]string{cmd}, args...), os.Environ())
}

//RandomStr randome string
func RandomStr(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	buf := make([]rune, n)
	for i := range buf {
		buf[i] = letters[rand.Intn(len(letters))]
	}
	return string(buf)
}

//RandomBytes random bytes
func RandomBytes(n uint) ([]byte, error) {
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	return buf, err
}
