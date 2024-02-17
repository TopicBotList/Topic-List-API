// routes/get_user_info.go

package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func GetUserInfo(c *fiber.Ctx) error {
	token := c.Query("token")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	user := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"result": "UNKNOWN_USER"})
	}

	user.ID = ""       // Remove MongoDB's "_id" field
	user.Token = ""    // Remove the token for security reasons
	user.Password = "" // Remove the password for security reasons

	return c.JSON(fiber.Map{"result": user})
}
