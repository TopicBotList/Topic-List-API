// routes/get_token.go

package routes

import (
	"context"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.topiclist.xyz/types"
)

func GetToken(c *fiber.Ctx) error {
	token := c.Query("token")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	user := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{"status": "LOGIN"})
	}

	zippyCode := generateRandomCode(6)
	zippyExpires := time.Now().Add(300 * time.Second).Unix()

	payload := bson.M{
		"zippyCode":    zippyCode,
		"zippyExpires": zippyExpires,
	}

	_, err = usersCollection.UpdateOne(context.Background(), bson.M{"token": token}, bson.M{"$set": payload})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating user with zippy code",
		})
	}

	return c.JSON(fiber.Map{"code": zippyCode})
}

func generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
