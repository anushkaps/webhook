package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/joho/godotenv"
)
type Resp struct {
	URL string `json:"url"`
	Event json.RawMessage `json:"event"`
}
func handler(w http.ResponseWriter, r *http.Request) {
	type DiscordMessage struct {
		Content string `json:"content"`
	}
	bodyBytes, err := io.ReadAll(r.Body)
    var resp Resp
	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}
	rurl := resp.URL
	revent := resp.Event
	rstring := string(revent)
	message := DiscordMessage{
		Content: rstring,
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Post(rurl, "application/json", bytes.NewBuffer(jsonData))
}

func main() {
	godotenv.Load()

	http.HandleFunc("/event", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}