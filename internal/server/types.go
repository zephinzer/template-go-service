package server

import (
	"github.com/gofiber/fiber/v2"
)

type Event struct {
	Data    any
	Message string
	Error   any
}

type Healthcheck func() error

type Logger func(content ...any)
type Loggerf func(format string, args ...any)

type Router func(app *fiber.App)
