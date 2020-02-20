package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"hello/hello"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return hello.Provider()
		},
	})
	log.Println("[TEST] test d'affichage")
}
