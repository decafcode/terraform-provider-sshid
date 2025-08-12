package provider

import (
	"github.com/decafcode/terraform-provider-sshid/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sshid": providerserver.NewProtocol6WithError(provider.New("test")()),
}
