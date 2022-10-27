package projectapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/lava/service"
)

func New() service.HttpRouter {
	return &handler{}
}

type handler struct {
}

func (h handler) Init() {
}

func (h handler) BasePrefix() string {
	return ""
}

func (h handler) Middlewares() []service.Middleware {
	return nil
}

func (h handler) Router(app *fiber.App) {
	app.Get("/test", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hello")
	})
}
