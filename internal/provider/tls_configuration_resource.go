package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	assetsv1 "github.com/ubikasec/terraform-provider-ubika/internal/apis/assets.ubika.io/v1beta"
	metav1 "github.com/ubikasec/terraform-provider-ubika/internal/apis/meta/v1beta"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TLSConfigurationResource{}
var _ resource.ResourceWithImportState = &TLSConfigurationResource{}

func NewTLSConfigurationResource() resource.Resource {
	return &TLSConfigurationResource{}
}

// TLSConfigurationResource defines the resource implementation.
type TLSConfigurationResource struct {
	client assetsv1.Client
}

func (r *TLSConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tls_configuration"
}

func (r *TLSConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "TLSConfiguration resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"metadata": GetObjectMetaResource(),
			"spec": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"protocol_min": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"protocol_max": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"ciphers": schema.SetAttribute{
						MarkdownDescription: "Ciphers for TLS 1.0 to 1.2",
						Required:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}

func (r *TLSConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = getClient(req, resp)
}

func (r *TLSConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating TLSConfiguration")

	// Read Terraform plan data into the model
	var plan *assetsv1.TLSConfigurationResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	tlsConfiguration, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create the resource
	tlsConfiguration, err := r.client.TLSConfiguration().Create(ctx, tlsConfiguration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tls configuration, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.TLSConfigurationResourceModel
	_, err = state.FromProto(tlsConfiguration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from tls configuration, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an tls configuration")

	// Save state data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TLSConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading TLSConfiguration")

	// Read Terraform prior state data into the model
	var state *assetsv1.TLSConfigurationResourceTFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState, diags := r.read(ctx, state.Metadata, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *TLSConfigurationResource) read(ctx context.Context, metaObjValue basetypes.ObjectValue, meta *metav1.ObjectMetaResourceTFModel) (assetsv1.TLSConfigurationResourceModel, diag.Diagnostics) {
	// get metadata from state
	if meta == nil {
		if diags := metaObjValue.As(ctx, &meta, basetypes.ObjectAsOptions{}); diags.HasError() {
			return assetsv1.TLSConfigurationResourceModel{}, diags
		}
	}

	tlsConfiguration, err := r.client.TLSConfiguration().Get(ctx, &metav1.GetOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		return assetsv1.TLSConfigurationResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read tls configuration %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}

	// update state from protobuf resource
	var state assetsv1.TLSConfigurationResourceModel
	_, err = state.FromProto(tlsConfiguration)
	if err != nil {
		return assetsv1.TLSConfigurationResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to get state from tls configuration %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}
	return state, nil
}

func (r *TLSConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan *assetsv1.TLSConfigurationResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	tlsConfiguration, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tlsConfiguration, err := r.client.TLSConfiguration().Update(ctx, tlsConfiguration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update tls configuration, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.TLSConfigurationResourceModel
	_, err = state.FromProto(tlsConfiguration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from tls configuration, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TLSConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting TLSConfiguration")
	var plan *assetsv1.TLSConfigurationResourceTFModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var meta metav1.ObjectMetaResourceTFModel
	if diags := plan.Metadata.As(ctx, &meta, basetypes.ObjectAsOptions{}); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	_, err := r.client.TLSConfiguration().Delete(ctx, &metav1.DeleteOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tls configuration, got error: %s", err))
		return
	}
}

func (r *TLSConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state, diags := r.read(ctx, basetypes.ObjectValue{}, getMeta(req, resp))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
