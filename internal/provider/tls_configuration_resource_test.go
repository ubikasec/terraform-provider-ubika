package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTLSConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTLSConfigurationResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_tls_configuration.test", "metadata.name", "tf-acc-test"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "ubika_tls_configuration.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.namespace", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccTLSConfigurationResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_tls_configuration.test", "metadata.namespace", "tf-acc-tests"),
				),
			},
			// // Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTLSConfigurationResourceConfig(name string, namespace string) string {
	return fmt.Sprintf(`
resource "ubika_tls_configuration" "test" {
  metadata = {
    name = %[1]q
    namespace = %[2]q
  }
  spec = {
	protocol_min= "TLS_1_0"
	protocol_max= "TLS_1_2"
	ciphers = ["ECDHE-ECDSA-AES128-SHA"]
  }
}
`, name, namespace)
}
