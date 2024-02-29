// routes/vote.go

package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func VoteBot(c *fiber.Ctx) error {
	botID := c.Params("botid")
	token := c.Query("token")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	botsCollection := db.Database("TopicBots").Collection("BotsDB1")
	votesCollection := db.Database("TopicBots").Collection("votesDB1")

	user := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"status": "LOGIN"})
	}

	data := types.Vote{}
	err = votesCollection.FindOne(context.Background(), bson.M{"token": token, "bot": botID}).Decode(&data)
	if err == nil {
		if data.End < time.Now().Unix() {
			votesCollection.DeleteOne(context.Background(), bson.M{"token": token, "bot": botID})
		} else {
			return c.JSON(fiber.Map{"status": "INVALID"})
		}
	}

	payload := types.Vote{
		Token: token,
		Bot:   botID,
		End:   time.Now().Unix() + 43200,
	}

	bot := types.Bots{}
	err = botsCollection.FindOne(context.Background(), bson.M{"id": botID}).Decode(&bot)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding server",
		})
	}

	_, err = botsCollection.UpdateOne(context.Background(),
		bson.M{"id": botID},
		bson.M{"$set": bson.M{"votes": bot.Votes + 1}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating server votes",
		})
	}

	_, err = votesCollection.InsertOne(context.Background(), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error inserting vote",
		})
	}

	return c.JSON(fiber.Map{"status": "OK"})
}
