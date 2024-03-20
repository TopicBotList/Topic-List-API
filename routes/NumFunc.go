package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UserNum(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	userCollection := db.Database("TopicBots").Collection("usersDB1")

	countOptions := options.Count()
	totaluser, err := userCollection.CountDocuments(context.Background(), bson.M{}, countOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting servers",
		})
	}

	return c.JSON(fiber.Map{"total_user": totaluser})
}

func StaffNum(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	userCollection := db.Database("TopicBots").Collection("usersDB1")

	// Define a filter to count documents where staff is true
	filter := bson.M{"staff": "true"}

	// Count documents in the collection
	totalStaff, err := userCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting staff",
		})
	}

	return c.JSON(fiber.Map{"total_staff": totalStaff})
}

func UnapprovedNum(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	userCollection := db.Database("TopicBots").Collection("botsDB1")

	// Define a filter to count documents where staff is true
	filter := bson.M{"approved": "false"}

	// Count documents in the collection
	UnapprovedNum, err := userCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting unpprovedbots",
		})
	}

	return c.JSON(fiber.Map{"total_unapprovedbots": UnapprovedNum})
}

func BotsNum(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botsCollection := db.Database("TopicBots").Collection("botsDB1")

	countOptions := options.Count()
	totalbots, err := botsCollection.CountDocuments(context.Background(), bson.M{}, countOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting servers",
		})
	}

	return c.JSON(fiber.Map{"total_bots": totalbots})
}

func CountServers(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	serversCollection := db.Database("TopicBots").Collection("serversDB1")

	countOptions := options.Count()
	totalServers, err := serversCollection.CountDocuments(context.Background(), bson.M{}, countOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting servers",
		})
	}

	return c.JSON(fiber.Map{"total_servers": totalServers})
}
