package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func AuthorizeZippy(c *fiber.Ctx) error {
	code := c.Query("code")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("tbServersDB1").Collection("usersDB1")

	user := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"zippyCode": code}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"status": "INVALID"})
	}

	if user.ZippyExpires < time.Now().Unix() {
		return c.JSON(fiber.Map{"status": "TIMEOUT"})
	}

	payload := bson.M{
		"zippyCode":    nil,
		"zippyExpires": nil,
	}

	_, err = usersCollection.UpdateOne(context.Background(), bson.M{"zippyCode": code}, bson.M{"$set": payload})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating user after zippy authorization",
		})
	}

	return c.JSON(fiber.Map{"status": "OK", "token": user.Token})
}
