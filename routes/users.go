// routes/users.go

package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.topiclist.xyz/types"
)

func checkToken(token string) bool {
	// Implement your token validation logic here
	// Return true if the token is valid, false otherwise
	return true
}

func UserEdit(c *fiber.Ctx) error {
	token := c.Query("token")
	if !checkToken(token) {
		return c.JSON(fiber.Map{"reply": "TOKEN_INVALID"})
	}

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	data := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&data)
	if err != nil {
		return c.JSON(fiber.Map{"reply": "TOKEN_INVALID"})
	}

	payload := bson.M{"bio": c.FormValue("bio")}
	_, err = usersCollection.UpdateOne(context.Background(), bson.M{"token": token}, bson.M{"$set": payload})
	if err != nil {
		return c.JSON(fiber.Map{"reply": "error updating user"})
	}

	return c.JSON(fiber.Map{"reply": "worked"})
}

func UserSettings(c *fiber.Ctx) error {
	token := c.Query("token")
	if !checkToken(token) {
		return c.JSON(fiber.Map{"reply": "TOKEN_INVALID"})
	}

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	data := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&data)
	if err != nil {
		return c.JSON(fiber.Map{"reply": "TOKEN_INVALID"})
	}

	fdata := make(map[string]string)
	// Note: You may customize the field names and types based on your actual User struct.
	payloadMap := bson.M{"field1": fdata["field1"], "field2": fdata["field2"]}
	_, err = usersCollection.UpdateOne(context.Background(), bson.M{"token": token}, bson.M{"$set": payloadMap}, options.Update().SetUpsert(true))
	if err != nil {
		return c.JSON(fiber.Map{"reply": "error updating user settings"})
	}

	data = types.User{}
	err = usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&data)
	if err != nil {
		return c.JSON(fiber.Map{"reply": "error fetching updated user data"})
	}

	data.ID = "" // Remove MongoDB's "_id" field
	return c.JSON(fiber.Map{"reply": "worked", "newdata": data})
}

func UserNotifications(c *fiber.Ctx) error {
	token := c.Query("token")
	if !checkToken(token) {
		return c.JSON(fiber.Map{"reply": "TOKEN_INVALID"})
	}

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	data := types.User{}
	err := usersCollection.FindOne(context.Background(), bson.M{"token": token}).Decode(&data)
	if err != nil {
		return c.JSON(fiber.Map{"reply": "TOKEN_INVALID"})
	}
	payloadMap := bson.M{"your_field": "your_value"}

	updateOptions := options.Update().SetUpsert(true)
	_, err = usersCollection.UpdateOne(context.Background(), bson.M{"token": token}, bson.M{"$set": payloadMap}, updateOptions)
	if err != nil {
		return c.JSON(fiber.Map{"reply": "error updating user notifications"})
	}

	return c.JSON(fiber.Map{"reply": data.Notifications})
}
