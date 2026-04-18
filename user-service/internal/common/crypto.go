package common

import (
	"crypto/md5"
	"encoding/hex"
)

const Salt = "abcd1234"

func Md5(input string) string {
	h := md5.Sum([]byte(input))
	return hex.EncodeToString(h[:])
}

func EncryptPassword(password string) string {
	return Md5(password + Salt)
}
