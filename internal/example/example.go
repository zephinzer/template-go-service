package example

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/sirupsen/logrus"
)

func handleCustomResponse() (customResponse, error) {
	logrus.Infof("custom response request received")
	return customResponse{
		Number: rand.Intn(10),
	}, nil
}

func handleWithInput(input string) (string, error) {
	logrus.Infof("received input: %s", input)
	if strings.Contains(input, "sudo") {
		return "hello there!", nil
	}
	return "", fmt.Errorf("handleWithInput: %w", ErrorBadInput)
}

func handleNotFound() error {
	logrus.Infof("not found request received")
	return fmt.Errorf("handleNotFound: %w", ErrorNotFound)
}

func handleOk() error {
	logrus.Infof("ok")
	return nil
}
