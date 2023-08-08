package api

import (
	"app/internal/example"

	"github.com/gofiber/fiber/v2"
)

// WithFiber godoc
// @title Example HTTP API
// @version 1.0.0
// @description This is an API
// @BasePath /
func WithFiber(app *fiber.App) {
	example.WithFiber(app)
}
