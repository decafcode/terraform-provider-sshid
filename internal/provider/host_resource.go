package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/crypto/ssh"
)

func newHostResource() resource.Resource {
	return &hostResource{}
}

type hostResource struct {
	ps *sshidProviderState
}

type hostResourceModel struct {
	Host      types.String `tfsdk:"host"`
	Port      types.Int32  `tfsdk:"port"`
	PublicKey types.String `tfsdk:"public_key"`
}

func (r *hostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *hostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pseudo-resource that connects to an SSH server and requests the server's public key when the resource is first created. This resource exists entirely within Terraform state. See provider documentation for further details.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Target SSH server's hostname",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"port": schema.Int32Attribute{
				Computed:            true,
				Default:             int32default.StaticInt32(22),
				MarkdownDescription: "Target SSH server's TCP port. Defaults to port 22.",
				Optional:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The public key that was returned by the SSH server when this resource was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *hostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ps, ok := req.ProviderData.(*sshidProviderState)

	if !ok {
		resp.Diagnostics.AddError("Internal error", "Invalid provider state type")

		return
	}

	r.ps = ps
}

func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data hostResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.captureHostKey(ctx, &data)

	if err != nil {
		resp.Diagnostics.AddError("Connection failed", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *hostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// no-op
}

func (r *hostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Invalid state", "Resource is immutable")
}

func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// no-op
}

func (r *hostResource) captureHostKey(ctx context.Context, data *hostResourceModel) error {
	if !data.PublicKey.IsUnknown() {
		return fmt.Errorf("resource is already configured, this shouldn't happen?")
	}

	var hostKey string

	config := ssh.ClientConfig{
		HostKeyAlgorithms: r.ps.HostKeyAlgorithms,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			hostKey = key.Type() + " " + base64.StdEncoding.EncodeToString(key.Marshal())

			return fmt.Errorf("dummy error")
		},
	}

	d := net.Dialer{}
	addr := net.JoinHostPort(data.Host.ValueString(), fmt.Sprintf("%d", data.Port.ValueInt32()))
	tcpConn, tcpErr := d.DialContext(ctx, "tcp", addr)

	if tcpErr != nil {
		return tcpErr
	}

	_, _, _, sshErr := ssh.NewClientConn(tcpConn, addr, &config)
	_ = tcpConn.Close()

	if hostKey == "" {
		// We didn't reach the host key callback, something else went wrong.

		if sshErr == nil {
			return fmt.Errorf("unknown error")
		} else {
			return sshErr
		}
	}

	data.PublicKey = types.StringValue(hostKey)

	tflog.Warn(ctx, "Captured a new SSH host key", map[string]any{
		"host":       data.Host.ValueString(),
		"port":       data.Port.ValueInt32(),
		"public_key": data.PublicKey.ValueString(),
	})

	return nil
}
