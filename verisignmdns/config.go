package verisignmdns

import (
	"fmt"
	"time"
)

const (
	VerisignTokenParamName   = "VERISIGN_MDNS_API_TOKEN"
	VerisignUrlParamName     = "VERISIGN_MDNS_API_URL"
	VerisignDebugParamName   = "VERISIGN_MDNS_DEBUG"
	VerisignTimeoutParamName = "VERISIGN_MDNS_TIMEOUT"
	userAgentFormat          = "terraform-provider-verisignmdns"
)

type Config struct {
	Token   string
	URL     string
	Debug   bool
	Timeout int
}

func (c *Config) NewClient() (*api_client, error) {

	timeout := time.Duration(0)
	if c.Timeout > 0 {
		timeout = time.Duration(c.Timeout) * time.Second
	}

	client, err := NewAPIClient(
		c.Token,
		c.URL,
		c.Debug,
    c.Timeout,
	)
	if err != nil {
		return nil, err
	}

	return client, err
}
