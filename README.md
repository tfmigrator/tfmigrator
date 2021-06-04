# tfmigrator-sdk

Go library to migrate Terraform Configuration and State with `terraform state mv` command and [hcledit](https://github.com/minamijoyo/hcledit).

## Relation to tfmigrator

tfmigrator-sdk is Go library. On the other hand, [tfmigrator](https://github.com/suzuki-shunsuke/tfmigrator) is CLI tool.
Originally we developed tfmigrator before tfmigrator-sdk.
In tfmigrator, we define rules for migration as YAML, so we don't have to write code with Go.
But when we migrate a number of resources, we found it is hard to write YAML.
We improved tfmigrator to write rules frexibly, but we feel the restriction of YAML.

tfmigrator-sdk is Go library, so we can implement the migration rules with Go.
This is very powerful.
tfmigrator-sdk provides high level API to implement a command to migrate resources.

## LICENSE

[MIT](LICENSE)
