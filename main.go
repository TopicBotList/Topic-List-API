package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/ravener/discord-oauth2"
	swagger "github.com/swaggo/fiber-swagger"
	"go.topiclist.xyz/configuration"
	"go.topiclist.xyz/database"
	_ "go.topiclist.xyz/docs"
	"go.topiclist.xyz/routes"
	"golang.org/x/oauth2"
)

//	@title			Topic-List API
//	@version		5.0
//	@description	A simple API for TopicList which handles all requests from TopicServers and TopicBots.
//	@termsOfService	https://topiclist.xyz/legal/tos

//	@contact.name	API Support
//	@contact.url	https://discord.gg/GJGbMXENtp
//	@contact.email	support@topiclist.xyz

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		api.topiclist.xyz
//	@BasePath	/swagger/index.html

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
		AllowOrigins:     "https://beta.topiclist.xyz,https://beta.topiclist.xyz,http://localhost:3001,http://localhost:3000,https://topic-bots.vercel.app,https://server.topiclist.xyz,https://servers.topiclist.xyz,https://topiclist.xyz",
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
	v1.Get("/private/server/cat/:cat", routes.FindServersByCategory)
	v1.Get("/private/server/vote/:serverid", routes.VoteServ)
	v1.Get("/private/server/:serverid/edit", routes.EditServer)

	//Add
	v1.Get("/private/add", routes.AddServer)

	//Reviews

	v1.Get("/reviews/:botid/add", routes.AddReview) //works
	v1.Get("/reviews/:botid/delete", routes.DeleteReview)

	//Users
	v1.Get("/private/user/get", routes.GetUserInfo)
	v1.Get("/private/user/:userid", routes.GetUser)
	v1.Get("/private/user/edit", routes.EditUser)
	v1.Get("/users/edit", routes.UserSettings)
	v1.Get("/users/settings", routes.UserSettings)
	v1.Get("/users/notifications", routes.UserNotifications)

	//Zippy
	v1.Get("/private/zippy/token", routes.GetToken)
	v1.Get("/private/zippy/authorize", routes.AuthorizeZippy)

	//Docs
	app.Get("/swagger/*", swagger.WrapHandler)

	//Bots
	v1.Get("/find_bots", routes.FindBots)
	v1.Get("/editbot/settings", routes.EditBotSettings)
	v1.Get("/delete/:botid", routes.DeleteBot)
	v1.Get("/users/edit", routes.UserEditBots)
	v1.Get("/bot", routes.BotRoute)
	v1.Get("/info", routes.InfoRoute)
	v1.Get("/vote/:botid", routes.VoteBot)
	v1.Get("/explore", routes.Explore)

	// Listen and serve
	port := configuration.GetConfig().Web.Port
	log.Fatal(app.Listen(":" + port))
}
