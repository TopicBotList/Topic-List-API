package routes

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func EditBotSettings(c *fiber.Ctx) error {
	token := c.Query("token")
	botID := c.Query("id")
	publicityName := c.FormValue("publicity[name]")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botsCollection := db.Database("TopicBots").Collection("botsDB1")

	var result map[string]interface{}
	err := botsCollection.FindOne(context.TODO(), bson.M{"id": botID}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"reply": "BOT_NOT_FOUND"})
	}

	// Check token
	if result["token"] != token {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"result": "TOKEN_INVALID"})
	}

	payload := bson.M{"publicity": publicityName}
	_, err = botsCollection.UpdateOne(context.TODO(), bson.M{"id": botID},
		bson.M{"$set": payload})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"reply": "worked", "id": botID})
}
