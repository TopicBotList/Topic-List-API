package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func GetServer(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	serverID := c.Params("serverid")

	serversCollection := db.Database("TopicBots").Collection("serversDB1")

	var server types.Server
	err := serversCollection.FindOne(context.Background(), bson.M{"id": serverID}).Decode(&server)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "NOT_FOUND"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Find the owner details
	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	var user types.User
	err = usersCollection.FindOne(context.Background(), bson.M{"token": server.Owner}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Modify server data
	server.OwnerName = user.Name
	server.OwnerID = user.ID
	server.OwnerAvatar = user.Avatar

	// Remove _id field
	server.ID = serverID

	return c.JSON(fiber.Map{"status": "OK", "server": server})
}
