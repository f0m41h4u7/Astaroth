package collector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCPU(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cpu, err := readCPUFile("../../tests/testdata/cpu.txt")
		require.Nil(t, err)
		require.Equal(t, 42055566250653, cpu)
	})
	t.Run("wrong file", func(t *testing.T) {
		_, err := readCPUFile("../../tests/testdata/bad_file.txt")
		require.NotNil(t, err)
	})
	t.Run("nonexistent file", func(t *testing.T) {
		_, err := readCPUFile("cpu.txt")
		require.NotNil(t, err)
	})
}

func TestLoadAvg(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		la, err := readLoadAvgFile("../../tests/testdata/loadavg.txt")
		require.Nil(t, err)
		require.Equal(t, "0.94", la[0])
		require.Equal(t, "1.49", la[1])
		require.Equal(t, "1.63", la[2])
		require.Equal(t, "1", la[3])
		require.Equal(t, "924", la[4])
	})
	t.Run("wrong file", func(t *testing.T) {
		_, err := readLoadAvgFile("../../tests/testdata/bad_file.txt")
		require.NotNil(t, err)
	})
	t.Run("nonexistent file", func(t *testing.T) {
		_, err := readLoadAvgFile("loadavg.txt")
		require.NotNil(t, err)
	})
}
