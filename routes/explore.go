package routes

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Explore(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	sortName := c.FormValue("sort[name]")
	resultsNameStr := c.FormValue("results[name]")
	resultsName, err := strconv.Atoi(resultsNameStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid results[name]",
		})
	}
	tagsName := c.FormValue("tags[name]")

	var filter bson.M
	switch sortName {
	case "votes":
		filter = bson.M{"publicity": "public"}
	case "latest":
		filter = bson.M{"publicity": "public"}
	case "oldest":
		filter = bson.M{"publicity": "public"}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sort parameter",
		})
	}

	collection := db.Database("TopicBots").Collection("botsDB1")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching data",
		})
	}
	defer cursor.Close(context.Background())

	var finaldata []bson.M
	for cursor.Next(context.Background()) {
		var bot bson.M
		if err := cursor.Decode(&bot); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error decoding data",
			})
		}

		if len(finaldata) >= resultsName {
			break
		}

		if _, ok := bot["_id"]; ok {
			delete(bot, "_id")
		}

		if tagsName == "all" || contains(bot["tags"], tagsName) {
			finaldata = append(finaldata, bot)
		}
	}

	return c.JSON(fiber.Map{"bots": finaldata})
}

func contains(tags interface{}, tag string) bool {
	tagList, ok := tags.([]string)
	if !ok {
		return false
	}
	for _, t := range tagList {
		if t == tag {
			return true
		}
	}
	return false
}
