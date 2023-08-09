package example

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExampleTests struct {
	suite.Suite
}

func TestExample(t *testing.T) {
	suite.Run(t, &ExampleTests{})
}

func (t ExampleTests) Test_handleCustomResponse() {
	observed, err := handleCustomResponse()
	t.Nil(err)
	t.LessOrEqual(observed.Number, 10)
	t.GreaterOrEqual(observed.Number, 0)
}

func (t ExampleTests) Test_handleWithInput() {
	observed, err := handleWithInput("sudo pass")
	t.Nil(err)
	t.EqualValues("hello there!", observed)
}

func (t ExampleTests) Test_handleNotFound() {
	err := handleNotFound()
	t.NotNil(err)
	t.Contains(err.Error(), "handleNotFound")
}

func (t ExampleTests) Test_handleOk() {
	err := handleOk()
	t.Nil(err)
}
