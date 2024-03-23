package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func EditUser(c *fiber.Ctx) error {
	token := c.Query("token")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	serversCollection := db.Database("TopicBots").Collection("serversDB1")

	data := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&data)
	if err != nil {
		return c.JSON(fiber.Map{"status": "NOT_FOUND"})
	}

	payload := bson.M{
		"bio":     c.FormValue("bio"),
		"longBio": c.FormValue("longBio"),
	}

	_, err = usersCollection.UpdateOne(context.Background(), bson.M{"token": token}, bson.M{"$set": payload})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating user",
		})
	}

	data = types.User{}
	err = usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching updated user data",
		})
	}

	userServersCursor, err := serversCollection.Find(context.Background(), bson.M{"owner": data.Token})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding user servers",
		})
	}
	defer userServersCursor.Close(context.Background())

	var fServer []types.Server
	for userServersCursor.Next(context.Background()) {
		var server types.Server
		err := userServersCursor.Decode(&server)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error decoding server data",
			})
		}

		user := types.User{}

		server.OwnerName = user.Name
		server.OwnerID = user.ID
		server.OwnerAvatar = user.Avatar
		server.Owner = "Redacted" // Remove the owner field
		fServer = append(fServer, server)
	}

	data.Owner = ""
	data.Token = "" // Remove the token for security reasons
	data.Servers = fServer

	return c.JSON(fiber.Map{"status": "OK", "result": data})
}
