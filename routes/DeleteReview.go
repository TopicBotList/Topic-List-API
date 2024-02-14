package routes

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteReview(c *fiber.Ctx) error {
	token := c.Query("token")
	botID := c.Params("botid")
	reviewID := c.Query("reviewid")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botsCollection := db.Database("TopicBots").Collection("botsDB1")
	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	var user map[string]interface{}
	err := usersCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"reply": "BOT_INVALID"})
	}

	var data map[string]interface{}
	err = botsCollection.FindOne(context.TODO(), bson.M{"id": botID}).Decode(&data)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"reply": "BOT_INVALID"})
	}

	reviews := data["reviews"].([]interface{})
	freviews := make([]interface{}, 0)

	for _, r := range reviews {
		review := r.(map[string]interface{})
		if review["id"].(string) == reviewID {
			if review["token"].(string) == token {
				continue
			} else {
				return c.Status(http.StatusForbidden).JSON(fiber.Map{"reply": "FORBIDDEN"})
			}
		}
		freviews = append(freviews, review)
	}

	_, err = botsCollection.UpdateOne(context.TODO(), bson.M{"id": botID},
		bson.M{"$set": bson.M{"reviews": freviews}})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"reply": "OK"})
}
