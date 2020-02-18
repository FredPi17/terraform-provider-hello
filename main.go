package main

import (
	"github.com/FredPi17/terraform-provider-hello/hello"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return hello.Provider()
		},
	})
}

func Provider() *schema.Provider {
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

	return nil, nil
}

func resourceServer() *schema.Resource {

	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"nom": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	nom := d.Get("nom").(string)
	d.SetId(nom)
	return resourceServerCreate(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
