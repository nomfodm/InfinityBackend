package utils

import (
	"crypto/md5"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashStringToBcrypt(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 5)
	return string(bytes), err
}

func VerifyStringHash(hash string, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}

func MD5ToString(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func AddHyphenToUUID(uuidRaw string) (uuidHyphened string) {
	return uuid.MustParse(uuidRaw).String()
}
