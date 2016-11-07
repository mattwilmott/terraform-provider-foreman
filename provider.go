package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["url"],
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["username"],
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["password"],
			},
		},

		/*ResourcesMap: map[string]*schema.Resource{
			"foreman_dns":    resourceDNS(),
		},*/

		ResourcesMap: map[string]*schema.Resource{
			"foreman_server": resourceServer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {

	descriptions = map[string]string{
		"url": "The Foreman server url. Example: \n" +
			"https://foreman.example.com/api/v2/",

		"username": "Foreman username with API access",

		"password": "Foreman password",
	}

}

// Returns the meta object (API connection/token)
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		URL:      d.Get("url").(string),
	}

	return config.Client()
}
