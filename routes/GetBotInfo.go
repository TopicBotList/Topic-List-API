package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func BotRoute(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botsCollection := db.Database("TopicBots").Collection("botsDB1")
	analyticsCollection := db.Database("TopicBots").Collection("analyticsDB1")

	botID := c.Params("botid")
	userResult := botsCollection.FindOne(context.TODO(), bson.M{"id": botID})
	var userDocument bson.M
	if err := userResult.Decode(&userDocument); err != nil {
		return c.JSON(fiber.Map{"result": "none"})
	}
	if userDocument == nil {
		return c.JSON(fiber.Map{"result": "none"})
	}

	userDocument["_id"] = nil

	requester := c.Query("requester")
	owner, _ := userDocument["owner"].(string)
	approved, _ := userDocument["approved"].(bool)

	if !approved {
		if owner == requester {
			return c.JSON(fiber.Map{"result": userDocument})
		}

		// Assuming you have a collection called "admins" in your MongoDB
		adminsCollection := db.Database("YourDatabaseName").Collection("admins")
		count, err := adminsCollection.CountDocuments(context.TODO(), bson.M{"username": requester})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check admin status"})
		}
		if count == 0 {
			return c.JSON(fiber.Map{"result": "none"})
		}
	}

	var analyticsDocument bson.M
	anuserResult := analyticsCollection.FindOne(context.TODO(), bson.M{"id": userDocument["id"]})
	if err := anuserResult.Decode(&analyticsDocument); err != nil {
		_, err := analyticsCollection.InsertOne(context.TODO(), bson.M{"id": userDocument["id"], "views": int32(0), "impressions": int32(0)})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert analytics data"})
		}
		anuserResult := analyticsCollection.FindOne(context.TODO(), bson.M{"id": userDocument["id"]})
		if err := anuserResult.Decode(&analyticsDocument); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve analytics data"})
		}
	}

	analyticsDocument["_id"] = nil

	if c.Query("userview") != "" {
		viewsInt32, _ := analyticsDocument["views"].(int32)
		views := int(viewsInt32)
		_, err := analyticsCollection.UpdateOne(context.TODO(), bson.M{"id": userDocument["id"]}, bson.M{"$set": bson.M{"views": views + 1}})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update analytics data"})
		}
	}

	// Redact token before sending the response
	userDocument["token"] = "redacted"

	// Handle reviews which could be of type primitive.A
	if reviews, ok := userDocument["reviews"].(primitive.A); ok {
		for _, r := range reviews {
			if review, ok := r.(bson.M); ok {
				review["token"] = "redacted"
			}
		}
	}

	return c.JSON(fiber.Map{"result": userDocument, "analytics": analyticsDocument})
}
