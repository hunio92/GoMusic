package server

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form["email"]
	pass := r.Form["password"]
	
	hash, err := bcrypt.GenerateFromPassword([]byte(pass[0]), bcrypt.DefaultCost)
	if err != nil {
		log.Println("hashing password error: ", err)
	}
	
	session, _ := s.Sessions.Get(r, "cookie-name")
	
	if email[0] == "asd@asd.com" {
		if err := bcrypt.CompareHashAndPassword(hash, []byte(pass[0])); err == nil {
			session.Values["email"] = email[0]
			session.AddFlash("")
			session.Save(r, w)
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}
	}
	
	session.AddFlash("Wrong user or password !")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Sessions.Get(r, "cookie-name")
		if _, ok := session.Values["email"]; !ok {
			
			session.AddFlash("Your cookie has expired !")
			session.Save(r, w)
			return
		}

		next.ServeHTTP(w, r)
	})
}