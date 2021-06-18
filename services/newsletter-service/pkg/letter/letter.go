package letter

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"

	"github.com/ankushChatterjee/jetpen/newsletter-service/datastore"
	"github.com/ankushChatterjee/jetpen/newsletter-service/models"
	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/rabbitmq"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type CreateRequest struct {
	Nid         string `json:"nid"`
	Content     string `json:"content"`
	Subject     string `json:"subject"`
	IsPublished bool   `json:"isPublished"`
}
type DraftRequest struct {
	Id      string `json:"id"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Nid 	string `json:"nid"`
}

type DeleteRequest struct {
	Id string `json:"id"`
}

func CreateLetter(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(CreateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	log.Println(req.Nid)
	owner, err := datastore.GetNewsLetterOwner(req.Nid, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error creating letter",
		})
	}
	if owner != username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Error creating letter",
		})
	}

	letter := new(models.Letter)
	letter.Content = req.Content
	letter.Nid = req.Nid
	letter.Id = uuid.NewV4().String()
	letter.Owner = username
	letter.IsPublished = req.IsPublished
	letter.Subject = req.Subject
	letter.PublishedAt.Time = time.Now().UTC()
	err = datastore.InsertLetter(letter, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Cannot save letter",
		})
	}
	if letter.IsPublished {
		subs, err := datastore.GetSubsAndTokens(req.Nid, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Cannot get sub",
			})
		}
		err = rabbitmq.PublishLetter(subs, req.Subject, req.Content, letter.Nid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Cannot Send letter",
			})
		}
	}
	return c.JSON(&fiber.Map{
		"id":letter.Id,
	})
}

func SendDraftLetter(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(DraftRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	owner, err := datastore.GetLetterOwner(req.Id, db)
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
	isPublished, err := datastore.IsLetterPublished(req.Id, db)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Error sending letter",
		})
	}
	if !isPublished {
		err := datastore.UpdateLetterContentSubjectAndPublish(req.Id, req.Subject, req.Content, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Cannot update letter",
			})
		}
		subs, err := datastore.GetSubsAndTokens(req.Nid, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Cannot get sub",
			})
		}
		err = rabbitmq.PublishLetter(subs, req.Subject, req.Content, req.Nid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Cannot Send letter",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Letter already published",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func SaveDraft(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(DraftRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	owner, err := datastore.GetLetterOwner(req.Id, db)
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
	isPublished, err := datastore.IsLetterPublished(req.Id, db)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Error sending letter",
		})
	}
	if !isPublished {
		err := datastore.UpdateLetterContentAndSubject(req.Id, req.Subject, req.Content, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Cannot udpate letter",
			})
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

func DeleteDraft(c *fiber.Ctx, db *sql.DB) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	req := new(DeleteRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Form Data Error",
		})
	}
	owner, err := datastore.GetLetterOwner(req.Id, db)
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
	isPublished, err := datastore.IsLetterPublished(req.Id, db)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Error sending letter",
		})
	}
	if !isPublished {
		err := datastore.DeleteLetter(req.Id, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Cannot udpate letter",
			})
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetLetters(c *fiber.Ctx, db *sql.DB) error {
	cursor := c.Query("cursor")
	limit := c.Query("limit")
	nid := c.Params("nid")

	if len(nid) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Bad nesletter id",
		})
	}

	l, err := strconv.Atoi(limit)
	if err != nil || l < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Bad Limit",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	owner, err := datastore.GetNewsLetterOwner(nid, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error getting letters",
		})
	}
	if owner != username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Auth Error",
		})
	}

	letters, nextCursor, err := datastore.GetLettersForNewsletter(nid, cursor, l, db)
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"letters":    letters,
		"nextCursor": nextCursor,
	})
}

func GetLetterContent(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Incorrect ID",
		})
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	owner, err := datastore.GetLetterOwner(id, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error getting letter content",
		})
	}
	if owner != username {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "Error getting letter content",
		})
	}
	content, err := datastore.GetLetterContent(id, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Error getting letter content",
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"content": content,
	})
}
