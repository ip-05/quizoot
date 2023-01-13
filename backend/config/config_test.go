package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	cfg := InitConfig("../config.example")

	t.Run("TestConfigServer", func(t *testing.T) {
		assert.Equal(t, false, cfg.Server.Secure, "should be equal")
		assert.Equal(t, "localhost", cfg.Server.Domain, "should be equal")
		assert.Equal(t, "http://localhost:1234", cfg.Server.Base, "should be equal")
		assert.Equal(t, "localhost", cfg.Server.Host, "should be equal")
		assert.Equal(t, int64(1234), cfg.Server.Port, "should be equal")
	})

	t.Run("TestConfigGoogle", func(t *testing.T) {
		assert.Equal(t, "id", cfg.Google.ClientId, "should be equal")
		assert.Equal(t, "secret", cfg.Google.ClientSecret, "should be equal")
	})

	t.Run("TestConfigSecrets", func(t *testing.T) {
		assert.Equal(t, "jwt", cfg.Secrets.Jwt, "should be equal")
	})

	t.Run("TestConfigFrontend", func(t *testing.T) {
		assert.Equal(t, "http://localhost:4321", cfg.Frontend.Base, "should be equal")
	})

	t.Run("TestConfigDB", func(t *testing.T) {
		assert.Equal(t, "localhost", cfg.Database.Host, "should be equal")
		assert.Equal(t, int64(5432), cfg.Database.Port, "should be equal")
		assert.Equal(t, "user", cfg.Database.User, "should be equal")
		assert.Equal(t, "password", cfg.Database.Password, "should be equal")
		assert.Equal(t, "name", cfg.Database.DbName, "should be equal")
		assert.Equal(t, false, cfg.Database.Secure, "should be equal")
	})

	t.Run("TestConfigInvalid", func(t *testing.T) {
		assert.Panics(t, func() {
			InitConfig("test_panic")
		}, "should panic")
	})
}
