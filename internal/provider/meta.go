package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	metav1 "github.com/ubikasec/terraform-provider-ubika/internal/client/meta/v1beta"
)

func GetObjectMetaResource() schema.Attribute {
	return schema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the resource",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "Namespace of the resource",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created": schema.Int64Attribute{
				Computed: true,
			},
			"updated": schema.Int64Attribute{
				Computed: true,
			},
			"version": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func getMeta(req resource.ImportStateRequest, resp *resource.ImportStateResponse) *metav1.ObjectMetaResourceTFModel {
	parts := strings.Split(req.ID, "/")
	var name, namespace string
	if len(parts) == 2 {
		namespace = parts[0]
		name = parts[1]
	} else {
		resp.Diagnostics.AddError("Inexpected input", "A namespace is required, ID must be in the form 'namespace/resource-name'")
	}

	return &metav1.ObjectMetaResourceTFModel{
		Name:      types.StringValue(name),
		Namespace: types.StringValue(namespace),
	}
}
