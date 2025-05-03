package main

import (
	"fmt"
	"glbackend/internal/config"
	"glbackend/internal/infrastructure"
	"glbackend/internal/infrastructure/database"
	"glbackend/internal/infrastructure/log"
	"glbackend/internal/infrastructure/router"
	"time"
)

// @title						glbackend API
// @version						1.0
// @BasePath					/api
//
// @securityDefinitions.apikey	<BearerAuth>
// @in							header
// @name						Authorization
// @host						localhost:8080
func main() {
	config, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("config: %+v\n", config)
	a := infrastructure.
		NewConfig(config).
		ContextTimeout(15 * time.Second).
		Logger(log.InstanceLogrusLogger).
		DBGSql(database.InstanceGPostgres).
		GTransaction(database.InstanceGPostgres)

	a.WebServer(router.InstanceGin)
	a.Start()
}
