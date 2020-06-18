package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

// Server is the HTTP server object
type Server struct {
	router       *mux.Router
	sessionStore *sessions.CookieStore
	authDetails  *authDetails
	DB           *gorm.DB
}

// NewServer creates a server
func NewServer() *Server {
	s := &Server{
		router:       mux.NewRouter().StrictSlash(true),
		sessionStore: sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET"))),
	}

	s.setupRoutes()
	s.setupAuth()

	return s
}

// Run runs the HTTP server
func (s *Server) Run(host string) {
	fmt.Println("Running server on " + host)
	log.Fatal(http.ListenAndServe(host, s.router))
}

func (s *Server) respondJSON(w http.ResponseWriter, payload interface{}, status int) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func (s *Server) respondError(w http.ResponseWriter, message string, status int) {
	s.respondJSON(w, map[string]string{"error": message}, status)
}
