package apis

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

/*
	this gen.go use to generate authenticate username, password and so on.
*/
const (
	PASSWORD_LEN = 24
)

// GeneratePassword generate random 24 bit password string
func GeneratePassword() string {
	return RandomString(PASSWORD_LEN)
}

// RandomString generate random string
func RandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// SHA224String generate pass string to it's hash
func SHA224String(password string) string {
	hash := sha256.New224()
	hash.Write([]byte(password))
	val := hash.Sum(nil)
	str := ""
	for _, v := range val {
		str += fmt.Sprintf("%02x", v)
	}
	return str
}

// GeneratePassAndHash generate random password and it's hash
func GeneratePassAndHash() (pass string, hash string) {
	pass = GeneratePassword()
	hash = SHA224String(pass)
	return
}
