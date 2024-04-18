package app

import (
	"lignis/internal/config"
	"lignis/internal/generated/api"
	"lignis/internal/repository"
	"lignis/internal/service/auth"
	"lignis/internal/storage"
)

type App struct {
	api.Handler

	config *config.Config

	auth *auth.Auth
	db   *storage.Database

	agentRepo      *repository.AgentRepo
	productRepo    *repository.ProductRepo
	saleRepo       *repository.SaleRepo
	userRepo       *repository.UserRepo
	acceptanceRepo *repository.AcceptanceRepo
	customerRepo   *repository.CustomerRepo
}

func NewApp() (*App, error) {
	var app App

	if err := app.initConfig(); err != nil {
		return &App{}, err
	}

	if err := app.initAuth(); err != nil {
		return &App{}, err
	}

	if err := app.initDB(); err != nil {
		return &App{}, err
	}

	app.initRepo()

	return &app, nil
}
