package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/zackradisic/youtube-rooms/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	s := server.NewServer()
	s.Run(":3000")
}
