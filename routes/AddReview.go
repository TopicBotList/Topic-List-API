package routes

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddReview adds a review for a bot.
//	@Summary		Add a review for a bot
//	@Description	Add a review for a bot identified by botid
//	@Tags			Reviews
//	@Accept			json
//	@Produce		json
//	@Param			token	query		string	true	"User token"
//	@Param			botid	path		string	true	"Bot ID"
//	@Param			content	formData	string	true	"Review content"
//	@Success		200		{object}	fiber.Map{"reply": "OK"}
//	@Failure		400		{object}	fiber.Map{"reply": "TOKEN_INVALID"}
//	@Failure		400		{object}	fiber.Map{"reply": "BOT_INVALID"}
//	@Failure		500		{object}	fiber.Map{"error": "Database connection not available"}
//	@Router			/reviews/{botid} [post]

func AddReview(c *fiber.Ctx) error {
	token := c.Query("token")
	botID := c.Params("botid")
	content := c.FormValue("content")

	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	botsCollection := db.Database("TopicBots").Collection("botsDB1")
	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	var user map[string]interface{}
	err := usersCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"reply": "TOKEN_INVALID"})
	}

	var data map[string]interface{}
	err = botsCollection.FindOne(context.TODO(), bson.M{"id": botID}).Decode(&data)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"reply": "BOT_INVALID"})
	}

	// Handling the reviews field correctly
	var reviews bson.A
	reviewsInterface := data["reviews"]
	if reviewsInterface == nil {
		reviews = bson.A{}
	} else {
		var ok bool
		reviews, ok = reviewsInterface.(bson.A)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Reviews field type assertion failed"})
		}
	}

	review := map[string]interface{}{
		"owner":   user["id"],
		"avatar":  user["avatar"],
		"name":    user["name"],
		"content": content,
		"token":   token,
		"id":      generateID(),
	}
	reviews = append(reviews, review)

	_, err = botsCollection.UpdateOne(context.TODO(), bson.M{"id": botID},
		bson.M{"$set": bson.M{"reviews": reviews}}, options.Update().SetUpsert(true))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"reply": "OK"})
}

// generateID generates a random string ID
func generateID() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	id := make([]byte, 10)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)

}
