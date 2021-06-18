package api

import (
	"github.com/ankushChatterjee/jetpen/sub-service/datastore"
	"github.com/ankushChatterjee/jetpen/sub-service/pkg/sub"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CreateAPI() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	db := datastore.GetConnection()
	app.Post("/", func(ctx *fiber.Ctx) error { return  sub.AddSubscription(ctx, db)})
	app.Get("/unsub", func(ctx *fiber.Ctx) error { return  sub.RemoveSubscription(ctx, db)})
	return app
}
