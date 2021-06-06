# tfmigrator

Go library to migrate Terraform Configuration and State with `terraform state mv` and `terraform state rm` command and [hcledit](https://github.com/minamijoyo/hcledit).

## Requirement

* Go
* Terraform

[hcledit](https://github.com/minamijoyo/hcledit) isn't needed.

## Compared with tfmigrator-cli

https://github.com/suzuki-shunsuke/tfmigrator-cli/issues/25

tfmigrator is Go library. On the other hand, [tfmigrator-cli](https://github.com/suzuki-shunsuke/tfmigrator-cli) is CLI tool.

Originally we have been developing tfmigrator-cli before tfmigrator.
In tfmigrator-cli, we define rules for migration with YAML, Go's [text/template](https://golang.org/pkg/text/template/), [expr](https://github.com/antonmedv/expr).
So we don't have to write code with Go.

But when we migrate many resources, it is hard to write YAML.
So we started to develop tfmigrator instead of tfmigrator-cli.

tfmigrator is Go package, so we can implement the migration rules with Go.
In tfmigrator, we can migrate Terraform Configuration and State by implementing the interface `Planner`.
This is very simple and powerful and flexible.
And we don't have to learn the configuration file format and expr's Language specification, so the learning cost is low.
Using `QuickRun` function, we can implement CLI for migration quickly.

And compared with tfmigrator v1.0.0, tfmigrator provides rich feature.

* Support to remove Resource
* Support to update Terraform Configuration files in place
* Don't have to install [hcledit](https://github.com/minamijoyo/hcledit) command

## Example

Please see [examples](examples).

## LICENSE

[MIT](LICENSE)
