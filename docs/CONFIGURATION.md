# Configuration

path | type | required | default | description
--- | --- | --- | --- | ---
.items | item | true | | 

## type: item

path | type | required | default | example | description
--- | --- | --- | --- | --- | ---
rule | bool expression | true | | `Type == "null_resource"` | If the result is `true`, the resource is proceeded by the item
state_out | string | true | | `foo/terraform.tfstate` |
tf_path | string | true | | `foo/resource.tf` |
resource_name | template | false | | `{{.Values.tags.Name}}` | If this isn't empty, the resource is renamed to this value
exclude | bool | false | false | | If this is true, resources which match the item are ignored 

## type: bool expression

[expr](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md) expression.
The expression must be returnes boolean (true or false).

## type: template

Go's [text/template](https://golang.org/pkg/text/template/)

[sprig](http://masterminds.github.io/sprig/) function can be used.

## expression and template parameter

The output of `terraform state -json` is passed.

path | type | example | description
--- | --- | --- | ---
.Name | string | `foo` | Terraform resource name
.Address | string | `aws_iam_user.foo` | Terraform resource address
.Type | string | `null_resource` | Terraform resource type
.Values | `map[string]interface{}` | `{"id": "xxx"}` | Terraform resource attributes
