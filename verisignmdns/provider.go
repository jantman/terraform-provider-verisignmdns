package verisignmdns

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(VerisignTokenParamName, nil),
				Description: "your Verisign MDNS ReST API token",
			},
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(VerisignAccountParamName, nil),
				Description: "Verisign MDNS Account ID to act on",
			},
			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(VerisignZoneParamName, nil),
				Description: "Zone name to manage records in",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(VerisignUrlParamName, "https://mdns.verisign.com/mdns-web/api/"),
				Description: "default Verisign MDNS API url",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(VerisignDebugParamName, false),
				Description: "Enabling this will cause lots of debug information to be printed to STDOUT by the API client.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(VerisignTimeoutParamName, "900"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"verisignmdns_rr": resourceRr(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		Token:     d.Get("token").(string),
		URL:       d.Get("api_url").(string),
		Debug:     d.Get("debug").(bool),
		Timeout:   d.Get("timeout").(int),
		AccountId: d.Get("account_id").(string),
		ZoneName:  d.Get("zone_name").(string),
	}

	return config.NewClient()
}
