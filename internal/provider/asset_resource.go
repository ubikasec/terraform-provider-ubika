package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	assetsv1 "github.com/ubikasec/terraform-provider-ubika/internal/apis/assets.ubika.io/v1beta"
	metav1 "github.com/ubikasec/terraform-provider-ubika/internal/apis/meta/v1beta"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AssetResource{}
var _ resource.ResourceWithImportState = &AssetResource{}

func NewAssetResource() resource.Resource {
	return &AssetResource{}
}

// AssetResource defines the resource implementation.
type AssetResource struct {
	client assetsv1.Client
}

func (r *AssetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset"
}

func (r *AssetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Asset resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"spec": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"application_module": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"security_mode": schema.StringAttribute{
								Optional: true,
							},
							"exception_profile": schema.StringAttribute{
								MarkdownDescription: "Exception profile (deprecated)",
								Optional:            true,
							},
						},
					},
					"api_module": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"security_mode": schema.StringAttribute{
								Optional: true,
							},
							"openapi": schema.StringAttribute{
								MarkdownDescription: "OpenAPI resource name",
								Required:            true,
							},
						},
					},
					"web_socket_module": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"security_mode": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					"geo_ip_module": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"security_mode": schema.StringAttribute{
								Optional: true,
							},
							"countries": schema.SetAttribute{
								Required:    true,
								ElementType: types.StringType,
							},
							"mode": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					"ip_reputation_module": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"security_mode": schema.StringAttribute{
								Optional: true,
							},
							"threats": schema.SetAttribute{
								Required:    true,
								ElementType: types.StringType,
							},
						},
					},
					"ip_blacklist_module": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"ip_blacklist": schema.StringAttribute{
								Required: true,
							},
							"security_mode": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					"custom_wkf_module": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"workflow": schema.StringAttribute{
								Optional: true,
							},
							"workflow_params": schema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
					"hostnames": schema.SetAttribute{
						Required:    true,
						ElementType: types.StringType,
					},
					"backend_url": schema.StringAttribute{
						MarkdownDescription: "Backend URL",
						Required:            true,
					},
					"backend_certificate_check": schema.StringAttribute{
						MarkdownDescription: "Check backend certificate",
						Optional:            true,
						Computed:            true,
					},
					"trusted_ip_address_header": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"exception_profile": schema.StringAttribute{
						MarkdownDescription: "Exception profile name",
						Optional:            true,
						Computed:            true,
					},
					"deployment_type": schema.StringAttribute{
						MarkdownDescription: "Deployment type (SAAS or SELF_HOSTED)",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"tls_mode": schema.StringAttribute{
						MarkdownDescription: "TLS mode (auto or custom)",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("NONE"),
					},
					"tls_material": schema.StringAttribute{
						MarkdownDescription: "TLS Material name",
						Optional:            true,
						Computed:            true,
					},
					"tls_configuration": schema.StringAttribute{
						MarkdownDescription: "TLS Configuration name",
						Optional:            true,
						Computed:            true,
					},
					"blocking_page": schema.StringAttribute{
						MarkdownDescription: "Blocking page name",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"unavailable_page": schema.StringAttribute{
						MarkdownDescription: "Unavailable page name",
						Optional:            true,
						Computed:            true,
					},
					"maintenance_page": schema.StringAttribute{
						MarkdownDescription: "Maintenance page name",
						Optional:            true,
						Computed:            true,
					},
					"maintenance_enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable maintenance page name",
						Optional:            true,
						Computed:            true,
					},
				},
			},
			"metadata": GetObjectMetaResource(),
			"status": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"state": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"redirected_hostnames": schema.SetAttribute{
								Computed:    true,
								ElementType: types.StringType,
							},
							"backend": schema.StringAttribute{
								Computed: true,
							},
							"dns": schema.StringAttribute{
								Computed: true,
							},
							"runningstate": schema.StringAttribute{
								Computed: true,
							},
						},
					},
					"tls": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"begins_on": schema.StringAttribute{
								Computed: true,
							},
							"expires_on": schema.StringAttribute{
								Computed: true,
							},
							"mode": schema.StringAttribute{
								Computed: true,
							},
						},
					},
					"service_address": schema.StringAttribute{
						MarkdownDescription: "Address of the service",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *AssetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = getClient(req, resp)
}

func (r *AssetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Asset")

	// Read Terraform plan data into the model
	var plan *assetsv1.AssetResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	asset, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create the resource
	asset, err := r.client.Asset().Create(ctx, asset)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create asset, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.AssetResourceModel
	_, err = state.FromProto(asset)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from asset, got error: %s", err))
		return
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 200 * time.Millisecond
	bo.MaxInterval = 10 * time.Second
	bo.MaxElapsedTime = 0

	retryFunc := func() error {
		asset, err := r.client.Asset().Get(ctx, &metav1.GetOptions{
			Name:      asset.GetMetadata().GetName(),
			Namespace: asset.GetMetadata().GetNamespace(),
		})
		if err != nil {
			return err
		}

		_, err = state.FromProto(asset)
		if err != nil {
			return err
		}

		if asset.GetStatus().GetServiceAddress() == "" {
			return errors.New("retry")
		}
		return nil
	}

	// Retry until context is cancelled
	err = backoff.Retry(retryFunc, bo)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get status from asset, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an asset")

	// Save state data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Asset")

	// Read Terraform prior state data into the model
	var state *assetsv1.AssetResourceTFModel
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

func (r *AssetResource) read(ctx context.Context, metaObjValue basetypes.ObjectValue, meta *metav1.ObjectMetaResourceTFModel) (assetsv1.AssetResourceModel, diag.Diagnostics) {
	// get metadata from state
	if meta == nil {
		if diags := metaObjValue.As(ctx, &meta, basetypes.ObjectAsOptions{}); diags.HasError() {
			return assetsv1.AssetResourceModel{}, diags
		}
	}

	asset, err := r.client.Asset().Get(ctx, &metav1.GetOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		return assetsv1.AssetResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read asset %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}

	// update state from protobuf resource
	var state assetsv1.AssetResourceModel
	_, err = state.FromProto(asset)
	if err != nil {
		return assetsv1.AssetResourceModel{}, []diag.Diagnostic{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to get state from asset %s/%s, got error: %s", meta.Name.ValueString(), meta.Namespace.ValueString(), err))}
	}

	return state, nil
}

func (r *AssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan *assetsv1.AssetResourceTFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to protobuf resource
	asset, diags := plan.ToProto(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	asset, err := r.client.Asset().Update(ctx, asset)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update asset, got error: %s", err))
		return
	}

	// generate state from protobuf resource
	var state assetsv1.AssetResourceModel
	_, err = state.FromProto(asset)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get state from asset, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Asset")
	var plan *assetsv1.AssetResourceTFModel

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
	_, err := r.client.Asset().Delete(ctx, &metav1.DeleteOptions{
		Name:      meta.Name.ValueString(),
		Namespace: meta.Namespace.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete asset, got error: %s", err))
		return
	}
}

func (r *AssetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
