package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExceptionProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExceptionProfileResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_exception_profile.test", "metadata.name", "tf-acc-test"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "ubika_exception_profile.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.namespace", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccExceptionProfileResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_exception_profile.test", "metadata.namespace", "tf-acc-tests"),
				),
			},
			// // Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExceptionProfileResourceConfig(name string, namespace string) string {
	return fmt.Sprintf(`
resource "ubika_exception_profile" "test" {
  metadata = {
    name = %[1]q
    namespace = %[2]q
  }
  spec = {
    rules = [
      {
        name = "my rule"
        filters = [
          "url == \"/url_1\"",
        ]
      },
      {
        name = "my second rule"
        filters = [
          "url == \"/url_2\"",
        ]
      }
    ]
  }
}
`, name, namespace)
}
