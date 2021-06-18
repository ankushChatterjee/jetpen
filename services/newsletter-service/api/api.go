package api

import (
	"github.com/ankushChatterjee/jetpen/newsletter-service/datastore"
	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/letter"
	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/newsletter"
	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/rabbitmq"
	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
)

func CreateApi() *fiber.App {
	app := fiber.New()
	db := datastore.GetConnection()
	rabbitmq.Init()
	app.Use(cors.New())

	app.Get("/newsletter/:id", func(c *fiber.Ctx) error { return newsletter.GetSingleNewsLetter(c, db) })
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(utils.GetEnvVar("JWT_SECRET")),
	}))
	app.Get("/newsletters", func(c *fiber.Ctx) error { return newsletter.GetNewsLetters(c, db) })
	app.Post("/newsletter/create", func(c *fiber.Ctx) error { return newsletter.CreateNewsLetter(c, db) })
	app.Post("/newsletter/delete", func(c *fiber.Ctx) error { return newsletter.DeleteNewsletter(c, db) })
	app.Post("/newsletter/edit-name", func(c *fiber.Ctx) error { return newsletter.EditNewsletterName(c, db) })
	app.Post("/newsletter/edit-description", func(c *fiber.Ctx) error { return newsletter.EditNewsletterDescription(c, db) })

	app.Get("/letters/:nid", func(c *fiber.Ctx) error { return letter.GetLetters(c, db) })
	app.Get("/letter/:id", func(c *fiber.Ctx) error { return letter.GetLetterContent(c, db) })
	app.Post("/letter/create", func(c *fiber.Ctx) error { return letter.CreateLetter(c, db) })
	app.Post("/letter/publish", func(c *fiber.Ctx) error { return letter.SendDraftLetter(c, db) })
	app.Post("/letter/save", func(c *fiber.Ctx) error { return letter.SaveDraft(c, db) })
	app.Post("/letter/delete", func(c *fiber.Ctx) error { return letter.DeleteDraft(c, db) })

	return app
}
