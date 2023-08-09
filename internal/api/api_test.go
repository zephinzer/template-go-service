package api

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

type ApiTests struct {
	suite.Suite
}

func TestApi(t *testing.T) {
	suite.Run(t, &ApiTests{})
}

func (t ApiTests) Test_WithFiber() {
	app := fiber.New()
	t.Empty(app.GetRoutes())
	WithFiber(app)
	t.NotEmpty(app.GetRoutes())
}
