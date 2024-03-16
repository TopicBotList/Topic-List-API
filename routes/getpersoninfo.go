package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InfoRoute(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	assetsCollection := db.Database("TopicBots").Collection("assetsDB1")
	itemsCollection := db.Database("TopicBots").Collection("itemsDB1")

	userID := c.Query("user")
	userResult := usersCollection.FindOne(context.TODO(), bson.M{"id": userID})
	var userDocument bson.M
	if err := userResult.Decode(&userDocument); err != nil {
		return c.JSON(fiber.Map{"result": "none"})
	}

	if userDocument == nil {
		return c.JSON(fiber.Map{"result": "none"})
	}

	userDocument["_id"] = nil

	userAssets := make([]bson.M, 0)

	cursor, err := assetsCollection.Find(context.TODO(), bson.M{"owner": userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user assets"})
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.Background()) {
		var asset bson.M
		if err := cursor.Decode(&asset); err != nil {
			continue
		}

		itemResult := itemsCollection.FindOne(context.TODO(), bson.M{"id": asset["id"]})
		var itemDocument bson.M
		if err := itemResult.Decode(&itemDocument); err != nil {
			assetsCollection.DeleteMany(context.TODO(), bson.M{"id": asset["id"]})
			continue
		}

		itemDocument["_id"] = nil
		userAssets = append(userAssets, itemDocument)
	}

	userDocument["token"] = nil
	userDocument["access_token"] = nil

	return c.JSON(fiber.Map{"result": userDocument, "assets": userAssets})
}
