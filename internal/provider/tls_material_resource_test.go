package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTLSMaterialResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTLSMaterialResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_tls_material.test", "metadata.name", "tf-acc-test"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTLSMaterialResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_tls_material.test", "metadata.namespace", "tf-acc-tests"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTLSMaterialResourceConfig(name string, namespace string) string {
	return fmt.Sprintf(`
resource "ubika_tls_material" "test" {
  metadata = {
    name = %[1]q
    namespace = %[2]q
  }
  spec = {
	certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIDdzCCAl+gAwIBAgIUNBzY3GlH/l5y9A8sYiePL9rxzlYwDQYJKoZIhvcNAQEL
BQAwHDEaMBgGA1UEAwwRc2l0ZTIuZXhhbXBsZS5jb20wHhcNMjMwMjA2MTI0MzMw
WhcNMjQwMjA2MTI0MzMwWjAcMRowGAYDVQQDDBFzaXRlMi5leGFtcGxlLmNvbTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAO89u+CHlwTdIkGonN0SFdR4
IQJO2Xu4NNRhfbxfw4HHgBxfwLYBFFkA4JYsOX3MQvhgR5chdpRXU2aZsD908kvn
qtlktcx1uaQnLPYozA3GoxJED9lmdJ5n84Pe6IGM1mGU3TCQJUQtZMA9OzdyMcQ0
Ig/HufXdKkGdgTFtFWJImKQi4u2O21bJnoSXixJTWmYFV1Zg6DOUN3BJmyzajVs9
tBDudOPWys0QVbD5u5kfkp3cQLs2VwCMTjtr0XlfwshFiD0tdFkG7g/Y4ocpegFr
AQB4ObFnMWzetFO03zrKWNZyiN2OmELoz6WYFONBA/Vn2bYdZtEgQA2qZ3sET1UC
AwEAAaOBsDCBrTAJBgNVHRMEAjAAMAsGA1UdDwQEAwIFoDAdBgNVHQ4EFgQUdoz0
uMTMjc5e/OlEbKsojlWcxGkwQQYDVR0jBDowOKEgpB4wHDEaMBgGA1UEAwwRc2l0
ZTIuZXhhbXBsZS5jb22CFDQc2NxpR/5ecvQPLGInjy/a8c5WMBwGA1UdEQQVMBOC
EXNpdGUzLmV4YW1wbGUuY29tMBMGA1UdJQQMMAoGCCsGAQUFBwMBMA0GCSqGSIb3
DQEBCwUAA4IBAQAc07MNMEaiTV4pvsC3ky9a8Qr8XMfVT6tp3lXm8EV2JmRcCtJX
nmrSBcMH/PegCeIaAfs/O33hPe7sNF/7ImNic6eI6TF7+3oGhit79Ed4kJ+MGEub
+/WjbSz+fCmDJvqVL5HjGzLKs/YdmDV1zaXdEkDFRyN15cXpmMrmiVM4TBilh8Hz
XvCAnOdwFBJPsSQGQ63EbWLUS3mZ8C031pG/lYk18rfqEuf97LBE/24iemaV+sK3
HHjqDhDQLA+xzngd/Wfvtz5Ndt7jU60ZfLmgxUvzNMMwYK9jrbRSUcKYifxWerk6
MiBpJuczzZB15xs2b7/okkKL/Iyb6mzKEnDV
-----END CERTIFICATE-----
EOT
	key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDvPbvgh5cE3SJB
qJzdEhXUeCECTtl7uDTUYX28X8OBx4AcX8C2ARRZAOCWLDl9zEL4YEeXIXaUV1Nm
mbA/dPJL56rZZLXMdbmkJyz2KMwNxqMSRA/ZZnSeZ/OD3uiBjNZhlN0wkCVELWTA
PTs3cjHENCIPx7n13SpBnYExbRViSJikIuLtjttWyZ6El4sSU1pmBVdWYOgzlDdw
SZss2o1bPbQQ7nTj1srNEFWw+buZH5Kd3EC7NlcAjE47a9F5X8LIRYg9LXRZBu4P
2OKHKXoBawEAeDmxZzFs3rRTtN86yljWcojdjphC6M+lmBTjQQP1Z9m2HWbRIEAN
qmd7BE9VAgMBAAECggEAcn9EHFgoEZ7Xlz7gG/lc0zvO6HDeKxTky8QAzhey1Liq
+mllLHXlzBbcZWKy/R616nMBsDXGr4X/lzn0nNoWX/d2z+YgD8ND3mkVXpX12p33
S3WhdDVRmMI9TB5xfhbuPvNkzKof+7oR3uMvJQEVCka9CnpW6HE5mP+NZh9Dva39
5N278PHVrYC/tcRKZ6FVbqAP6/8D61D5JUL6fB2M2oDIEFW5CLJ+SiVgXQPSJEIp
leKnSyCjlZFOt7puniT6wGKbdjL1Azvh9rbyJNE/d2m8rCrW+NF8YW5BW+DKBcq4
tAyZ6bfLdjzQgS8tX0J7J2IBXkegJzaBII21DFwF/QKBgQD2SjQldzGR9a7s1d/b
HCbTABX95KeeldZF7doJ7mdValm9JRum7y5ezTz7heepixOybkbEfD5GXxqMw+iI
vKxNsJHEM1DcxYl2TGShtEAu+M8CkMLR4HmAB0NXGZf+j7vrStnrP1xuULQGmXCn
gEx60z6mBMj9L7rkM3ViQRmuVwKBgQD4rGNFH/8PUgzVVGneqj7Nkzi+1v95FV1V
qb4pNfYaoYSznd7Kmxm8CdnpCDWXEdsE/N8nXQOEv8yxqjz7eKRE0dt7ry3JHApm
GrK1sFWCoKYU1WTzsMRAx1+1yIFV+8e/2jfGOJs7OYgF7XWfXemN8SOOTfMga0CU
BcTUTYmMMwKBgAkYLFzFzfrJy6/LJfz9SSG5exZo/xeICOWEJuz+9Knex3mIiUDc
zCWMBphzFV7ZV9za3ZobzGDk2iIgidGixeoIfxlzm6TLVowbvguKkGQro9lAZzFA
zLvBiOcKduZbgGDr3qglKtCYAk3MFLDat/hBHKscuh8/R/NtQwRhywT5AoGBAMne
l8+7w8DaKNTda/x4U/dqtjwmKMpQz64E1/x4c8r2p9VwKTZsZp3BmYaVRXvX4ObR
qQ45cevIEBGCU3MJYsEDY5uqgA6slryAm+bmuOQMKgbrnMI/E3JK56WYmXYFqQhT
y8c8mLehYoz9UekHwduaj/SrztzYdFo1vK1kLG8FAoGBAPH8sTk8iIALu5gjyDiB
+o4W5LkskcdstiYz+XxBr3zZimzHWBRkzqVLnAAKAq4f9NFikd7byd3w7W1WairN
FleV79iBI776fxFNLK5KgxxrZAakXof5O2CdOBr7rQ60XZHlpwCX8iPSVfM5+DAC
2zQ3jOlc6ihDljfgkBH/TORk
-----END PRIVATE KEY-----
EOT
  }
}
`, name, namespace)
}
