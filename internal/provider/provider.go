package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sshidProvider struct {
	version string
}

type sshidProviderModel struct {
	HostKeyAlgorithms types.List `tfsdk:"host_key_algorithms"`
}

type sshidProviderState struct {
	HostKeyAlgorithms []string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sshidProvider{
			version: version,
		}
	}
}

func (p *sshidProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	state := &sshidProviderState{}

	var data sshidProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(
		data.HostKeyAlgorithms.ElementsAs(ctx, &state.HostKeyAlgorithms, false)...,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = state
}

func (p *sshidProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *sshidProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

func (p *sshidProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sshid"
	resp.Version = p.version
}

func (p *sshidProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newHostResource,
	}
}

func (p *sshidProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host_key_algorithms": schema.ListAttribute{
				ElementType: types.StringType,
				MarkdownDescription: "An ordered list of public key type names (of the kind found in the second field of an entry in your `~/.ssh/authorized_keys` file) to request from remote SSH servers. If this is not specified then the default sequence of algorithms built in to the Go `crypto/ssh` library will be used.\n\n" +
					"  The first key type that the server supports will be used for SSH host key checks and any other host key types will be ignored. This is less secure than OpenSSH, which checks all of a remote host's known keys, but this deficiency is due to what appears to be a limitation in the API of Go's `crypto/ssh` module.\n\n" +
					"  The `crypto/ssh` module seems to start negotiations by requesting host keys based on NIST elliptic curves by default, so you might want to specify `[\"ssh-ed25519\"]` here to force the use of the less-dubious Ed25519 algorithm instead.",
				Optional: true,
			},
		},
	}
}
