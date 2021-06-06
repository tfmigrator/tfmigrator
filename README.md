# tfmigrator

[![Build Status](https://github.com/tfmigrator/tfmigrator/workflows/test/badge.svg)](https://github.com/tfmigrator/tfmigrator/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/tfmigrator/tfmigrator)](https://goreportcard.com/report/github.com/tfmigrator/tfmigrator)
[![GitHub last commit](https://img.shields.io/github/last-commit/tfmigrator/tfmigrator.svg)](https://github.com/tfmigrator/tfmigrator)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/tfmigrator/tfmigrator/main/LICENSE)

Go library to migrate Terraform Configuration and State with `terraform state mv` and `terraform state rm` command and [hcledit](https://github.com/minamijoyo/hcledit).

## Requirement

* Go
* Terraform

[hcledit](https://github.com/minamijoyo/hcledit) isn't needed.

## Compared with suzuki-shunsuke/tfmigrator

This is a Go package. On the other hand, [suzuki-shunsuke/tfmigrator](https://github.com/suzuki-shunsuke/tfmigrator) is CLI tool.

Originally we have been developing suzuki-shunsuke/tfmigrator before this package.
In suzuki-shunsuke/tfmigrator, we define rules for migration with YAML, Go's [text/template](https://golang.org/pkg/text/template/), and [expr](https://github.com/antonmedv/expr).
So we don't have to write code with Go.

But when we migrate many resources, it is hard to write the configuration file.
So we started to develop tfmigrator as Go package instead of suzuki-shunsuke/tfmigrator.

Using this package, we can implement the migration rules with Go freely.
By implementing the interface [Planner](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#Planner), we can migrate resources.
This is very simple and powerful and flexible.
And we don't have to learn the configuration file format and expr's Language specification, so the learning cost is low.
Using [QuickRun](https://pkg.go.dev/github.com/tfmigrator/tfmigrator/tfmigrator#QuickRun) function, we can implement CLI for migration quickly.

And compared with suzuki-shunsuke/tfmigrator v1.0.0, this package provides rich feature.

* Support to remove Resource
* Support to update Terraform Configuration files in place
* Don't have to install [hcledit](https://github.com/minamijoyo/hcledit) command

## Example

Please see [examples](examples).

## LICENSE

[MIT](LICENSE)
