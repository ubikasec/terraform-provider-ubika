package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	assetsv1 "github.com/ubikasec/terraform-provider-ubika/internal/client/assets.ubika.io/v1beta"
	metav1 "github.com/ubikasec/terraform-provider-ubika/internal/client/meta/v1beta"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TLSManualCreateResource{}
var _ resource.ResourceWithImportState = &TLSManualCreateResource{}

func NewTLSMaterialResource() resource.Resource {
	return &TLSManualCreateResource{}
}

// TLSMaterialResource defines the resource implementation.
type TLSManualCreateResource struct {
	client assetsv1.Client
}

func (r *TLSManualCreateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tls_material"
}

func (r *TLSManualCreateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "TLSMaterial resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"metadata": GetObjectMetaResource(),
			"spec": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"certificate": schema.StringAttribute{
						MarkdownDescription: "Certificate",
						Required:            true,
					},
					"chain": schema.StringAttribute{
						MarkdownDescription: "Chain",
						Optional:            true,
						Computed:            true,
					},
					"key": schema.StringAttribute{
						MarkdownDescription: "Key",
						Required:            true,
						Sensitive:           true,
					},
				},
			},
		},
	}
}

func (r *TLSManualCreateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = getClient(req, resp)
}

func (r *TLSManualCreateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating TLSMaterial")

	// Read Terraform plan data into the model
	var plan *assetsv1.TLSManualCreateResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	tlsManualCreate, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create the resource
	tlsMaterial, err := r.client.TLSConfiguration().CreateManualTLS(ctx, tlsManualCreate)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tls material, got error: %s", err))
		return
	}

	newTLSManualCreate := assetsv1.TLSManualCreate{
		ApiVersion: tlsMaterial.GetApiVersion(),
		Kind:       tlsManualCreate.GetKind(),
		Metadata: &metav1.ObjectMeta{
			Name:      tlsMaterial.GetMetadata().GetName(),
			Namespace: tlsMaterial.GetMetadata().GetNamespace(),
		},
		Spec: &assetsv1.TLSManualCreateSpec{
			Certificate: tlsMaterial.GetSpec().GetCertificate(),
			Chain:       tlsMaterial.GetSpec().GetChain(),
			Key:         tlsManualCreate.GetSpec().GetKey(),
		},
	}

	// generate state from protobuf resource
	var state assetsv1.TLSManualCreateResourceModel
	if err := state.FromProto(&newTLSManualCreate); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from tls material, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an tls material")

	// Save state data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TLSManualCreateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading TLSMaterial")

	// Read Terraform prior state data into the model
	var state *assetsv1.TLSManualCreateResourceTFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState, diags := r.read(ctx, state.Spec, state.Metadata, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *TLSManualCreateResource) read(ctx context.Context, specObjValue basetypes.ObjectValue, metaObjValue basetypes.ObjectValue, meta *metav1.ObjectMetaResourceTFModel) (assetsv1.TLSManualCreateResourceModel, diag.Diagnostics) {
	// get metadata from state
	if meta == nil {
		if diags := metaObjValue.As(ctx, &meta, basetypes.ObjectAsOptions{}); diags.HasError() {
			return assetsv1.TLSManualCreateResourceModel{}, diags
		}
	}

	var spec *assetsv1.TLSManualCreateSpecResourceTFModel
	if diags := specObjValue.As(ctx, &spec, basetypes.ObjectAsOptions{}); diags.HasError() {
		return assetsv1.TLSManualCreateResourceModel{}, diags
	}

	var key string
	if spec != nil {
		key = spec.Key.ValueString()
	}

	tlsMaterial, err := r.client.TLSConfiguration().GetTLSMaterial(ctx, &metav1.GetOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		return assetsv1.TLSManualCreateResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read tls material %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}

	newTLSManualCreate := assetsv1.TLSManualCreate{
		ApiVersion: tlsMaterial.GetApiVersion(),
		Kind:       "TLSManualCreate",
		Metadata: &metav1.ObjectMeta{
			Name:      tlsMaterial.GetMetadata().GetName(),
			Namespace: tlsMaterial.GetMetadata().GetNamespace(),
		},
		Spec: &assetsv1.TLSManualCreateSpec{
			Certificate: tlsMaterial.GetSpec().GetCertificate(),
			Chain:       tlsMaterial.GetSpec().GetChain(),
			Key:         key,
		},
	}

	// update state from protobuf resource
	var state assetsv1.TLSManualCreateResourceModel
	if err := state.FromProto(&newTLSManualCreate); err != nil {
		return assetsv1.TLSManualCreateResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to get state from tls material %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}
	return state, nil
}

func (r *TLSManualCreateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan *assetsv1.TLSManualCreateResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	tlsManualCreate, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tlsMaterial, err := r.client.TLSConfiguration().UpdateManualTLS(ctx, tlsManualCreate)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update tls material, got error: %s", err))
		return
	}

	newTLSManualCreate := assetsv1.TLSManualCreate{
		ApiVersion: tlsMaterial.GetApiVersion(),
		Kind:       tlsManualCreate.GetKind(),
		Metadata: &metav1.ObjectMeta{
			Name:      tlsMaterial.GetMetadata().GetName(),
			Namespace: tlsMaterial.GetMetadata().GetNamespace(),
		},
		Spec: &assetsv1.TLSManualCreateSpec{
			Certificate: tlsMaterial.GetSpec().GetCertificate(),
			Chain:       tlsMaterial.GetSpec().GetChain(),
			Key:         tlsManualCreate.GetSpec().GetKey(),
		},
	}

	// generate state from protobuf resource
	var state assetsv1.TLSManualCreateResourceModel
	if err := state.FromProto(&newTLSManualCreate); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from tls material, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TLSManualCreateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting TLSMaterial")
	var plan *assetsv1.TLSManualCreateResourceTFModel

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
	_, err := r.client.TLSConfiguration().DeleteTLSMaterial(ctx, &metav1.DeleteOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tls material, got error: %s", err))
		return
	}
}

func (r *TLSManualCreateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state, diags := r.read(ctx, basetypes.ObjectValue{}, basetypes.ObjectValue{}, getMeta(req, resp))
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
