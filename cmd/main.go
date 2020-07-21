package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zackradisic/youtube-rooms/internal/server"
)

func main() {
	clientID := os.Getenv("DISCORD_CLIENT_ID")
	if clientID == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	s := server.NewServer()
	s.Run(":8000")
}
