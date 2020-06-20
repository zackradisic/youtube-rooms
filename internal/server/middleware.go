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
		fmt.Println(session.Values)
		if !ok {
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
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		accessToken, ok := accessTokenRaw.(string)
		if !ok {
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user := &models.User{
			DiscordID: idString,
		}
		userAuth := &models.UserAuth{
			AccessToken: accessToken,
		}

		if err := s.DB.First(user).Error; err != nil {
			s.respondError(w, "Could not find that user", 404)
			return
		}
		userAuth.UserID = user.ID
		if err := s.DB.First(userAuth).Error; err != nil {
			s.respondError(w, "Invalid credentials", 401)
			return
		}

		ctx = context.WithValue(ctx, userKey, user)
		ctx = context.WithValue(ctx, userAuthKey, userAuth)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
