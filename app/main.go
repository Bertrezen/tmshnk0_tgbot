package main

import (
	"log"
	"tg_bot/bot"
	"tg_bot/cmd"
	rs "tg_bot/rest_api"
)

func main() {
	go bot.Work()

	if err := cmd.Root(rs.NewServer()).Execute(); err != nil {
		log.Fatal(err)
	}
}
