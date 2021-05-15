package main

import (
	"github.com/ankushChatterjee/jetpen/user-service/api"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
)

func main() {
	app := api.CreateApi()
	port := utils.GetEnvVar("HOST_PORT")
	if len(port) == 0 {
		port = "3001"
	}
	app.Listen(":" + port)
}
