// routes/sitemap.go

package routes

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func SitemapHandler(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botsCollection := db.Database("TopicBots").Collection("botsDB1")
	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	// Fetch bots
	botData, err := getBots(botsCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching bot data",
		})
	}

	// Fetch users
	userData, err := getUsers(usersCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching user data",
		})
	}

	// Filter and extract relevant information
	fbotData := make([]string, 0)
	for _, bot := range botData {
		if bot.Votes {
			fbotData = append(fbotData, bot.ID)
		}
	}

	fuserData := make([]string, 0)
	for _, user := range userData {
		fuserData = append(fuserData, user.ID)
	}

	// Construct and return the JSON response
	response := map[string]interface{}{
		"bots":  fbotData,
		"users": fuserData,
	}

	return c.JSON(response)
}

func getBots(collection *mongo.Collection) ([]types.Bots, error) {
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		// Add logging to see the actual error
		fmt.Println("Error fetching bots:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var bots []types.Bots
	if err := cursor.All(context.Background(), &bots); err != nil {
		// Add logging to see the actual error
		fmt.Println("Error decoding bots:", err)
		return nil, err
	}

	return bots, nil
}

func getUsers(collection *mongo.Collection) ([]types.User, error) {
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []types.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	return users, nil
}
