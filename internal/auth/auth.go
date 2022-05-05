package auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
)

func MakeMD5(in string) string {
	binHash := md5.Sum([]byte(in))
	return hex.EncodeToString(binHash[:])
}

func ComparePasswords(password string, original string) bool {
	return MakeMD5(password) == original
}

func GenerateAuthToken() string {
	b := make([]byte, 24)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func SetAuthCookie(c *gin.Context, token string) {
	c.SetCookie("auth", token, 3600, "/", "", false, false)
}
