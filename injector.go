//go:build wireinject
// +build wireinject

package main

import (
	"mini-wallet/app"
	"mini-wallet/handler"
	"mini-wallet/service"
	"mini-wallet/utils"

	querier "mini-wallet/db/repository"

	"github.com/go-chi/chi"
	"github.com/google/wire"
)

var walletSet = wire.NewSet(
	querier.NewRepository,
	handler.NewWalletHandler,
	service.NewWalletSvc,
)

func InitializeApp(
	route *chi.Mux,
	DB utils.PGXPool,
	config *utils.BaseConfig,
) (app.App, error) {
	wire.Build(
		walletSet,
		app.NewApp,
	)

	return nil, nil
}
