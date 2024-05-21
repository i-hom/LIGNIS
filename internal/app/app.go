package app

import (
	"fmt"
	"lignis/internal/config"
	"lignis/internal/generated/api"
	"lignis/internal/repository"
	"lignis/internal/service/auth"
	"lignis/internal/service/minio"
	"lignis/internal/storage"
)

type App struct {
	api.Handler

	config *config.Config

	auth  *auth.Auth
	db    *storage.Database
	minio *minio.MinioStorage

	agentRepo      *repository.AgentRepo
	productRepo    *repository.ProductRepo
	saleRepo       *repository.SaleRepo
	userRepo       *repository.UserRepo
	acceptanceRepo *repository.AcceptanceRepo
	customerRepo   *repository.CustomerRepo
	defectRepo     *repository.DefectRepo
	monthlyRepo    *repository.MonthlyRepo
}

func NewApp() (*App, error) {
	var app App

	if err := app.initConfig(); err != nil {
		return &App{}, err
	}
	fmt.Println("Successfully loaded .env file")
	if err := app.initAuth(); err != nil {
		return &App{}, err
	}

	if err := app.initDB(); err != nil {
		return &App{}, err
	}
	fmt.Println("Successfully connected to " + app.config.Mongo.MongoURI)
	if err := app.initMinio(); err != nil {
		return &App{}, err
	}
	fmt.Println("Successfully connected to " + app.config.Minio.Endpoint)
	app.initRepo()

	return &app, nil
}
