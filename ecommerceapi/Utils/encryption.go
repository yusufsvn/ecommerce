package Utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordByte := []byte(password)
	hashedpassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)

	if err != nil {
		log.Fatal("parola hashlenirken hatayla karşılaşıldı")
	}

	return string(hashedpassword), err
}

func CheckPasswordHash(password, hash string) bool {
	var err error = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashRefreshToken(token string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(token))
	if err != nil {
		return "", fmt.Errorf("hash algoritmasında hata %w", err)
	}

	hashed_token := hash.Sum(nil)

	hashed_string := hex.EncodeToString(hashed_token)

	return hashed_string, nil
}

func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// URL'de sorun yaratmayacak karakter seti ile string'e çeviriyoruz
	return base64.URLEncoding.EncodeToString(b), nil
}
