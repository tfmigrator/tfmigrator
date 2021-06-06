# Example 1

Run `terraform apply` to create `terraform.tfstate`.

```
$ terraform init
$ terraform apply
```

```console
$ cat main.tf
# comment
resource "null_resource" "foo" {}
```

Run dry run.

```console
$ go run main.go -dry-run -log-level debug main.tf
2021/06/06 12:34:04 [DEBUG] + terraform show -json
2021/06/06 12:34:04 [DEBUG] + hcledit block list -f main.tf
2021/06/06 12:34:04 [INFO] [DRYRUN] + terraform state rm null_resource.bar
2021/06/06 12:34:04 [DEBUG] [DRY RUN] + hcledit block rm -u resource.null_resource.bar
2021/06/06 12:34:04 [INFO] [DRYRUN] + terraform state mv -state-out foo/terraform.tfstate null_resource.foo null_resource.foo
2021/06/06 12:34:04 [DEBUG] + hcledit block get -f main.tf resource.null_resource.foo
2021/06/06 12:34:04 [DEBUG] [DRY RUN] + hcledit block rm -u resource.null_resource.foo
2021/06/06 12:34:04 [INFO] [DRYRUN] + terraform state mv -state-out foo/terraform.tfstate null_resource.zoo null_resource.zoo
2021/06/06 12:34:04 [DEBUG] + hcledit block get -f main.tf resource.null_resource.zoo
2021/06/06 12:34:04 [DEBUG] [DRY RUN] + hcledit block rm -u resource.null_resource.zoo
migrated_resources:
- source_address: null_resource.foo
  source_tf_file_path: main.tf
  new_tf_file_basename: main.tf
  dirname: foo
  state_basename: terraform.tfstate
- source_address: null_resource.zoo
  source_tf_file_path: main.tf
  new_tf_file_basename: main.tf
  dirname: foo
  state_basename: terraform.tfstate
removed_resources:
- address: null_resource.bar
  file_path: main.tf
not_migrated_resources:
- address: null_resource.yoo
  file_path: main.tf
```

Run.

```console
$ go run main.go -log-level debug main.tf
2021/06/06 12:35:38 [DEBUG] + terraform show -json
2021/06/06 12:35:38 [DEBUG] + hcledit block list -f main.tf
2021/06/06 12:35:38 [INFO] + terraform state rm null_resource.bar
Removed null_resource.bar
Successfully removed 1 resource instance(s).
2021/06/06 12:35:38 [DEBUG] + hcledit block rm -u resource.null_resource.bar
2021/06/06 12:35:38 [INFO] + terraform state mv -state-out foo/terraform.tfstate null_resource.foo null_resource.foo
Move "null_resource.foo" to "null_resource.foo"
Successfully moved 1 object(s).
2021/06/06 12:35:38 [DEBUG] + hcledit block get -f main.tf resource.null_resource.foo
2021/06/06 12:35:38 [DEBUG] + hcledit block rm -u resource.null_resource.foo
2021/06/06 12:35:38 [INFO] + terraform state mv -state-out foo/terraform.tfstate null_resource.zoo null_resource.zoo
Move "null_resource.zoo" to "null_resource.zoo"
Successfully moved 1 object(s).
2021/06/06 12:35:38 [DEBUG] + hcledit block get -f main.tf resource.null_resource.zoo
2021/06/06 12:35:38 [DEBUG] + hcledit block rm -u resource.null_resource.zoo
migrated_resources:
- source_address: null_resource.foo
  source_tf_file_path: main.tf
  new_tf_file_basename: main.tf
  dirname: foo
  state_basename: terraform.tfstate
- source_address: null_resource.zoo
  source_tf_file_path: main.tf
  new_tf_file_basename: main.tf
  dirname: foo
  state_basename: terraform.tfstate
removed_resources:
- address: null_resource.bar
  file_path: main.tf
not_migrated_resources:
- address: null_resource.yoo
  file_path: main.tf
```
