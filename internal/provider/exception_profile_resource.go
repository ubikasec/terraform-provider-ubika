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
var _ resource.Resource = &ExceptionProfileResource{}
var _ resource.ResourceWithImportState = &ExceptionProfileResource{}

func NewExceptionProfileResource() resource.Resource {
	return &ExceptionProfileResource{}
}

// ExceptionProfileResource defines the resource implementation.
type ExceptionProfileResource struct {
	client assetsv1.Client
}

func (r *ExceptionProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception_profile"
}

func (r *ExceptionProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ExceptionProfile resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"metadata": GetObjectMetaResource(),
			"spec": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"rules": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"filters": schema.ListAttribute{
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ExceptionProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = getClient(req, resp)
}

func (r *ExceptionProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating ExceptionProfile")

	// Read Terraform plan data into the model
	var plan *assetsv1.ExceptionProfileResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	exceptionProfile, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create the resource
	exceptionProfile, err := r.client.ExceptionProfile().Create(ctx, exceptionProfile)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create exception profile, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.ExceptionProfileResourceModel
	_, err = state.FromProto(exceptionProfile)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from exception profile, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an exception profile")

	// Save state data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ExceptionProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading ExceptionProfile")

	// Read Terraform prior state data into the model
	var state *assetsv1.ExceptionProfileResourceTFModel
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

func (r *ExceptionProfileResource) read(ctx context.Context, metaObjValue basetypes.ObjectValue, meta *metav1.ObjectMetaResourceTFModel) (assetsv1.ExceptionProfileResourceModel, diag.Diagnostics) {
	// get metadata from state
	if meta == nil {
		if diags := metaObjValue.As(ctx, &meta, basetypes.ObjectAsOptions{}); diags.HasError() {
			return assetsv1.ExceptionProfileResourceModel{}, diags
		}
	}

	exceptionProfile, err := r.client.ExceptionProfile().Get(ctx, &metav1.GetOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		return assetsv1.ExceptionProfileResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read exception profile %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}

	// update state from protobuf resource
	var state assetsv1.ExceptionProfileResourceModel
	_, err = state.FromProto(exceptionProfile)
	if err != nil {
		return assetsv1.ExceptionProfileResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to get state from exception profile %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}
	return state, nil
}

func (r *ExceptionProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan *assetsv1.ExceptionProfileResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	exceptionProfile, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	exceptionProfile, err := r.client.ExceptionProfile().Update(ctx, exceptionProfile)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update exception profile, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.ExceptionProfileResourceModel
	_, err = state.FromProto(exceptionProfile)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from exception profile, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ExceptionProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting ExceptionProfile")
	var plan *assetsv1.ExceptionProfileResourceTFModel

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
	_, err := r.client.ExceptionProfile().Delete(ctx, &metav1.DeleteOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete exception profile, got error: %s", err))
		return
	}
}

func (r *ExceptionProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
