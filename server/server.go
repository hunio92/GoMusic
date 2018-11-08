package server

import (
	"GoMusic/database"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/go-chi/chi"
)

var tpl *template.Template

type Server struct {
	Database *Database.MySqlDB
	Mux *chi.Mux
}

type Middleware func(http.Handler) http.Handler

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*.gohtml"))
}

func NewServer(db *Database.MySqlDB) *Server {
	s := &Server{
		Database: db,
		Mux: chi.NewRouter(),
	}
 
	s.configRoutes()

	return s
}

func (s *Server) configRoutes() {
	// Middlewares
	s.Mux.Use(Authentication)
	// Routes
	s.Mux.Get("/", s.handleLogin)
	
	s.Mux.Route("/home", func(homeRouter chi.Router) {
		homeRouter.Get("/home", s.handleHome)
	})
	
 
	// Static files
	workDir, _ := os.Getwd()
	fileServer(s.Mux, "/static", http.Dir(filepath.Join(workDir, "static")))
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
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
