package usi_test

import (
	"context"
	"log"
	"os"

	"github.com/kk-no/go-usi"
)

func Example() {
	ctx := context.Background()

	enginePath := os.Getenv("ENGINE_PATH")
	engine, err := usi.New(enginePath)
	if err != nil {
		log.Fatalln(err)
	}

	if err := engine.Connect(ctx); err != nil {
		os.Exit(0)
	}
}
