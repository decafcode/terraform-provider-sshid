package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

// Most of the tests for this Terraform are implemented as unit tests for the
// underlying Go function that implements its functionality, this just tests
// that the aforementioned Go function is integrated into the Terraform SDK
// correctly.

const x509Input = `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3m3JQFAEc2DqBGhA/3K68q5HwjSp
KASwcbiLtht5ne7CmyVRUT7qKVCphcmm81Hy6bzUR6PZLaMhToq8dTnXAA==
-----END PUBLIC KEY-----
`

const expectedOutput = "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBN5tyUBQBHNg6gRoQP9yuvKuR8I0qSgEsHG4i7YbeZ3uwpslUVE+6ilQqYXJpvNR8um81Eej2S2jIU6KvHU51wA="

func TestAccFromSpkiFunction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					output "test" {
						value = provider::sshid::fromspki("%s")
					}
				`, strings.ReplaceAll(x509Input, "\n", "\\n")),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact(expectedOutput),
					),
				},
			},
			{
				Config: `
					output "test" {
						value = provider::sshid::fromspki("invalid input")
					}
				`,
				ExpectError: regexp.MustCompile("expected input in PEM"),
			},
		},
	})
}
