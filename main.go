package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/ravener/discord-oauth2"
	"go.topiclist.xyz/configuration"
	"go.topiclist.xyz/database"
	"go.topiclist.xyz/routes"
	"golang.org/x/oauth2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "TopicList",
		AppName:       "A website used to list a Discord Server and Bots.",
	})

	config := configuration.GetConfig()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("config", config)
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
			"version": "5.0.0",
			"author":  "RanveerSoni",
			"links": fiber.Map{
				"status":  "https:/status.topiclist.xyz",
				"docs":    "https://docs.topiclist.xyz/",
				"support": "https://discord.gg/invite/Jad6TcdEet",
			},
		})
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://beta.topiclist.xyz,http://localhost:3000,https://servers.topiclist.xyz,https://topiclist.xyz",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	db, err := database.Connect(config.Database.Url)

	if err != nil {
		panic(err)
	} else {
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("db", db)
			return c.Next()
		})
	}

	v1 := app.Group("/")

	store := session.New()

	v1.Use(func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		c.Locals("session", sess)
		return c.Next()
	})

	v1.Use(func(c *fiber.Ctx) error {
		c.Locals("authConfig", &oauth2.Config{
			RedirectURL:  config.Client.Callback,
			ClientID:     config.Client.Id,
			ClientSecret: config.Client.Secret,
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
		})

		return c.Next()
	})

	v1.Get("/private/auth/login", routes.Login)
	v1.Get("/auth/callback", routes.Callback)
	v1.Get("/auth/logout", routes.Logout)
	v1.Get("/auth/@me", routes.GetCurrentUser)
	v1.Get("/sitemap", routes.SitemapHandler)
	v1.Get("/users/edit", routes.UserSettings)
	v1.Get("/users/settings", routes.UserSettings)
	v1.Get("/users/notifications", routes.UserNotifications)
	v1.Get("/private/server/all", routes.FindServers)
	v1.Get("/private/server/:serverid", routes.GetServer)
	v1.Get("/private/server/cat/:cat", routes.FindServersByCategory)
	v1.Get("/private/server/vote/:serverid", routes.VoteForServer)
	v1.Get("/private/server/:serverid/edit", routes.EditServer)
	v1.Get("/private/user/get", routes.GetUserInfo)
	v1.Get("/private/user/:userid", routes.GetUser)
	v1.Get("/private/user/edit", routes.EditUser)
	v1.Get("/private/zippy/token", routes.GetToken)
	v1.Get("/private/zippy/authorize", routes.AuthorizeZippy)
	v1.Get("/private/add", routes.AddServer)

	app.Listen(":" + config.Web.Port)
}
