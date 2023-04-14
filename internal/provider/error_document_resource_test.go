package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccErrorDocumentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccErrorDocumentResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_error_document.test", "metadata.name", "tf-acc-test"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "ubika_error_document.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.namespace", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccErrorDocumentResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_error_document.test", "metadata.namespace", "tf-acc-tests"),
				),
			},
			// // Delete testing automatically occurs in TestCase
		},
	})
}

func testAccErrorDocumentResourceConfig(name string, namespace string) string {
	return fmt.Sprintf(`
resource "ubika_error_document" "test" {
  metadata = {
    name = %[1]q
    namespace = %[2]q
  }
  spec = {
	page = <<EOT
<html>
  <head>
	<title>ErrorDocument SSI!</title>
  </head>
  <body>
	<h1>Formid'!</h1>

	Unique ID: {{BLOCK_ID}}<br/>
	Date: {{BLOCK_TIME}}<br/>
	<br/>
  </body>
</html>
EOT
	content_type = "text/html"
  }
}
`, name, namespace)
}
