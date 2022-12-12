package main

import (
	"context"
	"log"
	//"encoding/json"
	"fmt"
	"os"

	"github.com/shomali11/slacker"
	"github.com/lpernett/godotenv"
	//"github.com/tidwall/gjson"
)

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.Command("query - <city>", &slacker.CommandDefinition{
		Description: "enter city name",
		Examples:    []string{"Delhi"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("city")

			fmt.Println(query)
			response.Reply(query)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
