// routes/edit_server.go

package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func EditServer(c *fiber.Ctx) error {
	serverID := c.Params("serverid")
	token := c.Query("token")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	serversCollection := db.Database("TopicBots").Collection("serversDB1")

	data := types.Server{}
	err := serversCollection.FindOne(context.Background(), bson.M{"id": serverID}).Decode(&data)
	if err != nil {
		return c.JSON(fiber.Map{"status": "NOT_FOUND"})
	}

	user := types.User{}
	err = usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"status": "LOGIN"})
	}

	data.OwnerName = user.Name
	data.OwnerID = user.ID
	data.OwnerAvatar = user.Avatar
	data.ID = "" // Remove MongoDB's "_id" field

	if user.Token != data.Owner {
		return c.JSON(fiber.Map{"status": "NOT_ALLOWED", "data": data})
	}

	payload := bson.M{
		"summary":     c.FormValue("summary"),
		"description": c.FormValue("description"),
		"category":    c.FormValue("category"),
		"invite":      c.FormValue("invite"),
	}

	_, err = serversCollection.UpdateOne(context.Background(), bson.M{"id": serverID}, bson.M{"$set": payload})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating server",
		})
	}

	data = types.Server{}
	err = serversCollection.FindOne(context.Background(), bson.M{"id": serverID}).Decode(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching updated server data",
		})
	}

	data.OwnerName = user.Name
	data.OwnerID = user.ID
	data.OwnerAvatar = user.Avatar
	data.ID = "" // Remove MongoDB's "_id" field

	return c.JSON(fiber.Map{"status": "OK", "data": data})
}
