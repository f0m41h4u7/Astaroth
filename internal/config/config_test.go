package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		err := InitConfig("../../tests/testdata/config.json")
		require.Nil(t, err)
		res := string(RequiredMetrics.Metrics[LoadAvg])
		require.Equal(t, "on", res)
	})
}
