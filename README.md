# tfmigrator-sdk

Go library to migrate Terraform Configuration and State with `terraform state mv` and `terraform state rm` command and [hcledit](https://github.com/minamijoyo/hcledit).

## Requirement

* Go >= 1.16
* Terraform

[hcledit](https://github.com/minamijoyo/hcledit) isn't needed.

## Compared with tfmigrator

tfmigrator-sdk is Go library. On the other hand, [tfmigrator](https://github.com/suzuki-shunsuke/tfmigrator) is CLI tool.

Originally we developed tfmigrator before tfmigrator-sdk.
In tfmigrator, we define rules for migration with YAML, Go's [text/template](https://golang.org/pkg/text/template/), [expr](https://github.com/antonmedv/expr).
So we don't have to write code with Go.

But when we migrate a number of resources, it is hard to write YAML.
So we started to develop tfmigrator-sdk instead of tfmigrator.

tfmigrator-sdk is Go library, so we can implement the migration rules with Go.
In tfmigrator-sdk, we can migrate Terraform Configuration and State by implementing the interface `Migrator`.
This is very simple and powerful and flexible.
And we don't have to configuration file format and expr's Language specification, so the learning cost is low.
Using `QuickRun` function, we can implement CLI for migration quickly.

And compared with tfmigrator v1.0.0, tfmigrator-sdk provides rich feature.

* Support to remove Resource
* Support to update Terraform Configuration files in place
* Don't have to install [hcledit](https://github.com/minamijoyo/hcledit) command

## Example

Please see [examples](examples).

## LICENSE

[MIT](LICENSE)
