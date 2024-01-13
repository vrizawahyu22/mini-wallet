// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-chi/chi"
	"github.com/google/wire"
	"mini-wallet/app"
	"mini-wallet/db/repository"
	"mini-wallet/handler"
	"mini-wallet/service"
	"mini-wallet/utils"
)

// Injectors from injector.go:

func InitializeApp(route *chi.Mux, DB utils.PGXPool, config *utils.BaseConfig) (app.App, error) {
	repository := querier.NewRepository(DB)
	walletSvc := service.NewWalletSvc(repository)
	walletHandler := handler.NewWalletHandler(walletSvc)
	appApp := app.NewApp(route, walletHandler, config)
	return appApp, nil
}

// injector.go:

var walletSet = wire.NewSet(querier.NewRepository, handler.NewWalletHandler, service.NewWalletSvc)