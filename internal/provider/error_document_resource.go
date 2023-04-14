package provider

import (
	"context"
	"fmt"
	"strings"

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
var _ resource.Resource = &ErrorDocumentResource{}
var _ resource.ResourceWithImportState = &ErrorDocumentResource{}

func NewErrorDocumentResource() resource.Resource {
	return &ErrorDocumentResource{}
}

// ErrorDocumentResource defines the resource implementation.
type ErrorDocumentResource struct {
	client assetsv1.Client
}

func (r *ErrorDocumentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_error_document"
}

func (r *ErrorDocumentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ErrorDocument resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"metadata": GetObjectMetaResource(),
			"spec": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"page": schema.StringAttribute{
						MarkdownDescription: "Page",
						Required:            true,
					},
					"content_type": schema.StringAttribute{
						MarkdownDescription: "Content type",
						Required:            true,
					},
				},
			},
		},
	}
}

func (r *ErrorDocumentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(assetsv1.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *assetsv1.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ErrorDocumentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating ErrorDocument")

	// Read Terraform plan data into the model
	var plan *assetsv1.ErrorDocumentResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	errorDocument, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create the resource
	errorDocument, err := r.client.ErrorDocument().Create(ctx, errorDocument)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create error document, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.ErrorDocumentResourceModel
	_, err = state.FromProto(errorDocument)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from error document, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an error document")

	// Save state data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ErrorDocumentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading ErrorDocument")

	// Read Terraform prior state data into the model
	var state *assetsv1.ErrorDocumentResourceTFModel
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

func (r *ErrorDocumentResource) read(ctx context.Context, metaObjValue basetypes.ObjectValue, meta *metav1.ObjectMetaResourceTFModel) (assetsv1.ErrorDocumentResourceModel, diag.Diagnostics) {
	// get metadata from state
	if meta == nil {
		if diags := metaObjValue.As(ctx, &meta, basetypes.ObjectAsOptions{}); diags.HasError() {
			return assetsv1.ErrorDocumentResourceModel{}, diags
		}
	}

	errorDocument, err := r.client.ErrorDocument().Get(ctx, &metav1.GetOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		return assetsv1.ErrorDocumentResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read error document %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}

	// update state from protobuf resource
	var state assetsv1.ErrorDocumentResourceModel
	_, err = state.FromProto(errorDocument)
	if err != nil {
		return assetsv1.ErrorDocumentResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to get state from error document %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}
	return state, nil
}

func (r *ErrorDocumentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan *assetsv1.ErrorDocumentResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	errorDocument, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	errorDocument, err := r.client.ErrorDocument().Update(ctx, errorDocument)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update error document, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.ErrorDocumentResourceModel
	_, err = state.FromProto(errorDocument)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from error document, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ErrorDocumentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting ErrorDocument")
	var plan *assetsv1.ErrorDocumentResourceTFModel

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
	_, err := r.client.ErrorDocument().Delete(ctx, &metav1.DeleteOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete error document, got error: %s", err))
		return
	}
}

func (r *ErrorDocumentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	var name, namespace string
	if len(parts) == 2 {
		namespace = parts[0]
		name = parts[1]
	} else {
		resp.Diagnostics.AddError("Inexpected input", "A namespace is required, ID must be in the form 'namespace/resource-name'")
	}

	meta := metav1.ObjectMetaResourceTFModel{
		Name:      types.StringValue(name),
		Namespace: types.StringValue(namespace),
	}
	state, diags := r.read(ctx, basetypes.ObjectValue{}, &meta)
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
