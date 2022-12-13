package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/shomali11/slacker"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func weatherQuery(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + os.Getenv("OPENWEATHERMAP_API_KEY") + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func convertKtoC(tempK float64) float64 {
	tempC := tempK - 273.15
	return tempC
}

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("query - <city>", &slacker.CommandDefinition{
		Description: "enter city name",
		Examples:    []string{"Delhi"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("city")

			data, err := weatherQuery(query)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(data.Main.Kelvin)
			tempC := convertKtoC(data.Main.Kelvin)
			fmt.Println(tempC)

			tempstr := fmt.Sprintf("%.3f", tempC) + "C degrees"
			response.Reply(tempstr)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
