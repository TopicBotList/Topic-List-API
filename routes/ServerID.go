// routes/server.go

package routes

import (
	"context"
	"strconv" // Add this import statement

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func GetServer(c *fiber.Ctx) error {
	serverID := c.Params("serverid")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	serversCollection := db.Database("tbServersDB1").Collection("serversDB1")
	usersCollection := db.Database("tbServersDB1").Collection("usersDB1")

	data := types.Server{}
	err := serversCollection.FindOne(context.Background(), bson.M{"id": serverID}).Decode(&data)
	if err != nil {
		return c.JSON(fiber.Map{"status": "NOT_FOUND"})
	}

	user := types.User{}
	err = usersCollection.FindOne(context.Background(), bson.M{"token": data.Owner}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding owner user",
		})
	}

	data.OwnerName = user.Name
	data.OwnerID = user.ID
	data.OwnerAvatar = user.Avatar
	data.Owner = strconv.FormatBool(user.Token == c.Query("token")) // Convert boolean to string

	data.ID = "" // Remove MongoDB's "_id" field
	return c.JSON(fiber.Map{"status": "OK", "server": data})
}
