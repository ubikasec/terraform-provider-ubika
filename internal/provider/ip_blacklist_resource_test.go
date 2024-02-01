package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPBlacklitResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIPBlacklistResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_ip_blacklist.test", "metadata.namespace", "tf-acc-tests"),
					resource.TestCheckResourceAttr("ubika_ip_blacklist.test", "metadata.name", "tf-acc-test"),
					resource.TestCheckResourceAttr("ubika_ip_blacklist.test", "spec.ip_addresses.0", "192.0.2.42"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ubika_ip_blacklist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIPBlacklistResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_ip_blacklist.test", "metadata.namespace", "tf-acc-tests"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIPBlacklistResourceConfig(name string, namespace string) string {
	return fmt.Sprintf(`
resource "ubika_ip_blacklist" "test" {
  metadata = {
    name = %[1]q
    namespace = %[2]q
  }
  spec = {
	ip_addresses = ["192.0.2.42", "198.51.100.11"]
  }
}
`, name, namespace)
}
