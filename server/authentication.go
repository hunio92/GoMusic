package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) authentication(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form["email"]
	pass := r.Form["password"]

	hash, err := bcrypt.GenerateFromPassword([]byte(pass[0]), bcrypt.DefaultCost)
	if err != nil {
		log.Println("hashing password error: ", err)
	}

	if email[0] == "asd@asd.com" {
		if err := bcrypt.CompareHashAndPassword(hash, []byte(pass[0])); err == nil {
			uuid, err := encrypt(email[0])
			if err != nil {
				fmt.Printf("Encryp: %d", err)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			fmt.Println("uuid: ", uuid)
			session, _ := s.CookieStore.Get(r, "sessionCookie")
			session.Values["uuid"] = uuid
			session.Save(r, w)
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) isUserLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.CookieStore.Get(r, "sessionCookie")
		if _, ok := session.Values["uuid"]; !ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		uuid := session.Values["uuid"].(string)
		_, err := decrypt(uuid)
		if err != nil {
			fmt.Printf("Decode: %v", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func encrypt(data string) (string, error) {
	key, err := aes.NewCipher([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("Could not get key: %v", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("Could not generate init vector(IV): %v", err)
	}

	stream := cipher.NewCFBEncrypter(key, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(data string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(data)
	key, err := aes.NewCipher([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("Could not get key: %v", err)
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("key to short: %v", err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(key, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
