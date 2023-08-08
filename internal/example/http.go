package example

import (
	"app/internal/api/types"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type httpHandler func(fiber.Router)

func WithFiber(app fiber.Router) {
	addHandlers := []httpHandler{
		customResponseHttpHandler,
		notFoundHttpHandler,
		okHttpHandler,
		pathInputHttpHandler,
	}
	for _, addHandler := range addHandlers {
		addHandler(app)
	}
}

// customResponseHttpHandler godoc
// @Summary Example endpoint returning 200
// @Description Returns 200 and demonstrates documenting a custom response with Swagger
// @Tags 200
// @Produce json
// @Success 200 {object} customResponseHttpResponse
// @Router /api/example/200/custom-response [get]
func customResponseHttpHandler(app fiber.Router) {
	app.Get("/api/example/200/custom-response", func(c *fiber.Ctx) error {
		httpResponse := types.HttpResponse{}
		response, err := handleCustomResponse()
		if err != nil {
			httpResponse.Error = fmt.Errorf("failed to complete: %s", err)
			httpResponse.StatusCode = getHttpStatusCode(err)
		} else {
			httpResponse.Data = response
		}
		return httpResponse.WithFiber(c)
	})
}

// notFoundHttpHandler godoc
// @Summary Example endpoint returning 404
// @Description Returns 404 and demonstrates setting of status code
// @Tags 404
// @Produce json
// @Failure 404 {object} types.HttpResponse
// @Router /api/example/404 [get]
func notFoundHttpHandler(app fiber.Router) {
	app.Get("/api/example/404", func(c *fiber.Ctx) error {
		httpResponse := types.HttpResponse{}
		if err := handleNotFound(); err != nil {
			httpResponse.Error = fmt.Errorf("failed to complete: %s", err)
			httpResponse.StatusCode = getHttpStatusCode(err)
		}
		return httpResponse.WithFiber(c)
	})
}

// okHttpHandler godoc
// @Summary Example endpoint returning 200
// @Description Returns 200 and demonstrates a basic HTTP handler
// @Tags 200
// @Produce json
// @Success 200 {object} types.HttpResponse
// @Router /api/example/200 [get]
func okHttpHandler(app fiber.Router) {
	app.Get("/api/example/200", func(c *fiber.Ctx) error {
		httpResponse := types.HttpResponse{}
		if err := handleOk(); err != nil {
			httpResponse.Error = fmt.Errorf("failed to complete: %s", err)
			httpResponse.StatusCode = getHttpStatusCode(err)
		} else {
			message := "ok"
			httpResponse.Message = &message
		}
		return httpResponse.WithFiber(c)
	})
}

// pathInputHttpHandler godoc
// @Summary Example endpoint returning 200
// @Description Returns 200 and demonstrates path parameter usage
// @Tags 200
// @Produce json
// @Param input path string true "Any input you want"
// @Success 200 {object} types.HttpResponse
// @Router /api/example/200/with/{input} [get]
func pathInputHttpHandler(app fiber.Router) {
	app.Get("/api/example/200/with/:input", func(c *fiber.Ctx) error {
		httpResponse := types.HttpResponse{}
		response, err := handleWithInput(c.Params("input"))
		if err != nil {
			httpResponse.Error = fmt.Errorf("failed to complete: %s", err)
			httpResponse.StatusCode = getHttpStatusCode(err)
		} else {
			httpResponse.Message = &response
		}
		return httpResponse.WithFiber(c)
	})
}
