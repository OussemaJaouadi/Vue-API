package id_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/id"
)

func TestNewUUIDV7ReturnsUUIDV7(t *testing.T) {
	value := id.NewUUIDV7()

	uuidV7Pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	require.Regexp(t, uuidV7Pattern, value)
}
