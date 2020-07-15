package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zackradisic/youtube-rooms/internal/models"
)

type contextKey string

const userKey contextKey = "user"
const userAuthKey contextKey = "userAuth"

func (s *Server) checkAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.respondError(w, err.Error(), 403)
			return
		}

		ctx := r.Context()

		id, ok := session.Values["discord_id"]
		if !ok {
			fmt.Println("handleWS() -> Couldn't find discord id in session cookie")
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		idString, ok := id.(string)
		if !ok {
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		accessTokenRaw, ok := session.Values["access_token"]
		if !ok {
			fmt.Println("handleWS() -> Couldn't find access_token in session cookie")
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		accessToken, ok := accessTokenRaw.(string)
		if !ok {
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		fmt.Printf("Trying to authenticate: (%s)\n", idString)
		user := &models.User{}

		userAuth := &models.UserAuth{}

		if err := s.DB.Where("discord_id = ?", idString).Find(user).Error; err != nil {
			s.respondError(w, "Could not find that user", 404)
			return
		}

		if err := s.DB.Where("access_token = ? AND user_id = ?", accessToken, user.ID).Find(userAuth).Error; err != nil {
			s.respondError(w, "Invalid credentials", 401)
			return
		}

		ctx = context.WithValue(ctx, userKey, user)
		ctx = context.WithValue(ctx, userAuthKey, userAuth)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
