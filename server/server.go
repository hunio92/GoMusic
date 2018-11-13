package server

import (
	"GoMusic/database"
	"fmt"
	"github.com/gorilla/sessions"
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
	Mux      *chi.Mux
	Sessions  *sessions.CookieStore
}

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*.gohtml"))
	
}

func NewServer(db *Database.MySqlDB) *Server {
	store := sessions.NewCookieStore([]byte("COOKIE_KEY"))
	s := &Server{
		Database: db,
		Mux:      chi.NewRouter(),
		Sessions: store,
	}
 
	s.configRoutes()
	
	return s
}

func (s *Server) configRoutes() {
	// Public Routes
	s.Mux.HandleFunc("/", s.renderLogin)
	s.Mux.Post("/login", s.login)
	
	// Private Routes
	s.Mux.Route("/home", func(homeRouter chi.Router) {
		homeRouter.Use(s.authentication)
		homeRouter.HandleFunc("/", s.renderHome)
	})
	
	// Static files
	workDir, _ := os.Getwd()
	fileServer(s.Mux, "/static", http.Dir(filepath.Join(workDir, "static")))
}

func (s *Server) renderLogin(w http.ResponseWriter, r *http.Request) {
	session, err := s.Sessions.Get(r, "cookie-name")
	if err != nil {
		fmt.Println("Could not get session cookie")
	}
	
	var Message struct {
		Error string
	}
	if flashes := session.Flashes(); len(flashes) > 0 {
		flash := flashes[0].(string)
		Message.Error = flash
	}
	
	fmt.Println("error: ", Message.Error)
	err = tpl.ExecuteTemplate(w, "login.gohtml", Message)
	if err != nil {
		fmt.Printf("Error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) renderHome(w http.ResponseWriter, r *http.Request) {
	session, _ := s.Sessions.Get(r, "cookie-name")
	session.AddFlash("")
	session.Save(r, w)
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
