package main

import (
	"context"
	"flag"
	"log"

	tgClient "telegram_bot/clients/telegram"
	event_consumer "telegram_bot/consumer/event-consumer"
	"telegram_bot/events/telegram"
	"telegram_bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	// tgClient := telegram.New(tgBotHost, mustToken())

	// s:= files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}
	
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

	// fetcher = fetcher.New(tgClient)
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
