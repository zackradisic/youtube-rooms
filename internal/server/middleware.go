package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const discordIDKey contextKey = "discord_id"
const accessTokenKey contextKey = "access_token"

func (s *Server) checkAuthentication(next http.HandlerFunc) http.HandlerFunc {
	type requestBody struct {
		DiscordID string `json:"discordID"`
	}
	type sqlResult struct {
		DiscordID   string
		AccessToken string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var discordIDKey contextKey = "discord_id"
		var accessTokenKey contextKey = "access_token"

		decoder := json.NewDecoder(r.Body)
		rBody := &requestBody{}

		if err := decoder.Decode(rBody); err != nil {
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if rBody.DiscordID == "" {
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if r.Header.Get("Authorization") == "" {
			s.respondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.Contains(r.Header.Get("Authorization"), "Bearer ") {
			s.respondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
		result := &sqlResult{}
		query := "SELECT u.discord_id, ua.access_token FROM users u, user_auth ua WHERE ua.access_token=? AND u.discord_id=?"
		fmt.Println(accessToken, rBody.DiscordID)
		rows, err := s.DB.Raw(query, accessToken, rBody.DiscordID).Rows()
		if err != nil {
			s.respondError(w, "There was an error authorizing your request", 500)
			return
		}

		if !rows.Next() {
			fmt.Println("test")
			s.respondError(w, "There was an error authorizing your request", 500)
			return
		}

		rows.Scan(&result.DiscordID, &result.AccessToken)
		ctx = context.WithValue(ctx, discordIDKey, result.DiscordID)
		ctx = context.WithValue(ctx, accessTokenKey, result.AccessToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
