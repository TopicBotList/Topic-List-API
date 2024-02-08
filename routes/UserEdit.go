package routes

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserEditBots(c *fiber.Ctx) error {
	token := c.Query("token")
	bio := c.FormValue("bio")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	var user map[string]interface{}
	err := usersCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"reply": "TOKEN_INVALID"})
	}

	payload := bson.M{"bio": bio}
	_, err = usersCollection.UpdateOne(context.TODO(), bson.M{"token": token}, bson.M{"$set": payload})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"reply": "worked"})
}
