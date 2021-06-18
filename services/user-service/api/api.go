package api

import (
	"github.com/ankushChatterjee/jetpen/user-service/datastore"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/login"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/manage"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/rabbitmq"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/register"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
)

func CreateApi() *fiber.App {
	app := fiber.New()
	db := datastore.GetConnection()
	rabbitmq.Init()
	app.Use(cors.New())
	app.Post("/register", func(c *fiber.Ctx) error { return register.RegisterUser(c, db) })
	app.Get("/email-verification", func(c *fiber.Ctx) error { return register.EmailVerification(c, db) })
	app.Post("/login", func(c *fiber.Ctx) error { return login.Login(c, db) })
	app.Get("/username-exist", func(c *fiber.Ctx) error {return manage.DoesUsernameExist(c, db)})

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(utils.GetEnvVar("JWT_SECRET")),
	}))
	app.Get("/get-user", func(c *fiber.Ctx) error {return manage.GetUser(c, db)})
	app.Post("/edit-email", func(c *fiber.Ctx) error { return manage.EditEmail(c, db) })
	app.Post("/edit-password", func(c *fiber.Ctx) error { return manage.EditPassword(c, db) })
	app.Post("/edit-name", func(c *fiber.Ctx) error { return manage.EditName(c, db) })
	app.Post("/remove-user", func(c *fiber.Ctx) error { return manage.RemoveUser(c, db) })
	return app
}
