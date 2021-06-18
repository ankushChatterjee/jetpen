package sub

import (
	"database/sql"
	"github.com/ankushChatterjee/jetpen/sub-service/datastore"
	"github.com/ankushChatterjee/jetpen/sub-service/models"
	"github.com/ankushChatterjee/jetpen/sub-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AddRequest struct {
	Email string `json:"email"`
	Nid   string `json:"nid"`
}

func AddSubscription(c *fiber.Ctx, db *sql.DB) error {
	subscription := new(models.Subscription)
	req := new(AddRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}

	if !utils.ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Enter valid email",
		})
	}

	subscription.Email = req.Email
	subscription.Nid = req.Nid
	subscription.GenerateSubToken()
	subscription.GenerateID()
	err := datastore.InsertSubscription(subscription, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Couldn't get you subscribed",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func RemoveSubscription(c *fiber.Ctx, db *sql.DB) error {
	email := c.Query("email")
	nid := c.Query("nid")
	token := c.Query("token")
	sub, err := datastore.GetSubscriptionWithEmail(nid, email, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error" : "Error unsubscribing",
		})
	}
	if sub.SubToken != token {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error" : "Tokens do not match",
		})
	}
	err = datastore.DeleteSubscription(sub.Id, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error" : "Error unsubscribing",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
