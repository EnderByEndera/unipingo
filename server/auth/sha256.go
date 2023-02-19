package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSha256(str string) string {
	srcByte := []byte(str)
	sha256New := sha256.New()
	sha256Bytes := sha256New.Sum(srcByte)
	sha256String := hex.EncodeToString(sha256Bytes)
	return sha256String
}

func EncryptPassword(passwordInPlainText string) string {
	return GetSha256(passwordInPlainText)
}

func ComparePassword(passwordInPlainText, passwordSHA256 string) bool {
	return GetSha256(passwordInPlainText) == passwordSHA256
}
