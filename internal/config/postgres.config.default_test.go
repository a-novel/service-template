package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	postgrespresets "github.com/a-novel-kit/golib/postgres/presets"

	"github.com/a-novel/service-template/internal/config"
)

func TestNewPostgresPreset(t *testing.T) {
	t.Parallel()

	// The connection fields are golib's Connection.Options; the pool limits are this function's own.
	connection := postgrespresets.Connection{
		Host: "postgres", Port: 5432, User: "template", Password: "password", Database: "template",
	}
	preset := config.NewPostgresPreset(connection, 12, 8)

	require.Equal(t, 12, preset.MaxOpenConns)
	require.Equal(t, 8, preset.MaxIdleConns)
}
