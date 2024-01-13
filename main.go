package main

import (
	"mini-wallet/utils"

	"github.com/go-chi/chi"
)

func main() {
	utils.InitLogger()

	config := utils.CheckAndSetConfig("./config", "app")

	DBpool := utils.ConnectDBPool(config.DBConnString)
	DB := utils.ConnectDB(config.DBConnString)
	utils.RunMigrationPool(DB, config, true)

	r := chi.NewRouter()

	app, err := InitializeApp(r, DBpool, config)
	utils.LogAndPanicIfError(err, "failed when starting app")

	app.Start()
}
