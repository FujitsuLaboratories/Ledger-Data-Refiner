package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//go:generate sed -i -e /http_port/d ../config/config.ini
func TestSettings(t *testing.T) {
	require.Equal(t, AppMode, "release")
	require.Equal(t, LogOutput, "stdout")
	require.Equal(t, DBName, "ledgerdata_refiner")
	require.Equal(t, NetworkName, "test-fabric")

	// Default value
	require.Equal(t, HttpPort, ":9999")
	// Override with environment variable
	// require.Equal(t, DBPort, "15432")
}
