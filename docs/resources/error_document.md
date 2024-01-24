---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ubika_error_document Resource - terraform-provider-ubika"
subcategory: ""
description: |-
  ErrorDocument resource
---

# ubika_error_document (Resource)

ErrorDocument resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `metadata` (Attributes) (see [below for nested schema](#nestedatt--metadata))
- `spec` (Attributes) (see [below for nested schema](#nestedatt--spec))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--metadata"></a>
### Nested Schema for `metadata`

Required:

- `name` (String) Name of the resource
- `namespace` (String) Namespace of the resource

Read-Only:

- `created` (Number)
- `updated` (Number)
- `version` (Number)


<a id="nestedatt--spec"></a>
### Nested Schema for `spec`

Required:

- `content_type` (String) Content type
- `page` (String) Page