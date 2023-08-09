package types

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HttpResponse struct {
	// Data contains any logical data to be returned to
	// the client
	Data any `json:"data,omitempty"`

	// Error when set indicates that the response is not
	// successful
	Error error `json:"error,omitempty"`

	// Message is an optional message that will be overridden
	// with the error message if an error occurred
	Message *string `json:"message,omitempty"`

	// StatusCode is the HTTP status code
	StatusCode int `json:"-"`
}

func (hr HttpResponse) WithFiber(c *fiber.Ctx) error {
	statusCode := http.StatusOK
	if hr.Error != nil {
		statusCode = http.StatusInternalServerError
		logrus.Warnf("http[%v] returned: %s", statusCode, hr.Error)
		errorMessage := hr.Error.Error()
		hr.Message = &errorMessage
		hr.Error = nil
	}
	if hr.StatusCode > 0 {
		statusCode = hr.StatusCode
	}
	c = c.Status(statusCode)
	return c.JSON(hr)
}
