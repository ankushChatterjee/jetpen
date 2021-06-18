package login

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ankushChatterjee/jetpen/user-service/datastore"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	jwt "github.com/form3tech-oss/jwt-go"
)

type LoginDetails struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Login(c *fiber.Ctx, db *sql.DB) error {
	details := new(LoginDetails)
	if err := c.BodyParser(details); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	fmt.Println(details)
	user, err := datastore.GetUser(details.Username, db)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Invalid Login",
		})
	}
	log.Println(details.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(details.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Password did not match",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24 * 14).Unix()

	t, err := token.SignedString([]byte(utils.GetEnvVar("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
