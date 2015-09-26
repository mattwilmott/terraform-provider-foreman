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

		ResourcesMap: map[string]*schema.Resource{
			"foreman_dns":    resourceDNS(),
			"foreman_server": resourceServer(),
		},

		//ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"url": "The region where AWS operations will take place. Examples\n" +
			"are us-east-1, us-west-2, etc.",

		"username": "The access key for API operations. You can retrieve this\n" +
			"from the 'Security & Credentials' section of the AWS console.",

		"password": "The access key for API operations. You can retrieve this\n" +
			"from the 'Security & Credentials' section of the AWS console.",
	}
}

/*func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		URL:      d.Get("url").(string),
	}

	return config.Client()
}*/
