package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"

	"github.com/zackradisic/youtube-rooms/internal/models"
	"github.com/zackradisic/youtube-rooms/internal/ws"
)

func (s *Server) setupRoutes() {
	apiRouter := s.router.PathPrefix("/api/").Subrouter()
	s.addRoute(apiRouter, "GET", "/auth/discord/", s.handleBeginAuth())
	s.addRoute(apiRouter, "GET", "/auth/discord/callback", s.handleCompleteAuth())

	s.addRoute(s.router, "GET", "/ws", s.handleWS())

	s.router.PathPrefix("/").HandlerFunc(s.handleNonAPIRoute())
}

func (s *Server) addRoute(router *mux.Router, method string, path string, handler func(http.ResponseWriter, *http.Request)) {
	router.HandleFunc(path, handler).Methods(method)
}

func (s *Server) handleWS() http.HandlerFunc {
	type requestBody struct {
		roomName string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		rBody := &requestBody{}

		if err := decoder.Decode(rBody); err != nil {
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if rBody.roomName != "" {
			s.respondError(w, "Invalid room name", http.StatusBadRequest)
			return
		}

		ws.Serve(rBody.roomName, s.Hub, w, r)
	}
}

func (s *Server) handleCompleteAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.respondError(w, err.Error(), 403)
			return
		}

		state := session.Values["state"]
		if state == "" {
			s.respondError(w, "Unauthorized", 403)
			return
		}

		queryParams := r.URL.Query()
		if stateValues, ok := queryParams["state"]; ok {
			exists := false

			for _, val := range stateValues {
				if val == state {
					exists = true
					break
				}
			}

			if !exists {
				s.respondError(w, "Unauthorized", 403)
				return
			}

			if codeValues, ok := queryParams["code"]; ok {
				code := codeValues[0]
				authToken, err := s.getAuthorizationCode(code)
				fmt.Println(authToken)
				if err != nil {
					log.Fatal(err)
				}

				userInfo, err := s.getDiscordUserInfo(authToken.AccessToken, authToken.RefreshToken)
				if err != nil {
					s.respondError(w, "There was an error getting your Discord info", 500)
					return
				}

				user := &models.User{}
				if err := s.DB.Where(&models.User{DiscordID: userInfo.ID}).First(&user).Error; err != nil {
					if gorm.IsRecordNotFoundError(err) {
						user.LastDiscordUsername = userInfo.Username
						user.LastDiscordDiscriminator = userInfo.Discriminator
						if err = s.DB.Create(user).Error; err != nil {
							s.respondError(w, "Error accessing your user info", http.StatusInternalServerError)
							return
						}
					} else {
						s.respondError(w, "Error accessing your user info", http.StatusInternalServerError)
						return
					}
				}

				fmt.Println(user.ID, user.LastDiscordUsername, user.LastDiscordDiscriminator)
				return
			}
		}

		s.respondError(w, "Unauthorized", 403)
	}
}

func (s *Server) handleBeginAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.respondError(w, err.Error(), 403)
			return
		}

		state := base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))

		session.Values["state"] = state
		session.Save(r, w)

		param := url.Values{}
		param.Add("state", state)

		url := s.authDetails.OAuthURL.String() + "&" + param.Encode()
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleNonAPIRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/dist/index.html")
	}
}
