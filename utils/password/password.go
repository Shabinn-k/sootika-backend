package password

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	return string(hash),err
}

func ComparePassword(hash,password string)bool{
	err:=bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err==nil
}

func HashToken(token string)string{
	hashtoken:=sha256.Sum256([]byte(token))
	return hex.EncodeToString(hashtoken[:])
}

func CompareHashToken(token string,storedHash string)bool{
	hash:=HashToken(token)
	return subtle.ConstantTimeCompare(
		[]byte(hash),
		[]byte(storedHash),
	)==1
}