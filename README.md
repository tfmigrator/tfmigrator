# tfmigrator

CLI tool to migrate Terraform configuration and State

## Overview

`tfmigrator` is a CLI tool to migrate Terraform configuration and State.

## Requirement

* Terraform
* [hcledit](https://github.com/minamijoyo/hcledit)

## Install

Download a binary from the [replease page](https://github.com/suzuki-shunsuke/tfmigrator/releases).

## How to use

```
$ terraform show -json > state.json
$ vi tfmigrator.yaml
$ cat *.tf | tfmigrator run [-skip-state] state.json
```

## Configuration file

[CONFIGURATION.md](docs/CONFIGURATION.md)

example of tfmigrator.yaml

```yaml
items:
- rule: |
    "name" not in Values
  exclude: true
- rule: |
    Values.name contains "foo"
  state_out: foo/terraform.tfstate
  resource_name: "{{.Values.name}}"
  tf_path: foo/resource.tf
- rule: |
    Values.name contains "bar"
  state_out: bar/terraform.tfstate
  resource_name: "{{.Values.name}}"
  tf_path: bar/resource.tf
```

## LICENSE

[MIT](LICENSE)
