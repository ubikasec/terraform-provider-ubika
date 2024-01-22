package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	assetsv1 "github.com/ubikasec/terraform-provider-ubika/internal/apis/assets.ubika.io/v1beta"
	"github.com/ubikasec/terraform-provider-ubika/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Ensure UbikaProvider satisfies various provider interfaces.
var _ provider.Provider = &UbikaProvider{}

// UbikaProvider defines the provider implementation.
type UbikaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// UbikaProviderModel describes the provider data model.
type UbikaProviderModel struct {
	Host          types.String `tfsdk:"host"`
	Port          types.String `tfsdk:"port"`
	InsecureNoTLS types.Bool   `tfsdk:"insecure_no_tls"`
}

func (p *UbikaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ubika"
	resp.Version = p.version
}

func (p *UbikaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "API Host",
				Optional:            true,
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "API Port",
				Optional:            true,
			},
			"insecure_no_tls": schema.BoolAttribute{
				MarkdownDescription: "disable TLS ",
				Optional:            true,
			},
		},
	}
}

func (p *UbikaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data UbikaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Host.IsNull() || data.Host.Equal(types.StringValue("")) {
		data.Host = types.StringValue("api.ubika.io")
	}

	if data.Port.IsNull() || data.Port.Equal(types.StringValue("")) {
		data.Port = types.StringValue("443")
	}

	endpoint := net.JoinHostPort(data.Host.ValueString(), data.Port.ValueString())

	if data.InsecureNoTLS.IsNull() {
		data.InsecureNoTLS = types.BoolValue(false)
	}

	token, _, err := auth.GetToken(http.DefaultClient, ".appsecctl")
	if err != nil {
		resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to find authentication token, got error: %s", err))
		return
	}

	var transportCredentials credentials.TransportCredentials
	if data.InsecureNoTLS.ValueBool() {
		transportCredentials = insecure.NewCredentials()

	} else {
		transportCredentials = credentials.NewTLS(&tls.Config{})
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("dns:///%s", endpoint),
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithPerRPCCredentials(auth.NewPerRPCCredentials("bearer", token)),
	)
	if err != nil {
		resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to connect, got error: %s", err))
		return
	}

	client := assetsv1.NewGRPCClient(conn)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *UbikaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAssetResource,
		NewErrorDocumentResource,
		NewOpenAPIResource,
	}
}

func (p *UbikaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UbikaProvider{
			version: version,
		}
	}
}
