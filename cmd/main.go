package main

import (
	"log"
	"os"
	"tg_lists-and-groups/internal/tgbot"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loadind .env file: %v\n", err)
	}
	token := os.Getenv("TOKEN")
	/*	dbURL := os.Getenv("DATABASE_URL")
		if token == "" || dbURL == "" {
			log.Fatalf("TOKEN or DATABASE_URL not set\n")
		}
		if err := app.InitDB(dbURL); err != nil {
			log.Fatalf("Init DB error: %v\n", err)
		}
		defer app.CloseDB()*/

	// убрать после подключения бд
	if token == "" {
		log.Fatalf("TOKEN or DATABASE_URL not set\n")
	}

	// bot
	if err := tgbot.RunBot(token); err != nil {
		log.Fatalf("Bot exited with error: %v", err)
	}
}
