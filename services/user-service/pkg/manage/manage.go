package manage

import (
	"database/sql"
	"strings"

	"github.com/ankushChatterjee/jetpen/user-service/datastore"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"

	jwt "github.com/form3tech-oss/jwt-go"
)

type EditRequest struct {
	Content  string `json:"content"`
}

func EditName(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(EditRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	req.Content = strings.Trim(req.Content, " ")
	if len(req.Content) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Empty Content",
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
	req := new(EditRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	req.Content = strings.Trim(req.Content, " ")
	if len(req.Content) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Empty Content",
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
	req := new(EditRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
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
	req := new(EditRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
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

func GetUser(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	userObj, err := datastore.GetUser(username, db)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Invalid Login",
		})
	}
	return c.JSON(fiber.Map{
		"name": userObj.Name,
		"email":userObj.Email,
	})
}

func DoesUsernameExist(c *fiber.Ctx, db *sql.DB) error {
	username := c.Query("username")
	isTaken, err := datastore.IsUsernameTaken(username, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error checking username availablity",
		})
	}
	return c.JSON(fiber.Map{
		"isTaken": isTaken,
	})
}