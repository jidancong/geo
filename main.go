package main

import (
	"log"

	"github.com/jidancong/geo/config"
	"github.com/jidancong/geo/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}

	app.Run(cfg)

	select {}

}
