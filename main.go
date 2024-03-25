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
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Origin, X-Requested-With, Accept,x-client-key, x-client-token, x-client-secret, authorization",
		AllowCredentials: false,
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
	app.Get("/", func(c *fiber.Ctx) error {
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
	app.Get("/auth/login", routes.Login)
	app.Get("/auth/callback", routes.Callback)
	app.Get("/auth/logout", routes.Logout)
	app.Get("/auth/@me", routes.GetCurrentUser)

	//Servers
	app.Get("/private/server/@all", routes.FindServers)
	app.Get("/private/get/guilds", routes.GetGuilds)
	app.Get("/private/server/:serverid", routes.GetServer)
	app.Get("/private/server/cat/:cat", routes.FindServersByCategory)
	app.Get("/private/server/:serverid/edit", routes.EditServer)

	//Bots
	app.Get("/find_bots", routes.FindBots)
	app.Post("/editbot/settings", routes.EditBotSettings)
	app.Delete("/delete/:botid", routes.DeleteBot)
	app.Get("/bot/:botid", routes.BotRoute)
	app.Get("/info", routes.InfoRoute)
	app.Get("/explore", routes.Explore)

	//Add
	app.Post("/private/add", routes.AddServer)

	//Reviews
	app.Post("/reviews/:botid/add", routes.AddReview)
	app.Delete("/reviews/:botid/delete", routes.DeleteReview)

	//Users
	app.Get("/private/user/get", routes.GetUserInfo)
	app.Get("/private/user/:userid", routes.GetUser)
	app.Post("/private/user/edit", routes.EditUser)
	app.Post("/users/edit", routes.UserSettings)
	app.Post("/users/edit", routes.UserEditBots)
	app.Get("/users/settings", routes.UserSettings)
	app.Get("/users/notifications", routes.UserNotifications)

	//Zippy
	app.Get("/private/zippy/token", routes.GetToken)
	app.Post("/private/zippy/authorize", routes.AuthorizeZippy)

	// Admin Utils
	app.Get("/botnum", routes.BotsNum)
	app.Get("/usernum", routes.UserNum)
	app.Get("/servnum", routes.CountServers)
	app.Get("/staffnum", routes.StaffNum)
	app.Get("/team", routes.GetStaffUsers)
	app.Get("/unapprovedbotsnum", routes.UnapprovedNum)

	//Partner
	app.Get("/partners/@all", routes.GetAllPartner)

	//Posts
	app.Get("/posts/bots/@all", routes.BotsPostHandler)
	app.Get("/posts/serv/@all", routes.ServPostHandler)

	//Vote
	app.Post("/vote/:botid", routes.VoteBot)
	app.Post("/private/server/vote/:serverid", routes.VoteServ)

	// Listen and serve
	port := configuration.GetConfig().Web.Port
	log.Fatal(app.Listen(":" + port))
}
