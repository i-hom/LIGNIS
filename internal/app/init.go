package app

import (
	"context"
	"lignis/internal/config"
	"lignis/internal/generated/api"
	"lignis/internal/repository"
	"lignis/internal/service/auth"
	"lignis/internal/service/minio"
	"lignis/internal/storage"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func (a *App) initConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	config := config.Config{}
	err = envconfig.Process("", &config)
	if err != nil {
		return err
	}
	a.config = &config
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

func (a *App) initMinio() error {
	var err error
	a.minio, err = minio.NewMinio(a.config.Minio)
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
	a.defectRepo = repository.NewDefectRepo(a.db.GetCollection("defects"))
	a.monthlyRepo = repository.NewMonthlyRepo(a.db.GetCollection("monthly"))
}

func (a App) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
	user, err := a.auth.ValidateAndParseToken(t.Token)
	if err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, "user", user)
	return ctx, nil
}

func (a App) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{StatusCode: 404, Response: api.Error{Message: err.Error()}}
}
