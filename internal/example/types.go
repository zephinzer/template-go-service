package example

import "app/internal/api/types"

// customResponse is the response data for the custom response
// controller
type customResponse struct {
	// Number is a numerical property
	Number int `json:"number"`
}

// customResponseHttpResponse is the response data structure
// for the custom response HTTP handler
type customResponseHttpResponse struct {
	types.HttpResponse

	// Data demonstrates defining a customResponse
	Data customResponse `json:"data"`
}
