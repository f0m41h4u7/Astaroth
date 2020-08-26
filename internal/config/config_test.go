package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		_, err := ReadConfig("../../tests/testdata/config.json")
		require.Nil(t, err)
	})
}
