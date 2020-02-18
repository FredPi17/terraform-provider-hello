package hello

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() *terraform.ResourceProvider {
	return &schema.Provider{
		/*Schema: map[string]*schema.Schema {
		    "nom": &schema.Schema {
		        Type: schema.TypeString,
		        Description: "Your name",
		        Required: true,
		    },
		},*/
		DataSourcesMap: map[string]*schema.Resource{
			"hello_world": resourceServer(),
		},
		ConfigureFunc: providerConfigure,
	}

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{}

	return config.Client()
}
