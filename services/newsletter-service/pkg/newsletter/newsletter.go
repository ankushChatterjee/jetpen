package newsletter

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"

	"github.com/ankushChatterjee/jetpen/newsletter-service/datastore"
	"github.com/ankushChatterjee/jetpen/newsletter-service/models"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type CreateRequest struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}

type DeleteRequest struct {
	Nid string `json:"nid"`
}

type EditRequest struct {
	Nid     string `json:"nid"`
	Content string `json:"content"`
}

func CreateNewsLetter(c *fiber.Ctx, db *sql.DB) error {
	log.Println("CreateNEws")
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(CreateRequest)
	if err := c.BodyParser(req); err != nil {
		log.Println("Body Parser error")
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	if len(req.Name) == 0 {
		log.Println("Empty name error")
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Empty Name",
		})
	}

	newsletter := models.Newsletter{}
	newsletter.Name = req.Name
	newsletter.Description = req.Description
	newsletter.Owner = username
	newsletter.Id = uuid.NewV4().String()
	err := datastore.InsertNewsletter(&newsletter, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error Adding Newsletter",
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"id":newsletter.Id,
	})
}

func DeleteNewsletter(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(DeleteRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	owner, err := datastore.GetNewsLetterOwner(req.Nid, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error deleting newsletter",
		})
	}
	if owner != username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Error deleting newsletter",
		})
	}
	err = datastore.DeleteNewsletter(req.Nid, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Deleting newsletter",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func EditNewsletterName(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(EditRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	owner, err := datastore.GetNewsLetterOwner(req.Nid, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error deleting newsletter",
		})
	}
	if owner != username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Error deleting newsletter",
		})
	}
	err = datastore.UpdateNewsletterName(req.Nid, req.Content, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error Updating Name",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func EditNewsletterDescription(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(EditRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	owner, err := datastore.GetNewsLetterOwner(req.Nid, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error deleting newsletter",
		})
	}
	if owner != username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Error deleting newsletter",
		})
	}
	err = datastore.UpdateNewsletterDescription(req.Nid, req.Content, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error Updating descripition",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetNewsLetters(c *fiber.Ctx, db *sql.DB) error {
	cursor := c.Query("cursor")
	limit := c.Query("limit")
	l, err := strconv.Atoi(limit)
	if err != nil || l < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Bad Limit",
		})
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	newsletters, nextCursor, err := datastore.GetNewsLettersForUser(username, cursor, l, db)
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"newsletters": newsletters,
		"nextCursor":  nextCursor,
	})
}

func GetSingleNewsLetter(c *fiber.Ctx, db *sql.DB) error {
	nid := c.Params("id")
	if len(nid) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "ID is not provided properly",
		})
	}
	newsletter, err := datastore.GetNewsLetter(nid, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Unable to get newsletter",
		})
	}
	return c.JSON(newsletter)
}
