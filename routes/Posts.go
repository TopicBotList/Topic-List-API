package routes

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

// ServPostHandler handles requests for server posts
func ServPostHandler(c *fiber.Ctx) error {
	// Fetching the database client from context
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	servCollection := db.Database("TopicBots").Collection("ServPostsDB1")

	// Fetch all documents from the collection
	var servers []types.Post
	cursor, err := servCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error fetching server posts:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch server posts",
		})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var server types.Post
		if err := cursor.Decode(&server); err != nil {
			log.Println("Error decoding server post:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode server post",
			})
		}
		servers = append(servers, server)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error with cursor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cursor error",
		})
	}

	// Send the JSON response
	return c.JSON(servers)
}

// BotsPostHandler handles requests for bot posts
func BotsPostHandler(c *fiber.Ctx) error {
	// Fetching the database client from context
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botCollection := db.Database("TopicBots").Collection("BotPostsDB1")

	// Fetch all documents from the collection
	var bots []types.Post
	cursor, err := botCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error fetching bot posts:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch bot posts",
		})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var bot types.Post
		if err := cursor.Decode(&bot); err != nil {
			log.Println("Error decoding bot post:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode bot post",
			})
		}
		bots = append(bots, bot)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error with cursor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cursor error",
		})
	}

	// Send the JSON response
	return c.JSON(bots)
}
