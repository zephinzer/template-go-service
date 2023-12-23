package mongo

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type ConfigTests struct {
	suite.Suite
}

func TestConfig(t *testing.T) {
	suite.Run(t, &ConfigTests{})
}

func (t ConfigTests) Test_WithCobra() {
	c := cobra.Command{}
	t.False(c.HasFlags())
	WithCobra(&c)
	t.True(c.HasFlags())
}
