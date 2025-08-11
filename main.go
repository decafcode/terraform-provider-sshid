package main

import (
	"context"
	"flag"
	"log"

	"github.com/decafcode/terraform-provider-sshid/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "enable debugging")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/decafcode/sshid",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
