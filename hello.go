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

type Resp struct {
	URL string `json:"url"`
	Zen string `json:"zen"`
}

type PushEvent struct {
	Ref string `json:"ref"`

	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`

	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`

	HeadCommit struct {
		Message string `json:"message"`
	} `json:"head_commit"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	type DiscordMessage struct {
		Content string `json:"content"`
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var content string

	switch r.Header.Get("X-GitHub-Event") {
	case "ping":
		var resp Resp
		err = json.Unmarshal(bodyBytes, &resp)
		if err != nil {
			http.Error(w, "Invalid ping JSON", http.StatusBadRequest)
			return
		}
		content = resp.Zen

	case "push":
		var event PushEvent
		err = json.Unmarshal(bodyBytes, &event)
		if err != nil {
			http.Error(w, "Invalid push JSON", http.StatusBadRequest)
			return
		}

		content = fmt.Sprintf(
			"New push to %s\nBranch: %s\nPushed by: %s\nCommit: %s",
			event.Repository.Name,
			event.Ref,
			event.Pusher.Name,
			event.HeadCommit.Message,
		)

	default:
		fmt.Println("Unknown event:", r.Header.Get("X-GitHub-Event"))
		return
	}

	message := DiscordMessage{
		Content: content,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := http.Post(
		os.Getenv("DISCORD_WEBHOOK_URL"),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Discord error:", err)
		return
	}
	defer response.Body.Close()

	fmt.Println("Discord status:", response.StatusCode)
}

func main() {
	godotenv.Load()

	http.HandleFunc("/event", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}