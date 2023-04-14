terraform {
  required_providers {
    ubika = {
      source = "registry.terraform.io/ubika/ubika"
    }
  }
}

resource "ubika_asset" "example" {
  metadata = {
    namespace = "default"
    name      = "terraform-test-asset"
  }
  spec = {
    hostnames       = ["terraform.example.com"]
    backend_url     = "http://terraform.example.com/"
    deployment_type = "SAAS"
  }
}
