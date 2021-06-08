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

## Help

```console
$ go run main.go -help
tfmigrator - Migrate Terraform Configuration and State

Usage
  tfmigrator help
  tfmigrator [-help] [-dry-run] [-log-level debug] [-state ""] [Terraform Configuration file path ...]

Example

  $ ls *.tf | xargs tfmigrator -dry-run -log-level debug
```

## Dry Run

```console
$ go run main.go -dry-run
migrated_resources:
- source_address: null_resource.foo
  source_tf_file_path: main.tf
  new_address: null_resource.bar
not_migrated_resources: []
```

## Migrate

```console
$ go run main.go
2021/06/05 17:36:12 [INFO] + terraform state mv null_resource.foo null_resource.bar
Move "null_resource.foo" to "null_resource.bar"
Successfully moved 1 object(s).
```

Confirm that main.tf is updated without losing code comment.

```console
$ cat main.tf
# comment
resource "null_resource" "bar" {}
```
