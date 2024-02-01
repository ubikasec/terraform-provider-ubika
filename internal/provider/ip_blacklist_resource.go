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
var _ resource.Resource = &ipBlacklistResource{}
var _ resource.ResourceWithImportState = &ipBlacklistResource{}

func NewIPBlacklistResource() resource.Resource {
	return &ipBlacklistResource{}
}

// ipBlacklistResource defines the resource implementation.
type ipBlacklistResource struct {
	client assetsv1.Client
}

func (r *ipBlacklistResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_blacklist"
}

func (r *ipBlacklistResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "IP blacklist resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"metadata": GetObjectMetaResource(),
			"spec": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"ip_addresses": schema.SetAttribute{
						Required:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

func (r *ipBlacklistResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = getClient(req, resp)
}

func (r *ipBlacklistResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating IP blacklist")

	// Read Terraform plan data into the model
	var plan *assetsv1.IPBlacklistResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	ipBlacklist, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create the resource
	ipBlacklist, err := r.client.IPBlacklist().Create(ctx, ipBlacklist)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create IP blacklist, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.IPBlacklistResourceModel
	_, err = state.FromProto(ipBlacklist)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from IP blacklist, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an IP blacklist")

	// Save state data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ipBlacklistResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading IP blacklist")

	// Read Terraform prior state data into the model
	var state *assetsv1.IPBlacklistResourceTFModel
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

func (r *ipBlacklistResource) read(ctx context.Context, metaObjValue basetypes.ObjectValue, meta *metav1.ObjectMetaResourceTFModel) (assetsv1.IPBlacklistResourceModel, diag.Diagnostics) {
	// get metadata from state
	if meta == nil {
		if diags := metaObjValue.As(ctx, &meta, basetypes.ObjectAsOptions{}); diags.HasError() {
			return assetsv1.IPBlacklistResourceModel{}, diags
		}
	}

	ipBlacklist, err := r.client.IPBlacklist().Get(ctx, &metav1.GetOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		return assetsv1.IPBlacklistResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read IP blacklist %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}

	// update state from protobuf resource
	var state assetsv1.IPBlacklistResourceModel
	_, err = state.FromProto(ipBlacklist)
	if err != nil {
		return assetsv1.IPBlacklistResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to get state from IP blacklist %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}
	return state, nil
}

func (r *ipBlacklistResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan *assetsv1.IPBlacklistResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	ipBlacklist, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipBlacklist, err := r.client.IPBlacklist().Update(ctx, ipBlacklist)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update IP blacklist, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.IPBlacklistResourceModel
	_, err = state.FromProto(ipBlacklist)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from IP blacklist, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ipBlacklistResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting IP blacklist")
	var plan *assetsv1.IPBlacklistResourceTFModel

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
	_, err := r.client.IPBlacklist().Delete(ctx, &metav1.DeleteOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete IP blacklist, got error: %s", err))
		return
	}
}

func (r *ipBlacklistResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
