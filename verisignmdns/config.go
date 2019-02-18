package verisignmdns

const (
	VerisignTokenParamName   = "VERISIGN_MDNS_API_TOKEN"
	VerisignUrlParamName     = "VERISIGN_MDNS_API_URL"
	VerisignDebugParamName   = "VERISIGN_MDNS_DEBUG"
	VerisignTimeoutParamName = "VERISIGN_MDNS_TIMEOUT"
	VerisignAccountParamName = "VERISIGN_ACCOUNT_ID"
	VerisignZoneParamName    = "VERISIGN_ZONE_NAME"
	userAgentFormat          = "terraform-provider-verisignmdns"
)

type Config struct {
	Token     string
	URL       string
	Debug     bool
	Timeout   int
	AccountId string
	ZoneName  string
}

func (c *Config) NewClient() (*api_client, error) {

	client, err := NewAPIClient(
		c.Token,
		c.URL,
		c.AccountId,
		c.ZoneName,
		c.Debug,
    c.Timeout,
	)
	if err != nil {
		return nil, err
	}

	return client, err
}
