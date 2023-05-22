package main

import (
	"fmt"

	"github.com/ip-05/quizzus/api"
	"github.com/ip-05/quizzus/app/game"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/repo"
)

func main() {
	cfg := config.GetConfig()

	// Repository layer
	db := repo.New(cfg)
	gameRepo := repo.NewGameStore(cfg, db)

	// Business logic layer
	gameService := game.NewGameService(gameRepo)

	// TODO: Presentation layer
	// apiWeb := NewAuthController(gameService) ?
	// apiWs := s(gameService) ?

	r := api.InitWeb(cfg, db, gameService)

	r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}
