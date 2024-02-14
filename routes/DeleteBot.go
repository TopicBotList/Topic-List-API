package routes

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteBot(c *fiber.Ctx) error {
	token := c.Query("token")
	botID := c.Params("botid")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botsCollection := db.Database("TopicBots").Collection("botsDB1")

	var bot map[string]interface{}
	err := botsCollection.FindOne(context.TODO(), bson.M{"id": botID}).Decode(&bot)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"result": "BOT_NOT_FOUND"})
	}

	// Check token
	if bot["token"] != token {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"result": "TOKEN_INVALID"})
	}

	_, err = botsCollection.DeleteOne(context.TODO(), bson.M{"id": botID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"result": "VALID"})
}
