package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
