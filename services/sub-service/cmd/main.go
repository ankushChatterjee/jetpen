package main

import (
	"github.com/ankushChatterjee/jetpen/sub-service/api"
	"github.com/ankushChatterjee/jetpen/sub-service/pkg/utils"
)

func main() {
	app := api.CreateAPI()
	port := utils.GetEnvVar("HOST_PORT")
	if len(port) == 0 {
		port = "3003"
	}
	err := app.Listen(":" + port)
	if err != nil {
		panic(err)
	}
}
