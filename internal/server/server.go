package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/matthewhartstonge/argon2"

	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/zackradisic/youtube-rooms/internal/room"
	"github.com/zackradisic/youtube-rooms/internal/ws"

	"github.com/gorilla/mux"

	"golang.org/x/time/rate"
)

// Server is the HTTP server object
type Server struct {
	router       *mux.Router
	sessionStore *sessions.CookieStore
	authDetails  *authDetails
	DB           *gorm.DB
	Hub          *ws.Hub
	RoomManager  *room.Manager
	argon2       *argon2.Config
	limiter      *rate.Limiter
}

// NewServer creates a server
func NewServer() *Server {
	a := argon2.DefaultConfig()
	s := &Server{
		router:       mux.NewRouter().StrictSlash(true),
		sessionStore: sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET"))),
		argon2:       &a,
		limiter:      rate.NewLimiter(50, 20),
	}
	s.setupRoutes()
	s.setupAuth()

	db, err := s.setupDB()
	if err != nil {
		panic("Error initializing DB: " + err.Error())
	}

	s.DB = db
	s.RoomManager = room.NewManager(db)
	s.Hub = ws.NewHub(s.RoomManager)
	go s.Hub.Run()

	return s
}

// Run runs the HTTP server
func (s *Server) Run(host string) {
	fmt.Println("Running server on " + host)
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"})
	// mux.CORSMethodMiddleware(s.router)
	log.Fatal(http.ListenAndServe(host, handlers.CORS(originsOk, methodsOk, headersOk)(s.router)))
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
