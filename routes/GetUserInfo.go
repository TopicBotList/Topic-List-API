// routes/get_user_info.go

package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func GetUser(c *fiber.Ctx) error {
	userID := c.Params("userid")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	serversCollection := db.Database("TopicBots").Collection("serversDB1")

	data := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"id": userID}).Decode(&data)
	if err != nil {
		return c.JSON(fiber.Map{"result": "invalid"})
	}

	data.AccessToken = ""
	data.ID = data.ID // Remove MongoDB's "_id" field

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
		err = usersCollection.FindOne(context.Background(), bson.M{"token": server.Owner}).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error finding owner user",
			})
		}

		server.OwnerName = user.Name
		server.OwnerID = user.ID
		server.OwnerAvatar = user.Avatar
		server.ID = server.ID // Remove MongoDB's "_id" field
		server.Owner = ""
		fServer = append(fServer, server)
	}

	data.Owner = ""
	data.Token = "" // Remove the token for security reasons
	data.Servers = fServer

	return c.JSON(fiber.Map{"result": data})
}
