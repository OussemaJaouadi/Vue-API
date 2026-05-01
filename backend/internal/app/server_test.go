package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeRequestURI(t *testing.T) {
	require.Equal(t, "/events?ticket=[redacted]", sanitizeRequestURI("/events?ticket=secret-token"))
	require.Equal(t, "/events?foo=bar&ticket=[redacted]", sanitizeRequestURI("/events?foo=bar&ticket=secret-token"))
	require.Equal(t, "/healthz", sanitizeRequestURI("/healthz"))
	require.Equal(t, "/bad%zz", sanitizeRequestURI("/bad%zz"))
}
