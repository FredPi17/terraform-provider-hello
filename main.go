package main

import (
	"github.com/FredPi17/terraform-provider-hello/hello"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hello.Provider})
}
