package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.topiclist.xyz/types"
)

func FindServersByCategory(c *fiber.Ctx) error {
	category := c.Params("cat")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	serversCollection := db.Database("TopicBots").Collection("serversDB1")

	// Find servers with the specified category, sorted by votes in descending order
	cursor, err := serversCollection.Find(context.Background(), bson.M{"category": category}, options.Find().SetSort(bson.M{"votes": -1}))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding servers",
		})
	}
	defer cursor.Close(context.Background())

	var fdata []types.Server
	for cursor.Next(context.Background()) {
		var server types.Server
		err := cursor.Decode(&server)
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
		user.AccessToken = "" //remove access token field
		user.Token = ""
		server.OwnerAvatar = user.Avatar
		server.Owner = ""
		server.ID = server.ID // Remove MongoDB's "_id" field

		fdata = append(fdata, server)
	}

	// Fetch all users
	fusers := []types.User{}
	usersCursor, err := usersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding users",
		})
	}
	defer usersCursor.Close(context.Background())

	for usersCursor.Next(context.Background()) {
		var fuser types.User
		err := usersCursor.Decode(&fuser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error decoding user data",
			})
		}

		// Mask the token field
		fuser.Token = "********" // Masked value
		fuser.AccessToken = "********"
		fuser.ID = fuser.ID
		fuser.Password = "" // Remove the password for security reasons

		// Remove the owner field from the user
		fuser.Owner = "********" //remove this thindy
		fuser.Token = "********"
		fusers = append(fusers, fuser)
	}

	return c.JSON(fiber.Map{"status": "OK", "data": fdata, "users": fusers})
}
