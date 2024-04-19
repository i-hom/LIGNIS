package main

import (
	"fmt"
	"lignis/internal/app"
	"lignis/internal/generated/api"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	var service *app.App

	service, err := app.NewApp()
	if err != nil {
		fmt.Println(err.Error())
	}

	var config struct {
		Runtime      string `envconfig:"RUNTIME"`
		HTTPBindPort int    `envconfig:"HTTP_BIND_PORT"`
	}

	envconfig.MustProcess("lignis", &config)

	srv, err := api.NewServer(service, service)
	if err != nil {
		log.Panicf("Error creating server: %s", err)
	}
	fmt.Printf("Starting server on port %d\n", config.HTTPBindPort)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Runtime, config.HTTPBindPort), app.Middleware{Next: srv}); err != nil {
		log.Fatal(err)
	}
}
