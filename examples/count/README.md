# count

Run `terraform apply` to create `terraform.tfstate`.

```
$ terraform init
$ terraform apply
```

## Dry Run

```console
$ go run main.go -log-level debug -dry-run *.tf
2021/06/19 21:42:45 [DEBUG] + terraform show -json
2021/06/19 21:42:47 [DEBUG] + hcledit block list -f main.tf
2021/06/19 21:42:47 [INFO] [DRYRUN] + terraform state mv null_resource.foo[0] null_resource.foo0
2021/06/19 21:42:47 [INFO] [DRYRUN] + terraform state mv null_resource.foo[1] null_resource.foo1
2021/06/19 21:42:47 [INFO] [DRYRUN] + terraform state mv null_resource.foo[2] null_resource.foo2
2021/06/19 21:42:47 [INFO] [DRYRUN] + terraform state mv null_resource.foo[3] null_resource.foo3
migrated_resources:
- source_address: null_resource.foo[0]
  new_address: null_resource.foo0
- source_address: null_resource.foo[1]
  new_address: null_resource.foo1
- source_address: null_resource.foo[2]
  new_address: null_resource.foo2
- source_address: null_resource.foo[3]
  new_address: null_resource.foo3
removed_resources: []
not_migrated_resources: []
```

## Migrate

```console
$ go run main.go -log-level debug *.tf
2021/06/19 22:34:35 [DEBUG] + terraform show -json
2021/06/19 22:34:37 [DEBUG] + hcledit block list -f main.tf
2021/06/19 22:34:37 [INFO] + terraform state mv null_resource.foo[0] null_resource.foo0
2021/06/19 22:34:38 [INFO] + terraform state mv null_resource.foo[1] null_resource.foo1
2021/06/19 22:34:39 [INFO] + terraform state mv null_resource.foo[2] null_resource.foo2
2021/06/19 22:34:40 [INFO] + terraform state mv null_resource.foo[3] null_resource.foo3
migrated_resources:
- source_address: null_resource.foo[0]
  new_address: null_resource.foo0
- source_address: null_resource.foo[1]
  new_address: null_resource.foo1
- source_address: null_resource.foo[2]
  new_address: null_resource.foo2
- source_address: null_resource.foo[3]
  new_address: null_resource.foo3
removed_resources: []
not_migrated_resources: []
```
