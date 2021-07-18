# tfmigrator

[![Go Reference](https://pkg.go.dev/badge/github.com/tfmigrator/tfmigrator.svg)](https://pkg.go.dev/github.com/tfmigrator/tfmigrator)
[![Build Status](https://github.com/tfmigrator/tfmigrator/workflows/test/badge.svg)](https://github.com/tfmigrator/tfmigrator/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/tfmigrator/tfmigrator)](https://goreportcard.com/report/github.com/tfmigrator/tfmigrator)
[![GitHub last commit](https://img.shields.io/github/last-commit/tfmigrator/tfmigrator.svg)](https://github.com/tfmigrator/tfmigrator)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/tfmigrator/tfmigrator/main/LICENSE)

Go library to migrate Terraform Configuration and State with `terraform state mv` and `terraform state rm` command and [hcledit](https://github.com/minamijoyo/hcledit).

Using tfmigrator, we can do the following.

* migrate resource address
  * e.g. `null_resource.foo` => `null_resource.bar`
    * `terraform state mv null_resource.foo null_resource.bar`
    * `hcledit block mv -u -f main.tf resource.null_resource.foo resource.null_resource.bar`
  * e.g. `module.foo` => `module.bar`
* move resources in a Terraform Configuration file and State to the other file and State
  * e.g. `terraform.tfstate` => `foo/terraform.tfstate`
  * e.g. `main.tf` => `foo/main.tf`
* remove resources from State and Terraform Configuration
  * e.g. `terraform state rm null_resource.foo`
  * e.g. `hcledit block rm -u -f main.tf resource.null_resource.foo`

On the other hand, tfmigrator doesn't support the following things.

* change resource fields (attributes and blocks)
  * e.g. `hcledit attribute`

## Requirement

* Go
* Terraform

[hcledit](https://github.com/minamijoyo/hcledit) isn't needed.

## Example

Please see [examples](examples).

## Document

https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator

## Getting Started

[The example](https://github.com/tfmigrator/tfmigrator/blob/main/examples/example1/main.go) and [README](https://github.com/tfmigrator/tfmigrator/tree/main/examples/example1) tell us how to use tfmigrator and how tfmigrator works.
We can implement a command for migration with Go simply.
About the detail of [Planner](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#Planner), please see the document of [Source](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#Source) and [MigratedResource](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#MigratedResource).

[QuickRun](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#QuickRun) provides a high level API,
so if we use [QuickRun](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#QuickRun), we don't have to know about other API like [Runner](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#Runner).
But if we need low level API, please check other API like [Runner](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#Runner).

## Check tfmigrator/cli too

We have implemented CLI with this framework.

https://github.com/tfmigrator/cli

## Release Note (Change log)

https://github.com/tfmigrator/tfmigrator/releases

## LICENSE

[MIT](LICENSE)
