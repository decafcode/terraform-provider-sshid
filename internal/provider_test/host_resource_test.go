package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"gotest.tools/assert"
)

func TestAccExampleResource(t *testing.T) {
	f, err := newFramework()
	assert.NilError(t, err)
	defer f.Close()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "sshid_host" "test" {
						host = "%s"
						port = %d
					}
				`, f.Host(), f.Port()),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"sshid_host.test",
						tfjsonpath.New("public_key"),
						knownvalue.StringExact(f.PublicKey()),
					),
				},
			},
		},
	})
}
