package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.topiclist.xyz/types"
)

func FindServers(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	serversCollection := db.Database("TopicBots").Collection("serversDB1")

	// Find top servers
	topCursor, err := serversCollection.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.M{"votes": -1}))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding top servers",
		})
	}
	defer topCursor.Close(context.Background())

	var ftop []types.Server
	for topCursor.Next(context.Background()) {
		var server types.Server
		err := topCursor.Decode(&server)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error decoding top server",
			})
		}

		// Find the owner user
		user := types.User{}
		err = usersCollection.FindOne(context.Background(), bson.M{"token": server.Owner}).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error finding owner user",
			})
		}

		if user.Token == server.Owner { // Check if the owner's token matches the server owner
			user.AccessToken = "Redacted" // Remove access token field
			server.Owner = "Redacted"     // Remove the owner field
			server.OwnerName = user.Name
			server.OwnerID = user.ID
			server.OwnerAvatar = user.Avatar
			ftop = append(ftop, server)
		}
	}

	// Find latest servers
	latestCursor, err := serversCollection.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.M{"_id": -1}))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding latest servers",
		})
	}
	defer latestCursor.Close(context.Background())

	var flatest []types.Server
	for latestCursor.Next(context.Background()) {
		var server types.Server
		err := latestCursor.Decode(&server)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error decoding latest server",
			})
		}

		// Find the owner user
		user := types.User{}
		err = usersCollection.FindOne(context.Background(), bson.M{"token": server.Owner}).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error finding owner user",
			})
		}

		if user.Token == server.Owner { // Check if the owner's token matches the server owner
			user.AccessToken = "Redacted" // Remove access token field
			server.Owner = "Redacted"     // Remove the owner field
			server.OwnerName = user.Name
			server.OwnerID = user.ID
			server.OwnerAvatar = user.Avatar
			flatest = append(flatest, server)
		}
	}

	// Find users who own servers
	cursor, err := serversCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding servers",
		})
	}
	defer cursor.Close(context.Background())

	var serverOwners []string
	for cursor.Next(context.Background()) {
		var server types.Server
		err := cursor.Decode(&server)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error decoding server",
			})
		}
		serverOwners = append(serverOwners, server.Owner)
	}

	var ffusers []types.User
	for _, ownerToken := range serverOwners {
		var user types.User
		err := usersCollection.FindOne(context.Background(), bson.M{"token": ownerToken}).Decode(&user)
		if err != nil {
			continue // If user not found, skip to the next owner
		}
		user.AccessToken = "Redacted" // Remove access token field
		ffusers = append(ffusers, user)
	}

	return c.JSON(fiber.Map{"top": ftop, "latest": flatest, "users": ffusers})
}
