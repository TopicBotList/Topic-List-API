// private/add_server.go

package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddServer(c *fiber.Ctx) error {
	captchaResponse := c.FormValue("captcha")

	hCaptchaResp, err := http.PostForm("https://hcaptcha.com/siteverify", map[string][]string{
		"secret":   {os.Getenv("HCAPTCHA_SECRET")},
		"response": {captchaResponse},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "ERROR",
			"error":  "Failed to verify hCaptcha",
		})
	}
	defer hCaptchaResp.Body.Close()

	var hCaptchaResult map[string]interface{}
	if err := json.NewDecoder(hCaptchaResp.Body).Decode(&hCaptchaResult); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "ERROR",
			"error":  "Failed to decode hCaptcha response",
		})
	}

	if success, ok := hCaptchaResult["success"].(bool); !ok || !success {
		errorCodes, _ := hCaptchaResult["error-codes"].([]interface{})
		return c.JSON(fiber.Map{
			"status": "HERROR",
			"error":  errorCodes,
		})
	}

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "ERROR",
			"error":  "Database connection not available",
		})
	}

	serversCollection := db.Database("TopicBots").Collection("serversDB1")

	serverPayload := bson.M{
		"name":        c.FormValue("name"),
		"icon":        c.FormValue("icon"),
		"id":          c.FormValue("id"),
		"owner":       c.Query("token"),
		"votes":       0,
		"category":    c.FormValue("category"),
		"views":       0,
		"summary":     c.FormValue("summary"),
		"description": c.FormValue("description"),
		"invite":      c.FormValue("invite"),
	}

	_, err = serversCollection.InsertOne(context.Background(), serverPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "ERROR",
			"error":  "Failed to insert server into the database",
		})
	}

	return c.JSON(fiber.Map{
		"status": "OK",
		"server": c.FormValue("id"),
	})
}
