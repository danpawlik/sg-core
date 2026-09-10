package config

import (
	"net/url"
	"os"
	"strings"
)

// Secret is a config string holding a credential. It supports indirection so
// the value can be mounted from a Kubernetes Secret instead of being embedded
// in the config file:
//
//	password: "hunter2"                       # literal (discouraged)
//	password: "file:/run/secrets/es/password" # read from file, trailing newline trimmed
//	password: "env:ES_PASSWORD"               # read from environment
//
// Its String() is redacted, so a Secret logged via %v/%s (including
// logging.Metadata) never exposes the credential.
type Secret string

// UnmarshalYAML resolves file:/env: indirection at config parse time.
func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	switch {
	case strings.HasPrefix(v, "file:"):
		b, err := os.ReadFile(strings.TrimPrefix(v, "file:"))
		if err != nil {
			return err
		}
		v = strings.TrimRight(string(b), "\r\n")
	case strings.HasPrefix(v, "env:"):
		v = os.Getenv(strings.TrimPrefix(v, "env:"))
	}
	*s = Secret(v)
	return nil
}

// String returns a log-safe representation: for URIs the userinfo is stripped
// and the rest kept (useful for debugging), anything else becomes "***".
func (s Secret) String() string {
	if u, err := url.Parse(string(s)); err == nil && u.Host != "" {
		u.User = nil
		return u.String()
	}
	return "***"
}
