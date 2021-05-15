package api

import (
	"github.com/ankushChatterjee/jetpen/sub-service/datastore"
	"github.com/ankushChatterjee/jetpen/sub-service/pkg/sub"
	"github.com/gofiber/fiber/v2"
)

func CreateAPI() *fiber.App {
	app := fiber.New()
	db := datastore.GetConnection()
	app.Post("/sub", func(ctx *fiber.Ctx) error { return  sub.AddSubscription(ctx, db)})
	app.Get("/unsub", func(ctx *fiber.Ctx) error { return  sub.RemoveSubscription(ctx, db)})
	return app
}
