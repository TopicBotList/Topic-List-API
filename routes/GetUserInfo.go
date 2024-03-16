package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func GetUser(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	userCollection := db.Database("TopicBots").Collection("usersDB1")
	serverCollection := db.Database("TopicBots").Collection("serversDB1")

	userID := c.Params("userid")

	var user types.User
	err := userCollection.FindOne(context.TODO(), bson.M{"id": userID}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"result": "invalid"})
	}

	// Redacting sensitive user fields
	user.AccessToken = "Redacted"

	var servers []types.Server
	cursor, err := serverCollection.Find(context.TODO(), bson.M{"owner": user.Token})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.Background()) {
		var server types.Server
		if err := cursor.Decode(&server); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		// Assign user data to server object
		server.OwnerAvatar = user.Avatar
		server.OwnerName = user.Name
		server.OwnerID = user.ID
		server.Owner = ""

		servers = append(servers, server)
	}

	// Assign servers to user
	user.Servers = servers

	// Omitting token field from the response
	user.Token = ""

	return c.JSON(fiber.Map{"result": user})
}
