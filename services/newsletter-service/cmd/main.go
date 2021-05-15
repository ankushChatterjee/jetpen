package main

import (
	"github.com/ankushChatterjee/jetpen/newsletter-service/api"
	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/utils"
)

func main() {
	app := api.CreateApi()
	port := utils.GetEnvVar("HOST_PORT")
	if len(port) == 0 {
		port = "3002"
	}
	err := app.Listen(":" + port)
	if err != nil {
		panic(err)
	}
}
