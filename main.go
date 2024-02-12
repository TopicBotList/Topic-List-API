package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/ravener/discord-oauth2"
	"go.topiclist.xyz/configuration"
	"go.topiclist.xyz/database"
	"go.topiclist.xyz/routes"
	"golang.org/x/oauth2"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create Fiber app instance
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "TopicList",
		AppName:       "A website used to list a Discord Server and Bots.",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println("Error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	// Middleware: CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://beta.topiclist.xyz,http://localhost:3000,https://topic-bots.vercel.app,https://server.topiclist.xyz,https://servers.topiclist.xyz,https://topiclist.xyz",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Middleware: Database Connection
	db, err := database.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Disconnect(nil)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// Middleware: Session
	store := session.New()
	app.Use(func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		c.Locals("session", sess)
		return c.Next()
	})

	// Middleware: OAuth2 Configuration
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("authConfig", &oauth2.Config{
			RedirectURL:  configuration.GetConfig().Client.Callback,
			ClientID:     configuration.GetConfig().Client.Id,
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
		})
		return c.Next()
	})

	// Routes
	v1 := app.Group("/")
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
			"version": "5.0.0",
			"author":  "RanveerSoni",
			"links": fiber.Map{
				"status":  "https://status.topiclist.xyz",
				"docs":    "https://docs.topiclist.xyz/",
				"support": "https://discord.gg/invite/Jad6TcdEet",
			},
		})
	})
	//Shared
	v1.Get("/auth/login", routes.Login)
	v1.Get("/auth/callback", routes.Callback)
	v1.Get("/auth/logout", routes.Logout)
	v1.Get("/auth/@me", routes.GetCurrentUser)
	//Servers
	v1.Get("/private/server/all", routes.FindServers)
	v1.Get("/private/server/:serverid", routes.GetServer)
	v1.Get("/private/server/cat/:cat", routes.FindServersByCategory) //unknown state
	v1.Get("/private/server/vote/:serverid", routes.VoteForServer)   //unknown state
	v1.Get("/private/server/:serverid/edit", routes.EditServer)      //unknown state
	v1.Get("/private/user/get", routes.GetUserInfo)                  //unknown state
	v1.Get("/private/user/:userid", routes.GetUser)
	v1.Get("/private/user/edit", routes.EditUser)
	v1.Get("/private/zippy/token", routes.GetToken)           //unknown state
	v1.Get("/private/zippy/authorize", routes.AuthorizeZippy) //unknown state
	v1.Get("/private/add", routes.AddServer)
	//Bots
	v1.Get("/find_bots", routes.FindBots)
	v1.Get("/reviews/:botid/add", routes.AddReview)
	v1.Get("/reviews/:botid/delete", routes.DeleteReview)
	v1.Get("/editbot/settings", routes.EditBotSettings)
	v1.Get("/delete/:botid", routes.DeleteBot)
	v1.Get("/sitemap", routes.SitemapHandler)
	v1.Get("/users/edit", routes.UserEditBots)
	v1.Get("/bot", routes.BotRoute)
	v1.Get("/users/edit", routes.UserSettings)
	v1.Get("/users/settings", routes.UserSettings)
	v1.Get("/users/notifications", routes.UserNotifications)
	v1.Get("/info", routes.InfoRoute)

	// Listen and serve
	port := configuration.GetConfig().Web.Port
	log.Fatal(app.Listen(":" + port))
}
