package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecret(t *testing.T) {
	pwFile := path.Join(t.TempDir(), "password")
	require.NoError(t, os.WriteFile(pwFile, []byte("filepass\n"), 0600))
	t.Setenv("SG_TEST_PASSWORD", "envpass")

	conf := struct {
		Password Secret `yaml:"password"`
		FromFile Secret `yaml:"fromFile"`
		FromEnv  Secret `yaml:"fromEnv"`
		URI      Secret `yaml:"uri"`
	}{}
	input := strings.Join([]string{
		`password: "plainpass"`,
		`fromFile: "file:` + pwFile + `"`,
		`fromEnv: "env:SG_TEST_PASSWORD"`,
		`uri: "amqp://user:pass@localhost:5672"`,
	}, "\n")
	require.NoError(t, ParseConfig(strings.NewReader(input), &conf))

	t.Run("Test secret indirection", func(t *testing.T) {
		assert.Equal(t, Secret("plainpass"), conf.Password)
		assert.Equal(t, Secret("filepass"), conf.FromFile)
		assert.Equal(t, Secret("envpass"), conf.FromEnv)
	})

	t.Run("Test secret redaction", func(t *testing.T) {
		assert.Equal(t, "***", conf.Password.String())
		assert.Equal(t, "amqp://localhost:5672", conf.URI.String())
		// credentials must not leak through formatting of the whole config
		assert.NotContains(t, fmt.Sprintf("%v", conf), "pass@")
	})

	t.Run("Test missing secret file", func(t *testing.T) {
		err := ParseConfig(strings.NewReader(`password: "file:/nonexistent/password"`), &conf)
		assert.Error(t, err)
	})
}
