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
	Trace   string
	Timeout int
}

func (c *Config) NewClient() (*arukasClient, error) {

	timeout := time.Duration(0)
	if c.Timeout > 0 {
		timeout = time.Duration(c.Timeout) * time.Second
	}

	client, err := NewAPIClient(
		c.Token,
		c.URL,
		c.Trace != "",
    c.Timeout,
	})
	if err != nil {
		return nil, err
	}

	return client, err
}
