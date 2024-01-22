package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpenAPIResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOpenAPIResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_openapi.test", "metadata.name", "tf-acc-test"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "ubika_openapi.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.namespace", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccOpenAPIResourceConfig("tf-acc-test", "tf-acc-tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ubika_openapi.test", "metadata.namespace", "tf-acc-tests"),
				),
			},
			// // Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOpenAPIResourceConfig(name string, namespace string) string {
	return fmt.Sprintf(`
resource "ubika_openapi" "test" {
  metadata = {
    name = %[1]q
    namespace = %[2]q
  }
  spec = {
	source = <<EOT
{  "openapi": "3.0.0",  "info": { "version": "1.0.0", "title": "Swagger Petstore", "license": {"name": "MIT" }  },  "servers": [ {"url": "/" }  ],  "paths": { "/pets": {"get": {  "summary": "List all pets",  "operationId": "listPets",  "tags": [ "pets"  ],  "parameters": [ {"name": "limit","in": "query","description": "How many items to return at one time (max 100)","required": false,"schema": {  "type": "integer",  "maximum": 100,  "format": "int32"} }  ],  "responses": { "200": {"description": "A paged array of pets","headers": {  "x-next": { "description": "A link to the next page of responses", "schema": {"type": "string" }  }},"content": {  "application/json": { "schema": {"$ref": "#/components/schemas/Pets" }  }} }, "default": {"description": "unexpected error","content": {  "application/json": { "schema": {"$ref": "#/components/schemas/Error" }  }} }  }},"post": {  "summary": "Create a pet",  "operationId": "createPets",  "tags": [ "pets"  ],  "responses": { "201": {"description": "Null response" }, "default": {"description": "unexpected error","content": {  "application/json": { "schema": {"$ref": "#/components/schemas/Error" }  }} }  }} }, "/pets/{petId}": {"get": {  "summary": "Info for a specific pet",  "operationId": "showPetById",  "tags": [ "pets"  ],  "parameters": [ {"name": "petId","in": "path","required": true,"description": "The id of the pet to retrieve","schema": {  "type": "string"} }  ],  "responses": { "200": {"description": "Expected response to a valid request","content": {  "application/json": { "schema": {"$ref": "#/components/schemas/Pet" }  }} }, "default": {"description": "unexpected error","content": {  "application/json": { "schema": {"$ref": "#/components/schemas/Error" }  }} }  }} }  },  "components": { "schemas": {"Pet": {  "type": "object",  "required": [ "id", "name"  ],  "properties": { "id": {"type": "integer","format": "int64" }, "name": {"type": "string" }, "tag": {"type": "string" }  }},"Pets": {  "type": "array",  "maxItems": 100,  "items": { "$ref": "#/components/schemas/Pet"  }},"Error": {  "type": "object",  "required": [ "code", "message"  ],  "properties": { "code": {"type": "integer","format": "int32" }, "message": {"type": "string" }  }} }  }}
EOT
	}
}
`, name, namespace)
}
