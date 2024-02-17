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

func VoteForServer(c *fiber.Ctx) error {
	serverID := c.Params("serverid")
	token := c.Query("token")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")
	serversCollection := db.Database("TopicBots").Collection("serversDB1")
	votesCollection := db.Database("TopicBots").Collection("votesDB1")

	user := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"status": "LOGIN"})
	}

	data := types.Vote{}
	err = votesCollection.FindOne(context.Background(), bson.M{"token": token, "server": serverID}).Decode(&data)
	if err == nil {
		if data.End < time.Now().Unix() {
			votesCollection.DeleteOne(context.Background(), bson.M{"token": token, "server": serverID})
		} else {
			return c.JSON(fiber.Map{"status": "INVALID"})
		}
	}

	payload := types.Vote{
		Token:  token,
		Server: serverID,
		End:    time.Now().Unix() + 43200,
	}

	server := types.Server{}
	err = serversCollection.FindOne(context.Background(), bson.M{"id": serverID}).Decode(&server)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding server",
		})
	}

	_, err = serversCollection.UpdateOne(context.Background(),
		bson.M{"id": serverID},
		bson.M{"$set": bson.M{"votes": server.Votes + 1}})
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
