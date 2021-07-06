package goengage

import (
	"net/http"
	"os"
)

type (
	Config struct {
		Credentials *Credentials
		HTTPClient  *http.Client
	}
)

// NewConfig returns a blank config pointer. Call one of the With* methods to fill the config
func NewConfig() *Config {
	return &Config{}
}

// WithHttpClient sets a provided http client to to config
func (c *Config) WithHttpClient(client *http.Client) *Config {
	c.HTTPClient = client
	return c
}

// WithCredentials sets the config credential gotten from a credential provider
func (c *Config) WithCredentials(credentials *Credentials) *Config {
	c.Credentials = credentials
	return c
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type (
	Credentials struct {
		publicKey  string
		privateKey string
	}
)

// NewEnvCredentials reads and sets the API key and secret from runtime environment variables.
// It looks for ENGAGE_SO_PUBLIC_KEY and ENGAGE_SO_PRIVATE_KEY.
// The os env lookup sets blank string if the keys are not found.
func NewEnvCredentials() *Credentials {
	return &Credentials{
		publicKey:  os.Getenv("ENGAGE_SO_PUBLIC_KEY"),
		privateKey: os.Getenv("ENGAGE_SO_PRIVATE_KEY"),
	}
}

// NewStaticCredentials sets API key and secret from provided arguments
func NewStaticCredentials(publicKey, privateKey string) *Credentials {
	return &Credentials{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}
