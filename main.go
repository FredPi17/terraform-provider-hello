package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/rs/xid"
	"log"
	"os"
	_ "reflect"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"hello_world":  resourceServer(),
			"hello_create": resourceCreate(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var client Client
	return &client, nil
}

//Config is the config for the client.
type Config struct {
}

//Client is the client itself. Since we already have access to the shell no real provisioning needs to be done
type Client struct {
}

func resourceCreate() *schema.Resource {
	return &schema.Resource{
		Create: resourceHelloCreate,

		Schema: map[string]*schema.Schema{
			"create": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"working_directory": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  ".",
			},
			"output": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func resourceHelloCreate(d *schema.ResourceData, m interface{}) error {
	os.Stdout.WriteString("creating resource")
	return createHello(d, m, []string{"create"})
}

func createHello(d *schema.ResourceData, m interface{}, stack []string) error {
	os.Stdout.WriteString("Creating shell script resource")
	log.Printf("[DEBUG] Creating shell script resource...")
	printStackTrace(stack)
	l := d.Get("create").([]interface{})
	c := l[0].(map[string]interface{})
	command := c["create"].(string)
	environment := readEnvironmentVariables(vars)
	workingDirectory := d.Get("working_directory").(string)
	d.MarkNewResource()
	//obtain exclusive lock
	//shellMutexKV.Lock(shellScriptMutexKey)

	output := make(map[string]string)
	state := NewState(environment, output)
	newState, err := runCommand(command, state, environment, workingDirectory)
	if err != nil {
		os.Stderr.WriteString(err)
		return err
	}
	//shellMutexKV.Unlock(shellScriptMutexKey)

	//if create doesn't return a new state then must call the read operation
	if newState == nil {
		stack = append(stack, "read")
		if err := read(d, meta, stack); err != nil {
			return err
		}
	} else {
		d.Set("output", newState.Output)
	}

	//create random uuid for the id
	id := xid.New().String()
	d.SetId(id)
	return nil
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"nom": {
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
	d.SetId("")
	return nil
}
