---
layout: "docs"
page_title: "Command: import"
sidebar_current: "docs-commands-import"
description: |-
  The `terraform import` command is used to import existing resources into Terraform.
---

# Command: import

> **Hands-on:** Try the [Import Terraform Configuration](https://learn.hashicorp.com/tutorials/terraform/state-import?in=terraform/state&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) tutorial on HashiCorp Learn.

The `terraform import` command is used to
[import existing resources](/docs/import/index.html)
into Terraform.

## Usage

Usage: `terraform import [options] ADDRESS ID`

Import will find the existing resource from ID and import it into your Terraform
state at the given ADDRESS.

ADDRESS must be a valid [resource address](/docs/internals/resource-addressing.html).
Because any resource address is valid, the import command can import resources
into modules as well as directly into the root of your state.

ID is dependent on the resource type being imported. For example, for AWS
instances it is the instance ID (`i-abcd1234`) but for AWS Route53 zones
it is the zone ID (`Z12ABC4UGMOZ2N`). Please reference the provider documentation for details
on the ID format. If you're unsure, feel free to just try an ID. If the ID
is invalid, you'll just receive an error message.

~> Warning: Terraform expects that each remote object it is managing will be
bound to only one resource address, which is normally guaranteed by Terraform
itself having created all objects. If you import existing objects into Terraform,
be careful to import each remote object to only one Terraform resource address.
If you import the same object multiple times, Terraform may exhibit unwanted
behavior. For more information on this assumption, see
[the State section](/docs/state/).

The command-line flags are all optional. The list of available flags are:

* `-backup=path` - Path to backup the existing state file. Defaults to
  the `-state-out` path with the ".backup" extension. Set to "-" to disable
  backups.

* `-config=path` - Path to directory of Terraform configuration files that
  configure the provider for import. This defaults to your working directory.
  If this directory contains no Terraform configuration files, the provider
  must be configured via manual input or environmental variables.

* `-input=true` - Whether to ask for input for provider configuration.

* `-lock=true` - Lock the state file when locking is supported.

* `-lock-timeout=0s` - Duration to retry a state lock.

* `-no-color` - If specified, output won't contain any color.

* `-parallelism=n` - Limit the number of concurrent operation as Terraform
  [walks the graph](/docs/internals/graph.html#walking-the-graph). Defaults
  to 10.

* `-provider=provider` - **Deprecated** Override the provider configuration to
use when importing the object. By default, Terraform uses the provider specified
in the configuration for the target resource, and that is the best behavior in most cases.

* `-state=path` - Path to the source state file to read from. Defaults to the
  configured backend, or "terraform.tfstate".

* `-state-out=path` - Path to the destination state file to write to. If this
  isn't specified the source state file will be used. This can be a new or
  existing path.

* `-var 'foo=bar'` - Set a variable in the Terraform configuration. This flag
  can be set multiple times. Variable values are interpreted as
  [HCL](/docs/configuration/syntax.html#HCL), so list and map values can be
  specified via this flag. This is only useful with the `-config` flag.

* `-var-file=foo` - Set variables in the Terraform configuration from
  a [variable file](/docs/configuration/variables.html#variable-files). If
  a `terraform.tfvars` or any `.auto.tfvars` files are present in the current
  directory, they will be automatically loaded. `terraform.tfvars` is loaded
  first and the `.auto.tfvars` files after in alphabetical order. Any files
  specified by `-var-file` override any values set automatically from files in
  the working directory. This flag can be used multiple times. This is only
  useful with the `-config` flag.

## Provider Configuration

Terraform will attempt to load configuration files that configure the
provider being used for import. If no configuration files are present or
no configuration for that specific provider is present, Terraform will
prompt you for access credentials. You may also specify environmental variables
to configure the provider.

The only limitation Terraform has when reading the configuration files
is that the import provider configurations must not depend on non-variable
inputs. For example, a provider configuration cannot depend on a data
source.

As a working example, if you're importing AWS resources and you have a
configuration file with the contents below, then Terraform will configure
the AWS provider with this file.

```hcl
variable "access_key" {}
variable "secret_key" {}

provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}
```

## Example: Import into Resource

This example will import an AWS instance into the `aws_instance` resource named `foo`:

```shell
$ terraform import aws_instance.foo i-abcd1234
```

## Example: Import into Module

The example below will import an AWS instance into the `aws_instance` resource named `bar` into a module named `foo`:

```shell
$ terraform import module.foo.aws_instance.bar i-abcd1234
```

## Example: Import into Resource configured with count

The example below will import an AWS instance into the first instance of the `aws_instance` resource named `baz` configured with
[`count`](/docs/configuration/resources.html#count-multiple-resource-instances-by-count):

```shell
$ terraform import 'aws_instance.baz[0]' i-abcd1234
```

## Example: Import into Resource configured with for_each

The example below will import an AWS instance into the `"example"` instance of the `aws_instance` resource named `baz` configured with
[`for_each`](/docs/configuration/resources.html#for_each-multiple-resource-instances-defined-by-a-map-or-set-of-strings):

Linux, Mac OS, and UNIX:

```shell
$ terraform import 'aws_instance.baz["example"]' i-abcd1234
```

PowerShell:

```shell
$ terraform import 'aws_instance.baz[\"example\"]' i-abcd1234
```

Windows `cmd.exe`:

```shell
$ terraform import aws_instance.baz[\"example\"] i-abcd1234
```
