package manage

import (
	"database/sql"

	"github.com/ankushChatterjee/jetpen/user-service/datastore"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"

	jwt "github.com/form3tech-oss/jwt-go"
)

type JwtRequest struct {
	Content  string `form:"content"`
	Username string `form:"username"`
}

const EMAIL = 0
const PASSWORD = 1
const NAME = 2

func EditName(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(JwtRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	if username != req.Username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "UnAuthorized",
		})
	}
	err := datastore.UpdateNameOfUser(username, req.Content, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error Updating Name",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func EditPassword(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(JwtRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	if username != req.Username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "UnAuthorized",
		})
	}
	err := datastore.UpdateUserPassword(username, req.Content, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error Updating Name",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func EditEmail(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(JwtRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	if username != req.Username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "UnAuthorized",
		})
	}
	if utils.ValidateEmail(req.Content) {
		err := datastore.UpdateUserEmail(username, req.Content, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Error Updating Name",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Invalid Email",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func RemoveUser(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(JwtRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	if username != req.Username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "UnAuthorized",
		})
	}
	err := datastore.DeleteUser(username, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error Updating Name",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
