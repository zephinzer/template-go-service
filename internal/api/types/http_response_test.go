package types

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type HttpResponseTests struct {
	suite.Suite
}

func TestApi(t *testing.T) {
	suite.Run(t, &HttpResponseTests{})
}

func (t HttpResponseTests) Test_WithFiber() {
	app := fiber.New()
	bodyData := "hello world"
	errorMessage := "error"
	customMessage := "this will be overridden by the error"
	statusCode := http.StatusTeapot

	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	r := HttpResponse{
		Data:       bodyData,
		Error:      errors.New(errorMessage),
		Message:    &customMessage,
		StatusCode: statusCode,
	}
	r.WithFiber(c)
	t.EqualValues(statusCode, c.Response().StatusCode())
	t.EqualValues(
		fmt.Sprintf(`{"data":"%s","message":"%s"}`, bodyData, errorMessage),
		string(c.Response().Body()),
		"message field should be overridden if there are errors",
	)
	r = HttpResponse{
		Data:       bodyData,
		Message:    &customMessage,
		StatusCode: statusCode,
	}
	r.WithFiber(c)
	t.EqualValues(statusCode, c.Response().StatusCode())
	t.EqualValues(
		fmt.Sprintf(`{"data":"%s","message":"%s"}`, bodyData, customMessage),
		string(c.Response().Body()),
		"message field should not be overridden if there are no errors",
	)
}
