package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAssetResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_asset.test", "metadata.name", "tf-acc-test"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "ubika_asset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.namespace", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccAssetResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_asset.test", "metadata.namespace", "tf-acc-tests"),
				),
			},
			// // Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAssetResourceConfig(name string, namespace string) string {
	return fmt.Sprintf(`
resource "ubika_asset" "test" {
  metadata = {
    name = %[1]q
    namespace = %[2]q
  }
  spec = {
    hostnames = ["tf-acc-test.example.com"]
	backend_url = "https://tf-acc-test.example.com/"
	deployment_type = "SAAS"
  }
}
`, name, namespace)
}
