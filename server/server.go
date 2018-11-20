package server

import (
	"GoMusic/database"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/sessions"

	"github.com/go-chi/chi"
)

var tpl *template.Template

type Server struct {
	Database    *Database.Repo
	Mux         *chi.Mux
	CookieStore *sessions.CookieStore
}

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*.gohtml"))

}

func NewServer(db *Database.Repo) *Server {
	newCookie := sessions.NewCookieStore([]byte(os.Getenv("COOKIE_KEY")))
	s := &Server{
		Database:    db,
		Mux:         chi.NewRouter(),
		CookieStore: newCookie,
	}

	s.CookieStore.MaxAge(60) // cookie expires: 1 min
	s.configRoutes()

	return s
}

func (s *Server) configRoutes() {
	// Public Routes
	s.Mux.HandleFunc("/", s.renderLogin)
	s.Mux.Post("/login", s.authentication)

	// Private Routes
	s.Mux.Route("/home", func(homeRouter chi.Router) {
		homeRouter.Use(s.isUserLoggedIn)
		homeRouter.HandleFunc("/", s.renderHome)
	})

	// Static files
	workDir, _ := os.Getwd()
	fileServer(s.Mux, "/static", http.Dir(filepath.Join(workDir, "static")))
}

func (s *Server) renderLogin(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) renderHome(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "home.gohtml", nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
