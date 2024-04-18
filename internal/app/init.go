package app

import (
	dev "lignis/build"
	"lignis/internal/repository"
	"lignis/internal/service/auth"
	"lignis/internal/storage"

	"github.com/kelseyhightower/envconfig"
)

func (a *App) initConfig() error {
	if dev.Dev {
		err := dev.LoadConfig()
		if err != nil {
			return err
		}
	}

	err := envconfig.Process("", a.config)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initDB() error {
	var err error
	a.db, err = storage.NewMongo(a.config.Mongo)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initAuth() error {
	a.auth = auth.NewAuth(a.config.Auth)
	return nil
}

func (a *App) initRepo() {
	a.agentRepo = repository.NewAgentRepo(a.db.GetCollection("agents"))
	a.productRepo = repository.NewProductRepo(a.db.GetCollection("products"))
	a.saleRepo = repository.NewSaleRepo(a.db.GetCollection("sales"))
	a.userRepo = repository.NewUserRepo(a.db.GetCollection("users"))
	a.acceptanceRepo = repository.NewAcceptanceRepo(a.db.GetCollection("acceptances"))
	a.customerRepo = repository.NewCustomerRepo(a.db.GetCollection("customers"))
}
