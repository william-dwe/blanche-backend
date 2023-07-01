package main

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cache"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cronjob"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/db"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/server"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Blanche API Backend Started")
	log.Info().Msg("CONFIG: " + config.Config.ENVConfig.Mode)

	dbErr := db.Connect()
	if dbErr != nil {
		log.Fatal().Msg("error connecting to DB")
	}
	rdbErr := cache.ConnectRDB()
	if rdbErr != nil {
		log.Fatal().Msg("error connecting to RDB")
	}

	cronjob.Init()
	server.Init()

	log.Info().Msg("Blanche API Backend is ready!")
}
