# Module 1

Run `terraform apply` to create `terraform.tfstate`.

```
$ terraform init
$ terraform apply
```

Run dry run.

```console
$ go run main.go -dry-run -log-level debug main.tf                   
2021/06/07 19:31:17 [DEBUG] + terraform show -json
2021/06/07 19:31:19 [DEBUG] + hcledit block list -f main.tf
2021/06/07 19:31:19 [INFO] [DRYRUN] + terraform state mv module.foo module.bar
2021/06/07 19:31:19 [DEBUG] [DRY RUN] + hcledit block mv -u module.foo module.bar
migrated_resources:
- source_address: module.foo
  source_tf_file_path: main.tf
  new_address: module.bar
removed_resources: []
not_migrated_resources:
- address: null_resource.foo
  file_path: main.tf
  attributes:
    id: "8064984026557147891"
    triggers: null
```

Run.

```console
$ go run main.go -log-level debug main.tf                   
2021/06/07 19:40:42 [DEBUG] + terraform show -json
2021/06/07 19:40:44 [DEBUG] + hcledit block list -f main.tf
2021/06/07 19:40:44 [INFO] + terraform state mv module.foo module.bar
2021/06/07 19:40:44 [DEBUG] + hcledit block mv -u module.foo module.bar
migrated_resources:
- source_address: module.foo
  source_tf_file_path: main.tf
  new_address: module.bar
removed_resources: []
not_migrated_resources:
- address: null_resource.foo
  file_path: main.tf
  attributes:
    id: "8064984026557147891"
    triggers: null
```
