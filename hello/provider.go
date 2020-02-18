package hello

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		/*Schema: map[string]*schema.Schema {
		    "nom": &schema.Schema {
		        Type: schema.TypeString,
		        Description: "Your name",
		        Required: true,
		    },
		},*/
		ResourcesMap: map[string]*schema.Resource{
			"hello_world": resourceServer(),
		},
		ConfigureFunc: providerConfigure,
	}

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{}

	return config.Client()
}