package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

type Listener[T any] func(T)

type ForecastResponse struct {
	City struct {
		Name string `json:"name"`
	} `json:"city"`

	List []struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`

		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"list"`
}

func main() {
	godotenv.Load()

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	discordURL := os.Getenv("DISCORD_WEBHOOK_URL")

	serverurl := "http://api.openweathermap.org/data/2.5/forecast?id=524901&appid=" + apiKey
	response, err := http.Get(serverurl)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	var forecast ForecastResponse
	err = json.Unmarshal(data, &forecast)
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}
	first := forecast.List[0]

	fmt.Println("City:", forecast.City.Name)
	fmt.Println("Temperature:", first.Main.Temp, "°C")
	fmt.Println("Humidity:", first.Main.Humidity, "%")

	type DiscordMessage struct {
		Content string `json:"content"`
	}
	message := DiscordMessage{
		Content: fmt.Sprintf(
			"Weather in %s\nTemperature: %.1f\nHumidity: %d%%",
			forecast.City.Name,
			first.Main.Temp,
			first.Main.Humidity,
		),
	}
	fmt.Println(message)
	jsonData, err := json.Marshal(message)
	fmt.Println(jsonData)
	fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}

	http.Post(discordURL,"application/json", bytes.NewBuffer(jsonData))
}
