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
                  "api_url": {
                           Type:        schema.TypeString,
                           Optional:    true,
                           DefaultFunc: schema.EnvDefaultFunc(VerisignUrlParamName, "https://mdns.verisign.com/mdns-web/api/"),
                           Description: "default Arukas API url",
                  },
                  "trace": {
                         Type:        schema.TypeString,
                         Optional:    true,
                         DefaultFunc: schema.EnvDefaultFunc(VerisignDebugParamName, ""),
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
		Token:   d.Get("token").(string),
		URL:     d.Get("api_url").(string),
		Trace:   d.Get("trace").(string),
		Timeout: d.Get("timeout").(int),
	}

  return config.NewClient()
}
