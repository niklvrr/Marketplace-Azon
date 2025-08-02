package app

import (
	"fmt"
	"log"
	"myMarketplace/internal/config"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	// TODO logger

	// TODO router

	// TODO db

}
