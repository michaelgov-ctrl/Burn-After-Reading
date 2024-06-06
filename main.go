package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
	MessageBucket
}

/*
	example post:
		Invoke-WebRequest -Uri http://localhost:80/v1/messages -Method Post -Body '{"content": "hello"}'

	example get:
		Invoke-WebRequest -Uri http://localhost:80/v1/messages?uuid=c156c173-ab45-45c4-a168-5530f7b21887
*/

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|test|prod)")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
		MessageBucket: MessageBucket{
			Bucket: make(map[uuid.UUID]string),
		},
	}

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
