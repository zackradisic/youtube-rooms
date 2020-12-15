package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/matthewhartstonge/argon2"

	"github.com/zackradisic/youtube-rooms/internal/models"
	"github.com/zackradisic/youtube-rooms/internal/room"
	"github.com/zackradisic/youtube-rooms/internal/ws"
)

func (s *Server) setupRoutes() {
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})
	apiRouter := s.router.PathPrefix("/api/").Subrouter().StrictSlash(true)
	s.addRoute(apiRouter, "GET", "/auth/discord/", s.handleBeginAuth())
	s.addRoute(apiRouter, "GET", "/auth/discord/callback", s.rateLimited(s.handleCompleteAuth()))
	s.addRoute(apiRouter, "POST", "/rooms/verify/", s.rateLimited(s.handleVerifyPassword()))
	s.addRoute(apiRouter, "GET", "/rooms/", s.rateLimited(s.handleGetRooms()))
	s.addRoute(apiRouter, "GET", "/me", s.checkAuthentication(s.meRoute()))

	// s.addRoute(s.router, "POST", "/test", s.handleVerifyPassword())

	s.addRoute(s.router, "GET", "/ws", s.rateLimited(s.checkAuthentication(s.handleWS())))

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

func (s *Server) meRoute() http.HandlerFunc {
	type response struct {
		ID            string `json:"discordId"`
		Username      string `json:"discordUsername"`
		Discriminator string `json:"discordDiscriminator"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "https://theatre.theradisic.com")
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.respondError(w, err.Error(), 403)
			return
		}

		accessTokenRaw := session.Values["access_token"]
		if accessToken, ok := accessTokenRaw.(string); ok {

			if accessToken == "" {
				s.respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userInfo, err := s.getDiscordUserInfo(accessToken, "")
			if err != nil {
				fmt.Println(err.Error())
				s.respondError(w, "Internal server error occurred", http.StatusInternalServerError)
				return
			}

			res := &response{}
			res.ID = userInfo.ID
			res.Username = userInfo.Username
			res.Discriminator = userInfo.Discriminator

			s.respondJSON(w, res, http.StatusOK)
			return
		}
	}
}

func (s *Server) handleVerifyPassword() http.HandlerFunc {

	type jsonData struct {
		RoomName string `json:"roomName"`
		Password string `json:"password"`
	}

	type requestBody struct {
		Data jsonData `json:"data"`
	}

	type responseBody struct {
		Success bool `json:"success"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		decoder := json.NewDecoder(r.Body)
		body := &jsonData{}
		err := decoder.Decode(body)

		if err != nil {
			fmt.Println("failed")
			s.respondError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		rm, err := s.RoomManager.GetRoom(body.RoomName)
		if err != nil {
			fmt.Println(err.Error())
			s.respondError(w, "Couldn't find that room", http.StatusNotFound)
			return
		}

		ok, err := argon2.VerifyEncoded([]byte(body.Password), []byte(rm.Model.HashedPassword))
		if err != nil {
			s.respondError(w, "Internal server error", 500)
			return
		}

		s.respondJSON(w, &responseBody{Success: ok}, http.StatusOK)
	}
}

func (s *Server) handleGetRooms() http.HandlerFunc {
	type jsonRoom struct {
		ID                  uint   `json:"id"`
		Name                string `json:"name"`
		IsPasswordProtected bool   `json:"passwordProtected"`
		UsersCount          int    `json:"usersCount"`
	}

	type jsonResponse struct {
		Rooms []jsonRoom `json:"rooms"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		jr := &jsonResponse{
			Rooms: []jsonRoom{},
		}

		params := r.URL.Query()
		roomName, ok := params["name"]
		var rooms *[]models.Room
		var err error
		if !ok || len(roomName) == 0 {
			rooms, err = s.getRooms("")

		} else {
			rooms, err = s.getRooms(roomName[0])
		}

		if err != nil {
			s.respondError(w, "Error retrieving rooms", http.StatusInternalServerError)
			return
		}

		serverRooms := s.RoomManager.GetRooms()
		for _, m := range *rooms {
			userCount := 0
			for _, r := range serverRooms {
				if m.ID == r.Model.ID {
					userCount = len(r.GetUsers())
				}
			}

			jr.Rooms = append(jr.Rooms, jsonRoom{
				ID:                  m.ID,
				Name:                m.Name,
				IsPasswordProtected: m.HashedPassword != "",
				UsersCount:          userCount,
			})
		}

		s.respondJSON(w, jr, http.StatusOK)
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

		if roomName[0] == "" {
			fmt.Println("handleWS() -> Invalid room name")
			s.respondError(w, "Invalid room name", http.StatusBadRequest)
			return
		}

		rm, err := s.RoomManager.GetRoom(roomName[0])
		if err != nil {
			fmt.Println("handleWS() -> Couldn't find room")
			s.respondError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if rm.Model.HashedPassword != "" {
			roomPassword, ok := params["roomPassword"]
			if !ok {
				fmt.Println("handleWS() -> Bad request")
				s.respondError(w, "Bad Request", http.StatusBadRequest)
				return
			}
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
				if err := rm.AddUser(user); err != nil {
					fmt.Println("handleWS() -> Already in room")
					s.respondError(w, "Already in the room", http.StatusBadRequest)
					return
				}
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
					fmt.Println(err)
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

				http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusPermanentRedirect)
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
