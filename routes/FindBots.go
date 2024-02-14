package routes

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindBots(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botsCollection := db.Database("TopicBots").Collection("botsDB1")
	var data *mongo.Cursor
	var err error

	owner := c.Query("owner")
	limit := c.Query("limit")

	// Find bots based on owner if owner parameter is provided
	if owner != "" {
		data, err = botsCollection.Find(context.TODO(), bson.M{
			"publicity": "public",
			"owner":     owner,
		}, options.Find().SetSort(bson.D{{Key: "votes", Value: -1}}))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	} else if limit != "" {
		limitVal, err := strconv.Atoi(limit)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit parameter"})
		}

		data, err = botsCollection.Find(context.TODO(), bson.M{"publicity": "public"},
			options.Find().SetSort(bson.D{{Key: "votes", Value: -1}}).SetLimit(int64(limitVal)))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		data, err = botsCollection.Find(context.TODO(), bson.M{"publicity": "public"},
			options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	botsData := make([]interface{}, 0)
	for data.Next(context.Background()) {
		var bot map[string]interface{}
		if err := data.Decode(&bot); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		delete(bot, "_id")
		botsData = append(botsData, bot)
	}

	// Reset cursor to reuse it for the second loop
	data, err = botsCollection.Find(context.TODO(), bson.M{"publicity": "public"},
		options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	lastBotsData := make([]interface{}, 0)
	for data.Next(context.Background()) {
		var bot map[string]interface{}
		if err := data.Decode(&bot); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		delete(bot, "_id")
		lastBotsData = append(lastBotsData, bot)
	}

	return c.JSON(fiber.Map{"bots": botsData, "lbots": lastBotsData})
}
