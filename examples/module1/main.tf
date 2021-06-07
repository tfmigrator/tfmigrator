# comment
resource "null_resource" "foo" {}

module "foo" {
  source = "./foo"
}
