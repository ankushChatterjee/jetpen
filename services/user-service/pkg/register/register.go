package register

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ankushChatterjee/jetpen/user-service/datastore"
	"github.com/ankushChatterjee/jetpen/user-service/models"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/rabbitmq"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type UserData struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Name     string `form:"name"`
	Password string `form:"password"`
}

func (u *UserData) trimAll() {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)
	u.Name = strings.TrimSpace(u.Name)
}

func (u *UserData) IsAnyBlank() bool {
	u.trimAll()
	return len(u.Name) == 0 || len(u.Name) == 0 || len(u.Password) == 0 || len(u.Email) == 0
}

func RegisterUser(c *fiber.Ctx, db *sql.DB) error {
	userData := new(UserData)
	if err := c.BodyParser(userData); err != nil {
		return err
	}
	if userData.IsAnyBlank() {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "All values not entered",
		})
	}
	fmt.Println(userData)
	if !utils.ValidateEmail(userData.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Email is Invalid",
		})
	}

	tempUser := models.NewTempUser(userData.Username, userData.Email, userData.Name, userData.Password)
	if tempUser == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "We have some error, please try later",
		})
	}
	err := datastore.InsertTempUser(tempUser, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "We have some error, please try later",
		})
	}
	err = rabbitmq.PublishVerificationEmail(tempUser.Email, tempUser.Token, tempUser.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error Sending Email",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func EmailVerification(c *fiber.Ctx, db *sql.DB) error {
	token := c.Query("token")
	username := c.Query("username")
	fmt.Println(token, username)
	tempUser, err := datastore.GetTempUser(username, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "There is some error",
		})
	}
	fmt.Println(tempUser)
	if tempUser.Token == token {
		err = datastore.DeleteTempUser(username, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "There is some error",
			})
		}
		user := models.NewUserFromTempUser(tempUser)
		err = datastore.InsertUser(user, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "There is some error",
			})
		}
		return c.SendStatus(fiber.StatusOK)
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}
