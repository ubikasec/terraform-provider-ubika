package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	assetsv1 "github.com/ubikasec/terraform-provider-ubika/internal/client/assets.ubika.io/v1beta"
)

func getClient(req resource.ConfigureRequest, resp *resource.ConfigureResponse) assetsv1.Client {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	client, ok := req.ProviderData.(assetsv1.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *assetsv1.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil
	}

	return client
}
