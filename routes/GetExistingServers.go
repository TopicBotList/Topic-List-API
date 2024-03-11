package routes

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetGuilds(c *fiber.Ctx) error {
	// Retrieve the MongoDB client from fiber context
	db, ok := c.Locals("db").(*mongo.Client)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not available",
		})
	}

	// Access the users collection
	usersCollection := db.Database("TopicBots").Collection("usersDB1")

	// Fetch the token from usersDB1
	var user struct {
		Token string `bson:"token"`
	}
	err := usersCollection.FindOne(context.Background(), bson.M{}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	token := user.Token

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/users/@me/guilds", nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if resp.StatusCode != http.StatusOK {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch guilds"})
	}

	var guilds []Guild
	if err := json.Unmarshal(body, &guilds); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"guilds": guilds})
}
