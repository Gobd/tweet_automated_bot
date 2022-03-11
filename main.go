package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
	"github.com/SoyPete/tweet_automated_bot/internal/botguts"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("we are live"))
}

func main() {
	ctx := context.Background()
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: remove and setup permanent datastore
	defer db.Close(ctx)

	// TODO: ここでbotを作成する
	bot := botguts.NewAutoBot(db, client)
	err = bot.TweetYoutubeVideo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client.RunDiscordBot()

	http.HandleFunc("/health", healthCheck)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
