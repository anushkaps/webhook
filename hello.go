package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	type DiscordMessage struct {
		Content string `json:"content"`
	}
	bodyBytes, err := io.ReadAll(r.Body)
	rstring := string(bodyBytes)
	message := DiscordMessage{
		Content: rstring,
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Post(os.Getenv("DISCORD_WEBHOOK_URL"), "application/json", bytes.NewBuffer(jsonData))
}

func main() {
	godotenv.Load()

	http.HandleFunc("/event", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}