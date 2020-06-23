package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/matthewhartstonge/argon2"

	"github.com/zackradisic/youtube-rooms/internal/models"
	"github.com/zackradisic/youtube-rooms/internal/room"
	"github.com/zackradisic/youtube-rooms/internal/ws"
)

func (s *Server) setupRoutes() {
	apiRouter := s.router.PathPrefix("/api/").Subrouter()
	s.addRoute(apiRouter, "GET", "/auth/discord/", s.handleBeginAuth())
	s.addRoute(apiRouter, "GET", "/auth/discord/callback", s.handleCompleteAuth())

	s.addRoute(s.router, "GET", "/ws", s.checkAuthentication(s.handleWS()))
	s.addRoute(apiRouter, "GET", "/test", s.testRoute())

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/build/static"))))
	s.router.PathPrefix("/").HandlerFunc(s.handleNonAPIRoute())
}

func (s *Server) addRoute(router *mux.Router, method string, path string, handler func(http.ResponseWriter, *http.Request)) {
	router.HandleFunc(path, handler).Methods(method)
}

func (s *Server) testRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encoded, err := s.argon2.HashEncoded([]byte("test123"))
		if err != nil {
			fmt.Println(err)
		}
		rm := &models.Room{
			OwnerID:        1,
			HashedPassword: string(encoded),
			Name:           "zack's room",
		}
		if err := s.DB.Create(rm).Error; err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) handleWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		roomName, ok := params["roomName"]
		if !ok {
			s.respondError(w, "Bad Request", http.StatusBadRequest)
			return
		}

		roomPassword, ok := params["roomPassword"]
		if !ok {
			s.respondError(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if roomName[0] == "" {
			s.respondError(w, "Invalid room name", http.StatusBadRequest)
			return
		}

		rm, err := s.RoomManager.GetRoom(roomName[0])
		if err != nil {
			s.respondError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if rm.Model.HashedPassword != "" {
			if roomPassword[0] == "" {
				s.respondError(w, "Invalid room password", 403)
				return
			}

			ok, err := argon2.VerifyEncoded([]byte(roomPassword[0]), []byte(rm.Model.HashedPassword))
			if err != nil {
				s.respondError(w, "Internal server error", 500)
				return
			}

			if !ok {
				s.respondJSON(w, "Invalid room password", 403)
				return
			}
		}

		ctx := r.Context()
		if u, ok := ctx.Value(userKey).(*models.User); ok {
			if ua, ok := ctx.Value(userAuthKey).(*models.UserAuth); ok {
				user := room.NewUser(u, ua)
				rm.AddUser(user)
				ws.Serve(user, s.Hub, w, r)
				return
			}
		}

		s.respondError(w, "Internal server error", 500)
		return
	}
}

func (s *Server) handleCompleteAuth() http.HandlerFunc {
	type response struct {
		ID            string     `json:"id"`
		Username      string     `json:"username"`
		Discriminator string     `json:"discriminator"`
		Avatar        string     `json:"avatar"`
		Auth          *AuthToken `json:"auth"`
	}
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

				if err != nil {
					s.respondError(w, "Invalid credentials provided", 403)
					return
				}

				userInfo, err := s.getDiscordUserInfo(authToken.AccessToken, authToken.RefreshToken)
				if err != nil {
					s.respondError(w, "There was an error getting your Discord info", 500)
					return
				}

				user, err := s.createUser(userInfo, authToken)
				if err != nil {
					s.respondError(w, "There was an internal server error", 500)
					return
				}

				re := &response{}
				re.ID = user.DiscordID
				re.Username = user.LastDiscordUsername
				re.Discriminator = user.LastDiscordDiscriminator
				re.Avatar = userInfo.Avatar
				re.Auth = authToken

				session.Values["discord_id"] = re.ID
				session.Values["access_token"] = authToken.AccessToken
				session.Save(r, w)

				s.respondJSON(w, re, http.StatusOK)
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
		file := r.URL.Path
		fmt.Println("./frontend/build" + file)
		http.ServeFile(w, r, "./frontend/build"+file)
	}
}
