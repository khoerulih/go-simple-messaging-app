package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/khoerulih/go-simple-messaging-app/app/ws"
	"github.com/khoerulih/go-simple-messaging-app/pkg/database"
	"github.com/khoerulih/go-simple-messaging-app/pkg/env"
	"github.com/khoerulih/go-simple-messaging-app/pkg/router"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	database.SetupDatabase()
	database.SetupMongoDB()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())

	go ws.ServeWSMessaging(app)

	router.InstallRouter(app)

	return app
}
