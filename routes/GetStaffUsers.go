package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Define a function to handle the route
func GetStaffUsers(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	collection := db.Database("TopicBots").Collection("usersDB1")

	// Define a filter to find documents where staff is true
	filter := bson.M{"staff": "true"}

	// Find documents in the collection
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding documents",
		})
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and append each document to results
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding documents",
		})
	}

	// Remove the "token" field from each document
	for _, result := range results {
		delete(result, "token")
		delete(result, "_id")
		delete(result, "access_token")
	}

	// Return the results as JSON
	return c.JSON(results)
}
