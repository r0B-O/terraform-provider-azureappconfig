---
layout: "docs"
page_title: "Managing Plugins - Terraform CLI"
---

# Managing Plugins

Terraform relies on plugins called "providers" in order to manage various types
of resources. (For more information about providers, see
[Providers](/docs/configuration/blocks/providers/index.html) in the Terraform
language docs.)

-> **Note:** Providers are currently the only plugin type most Terraform users
will interact with. Terraform also supports third-party provisioner plugins, but
we discourage their use.

Terraform downloads and/or installs any providers
[required](/docs/configuration/provider-requirements.html) by a configuration
when [initializing](/docs/cli/init/index.html) a working directory. By default,
this works without any additional interaction but requires network access to
download providers from their source registry.

You can configure Terraform's provider installation behavior to limit or skip
network access, and to enable use of providers that aren't available via a
networked source. Terraform also includes some commands to show information
about providers and to reduce the effort of installing providers in airgapped
environments.

## Configuring Plugin Installation

Terraform's configuration file includes options for caching downloaded plugins,
or explicitly specifying a local or HTTPS mirror to install plugins from. For
more information, see [CLI Config File](/docs/commands/cli-config.html).

## Getting Plugin Information

Use the [`terraform providers`](/docs/commands/providers.html) command to get information
about the providers required by the current working directory's configuration.

Use the [`terraform providers schema`](/docs/commands/providers/schema.html) command to
get machine-readable information about the resources and configuration options
offered by each provider.

## Managing Plugin Installation

Use the [`terraform providers mirror`](/docs/commands/providers/mirror.html) command to
download local copies of every provider required by the current working
directory's configuration. This directory will use the nested directory layout
that Terraform expects when installing plugins from a local source, so you can
transfer it directly to an airgapped system that runs Terraform.

Use the [`terraform providers lock`](/docs/commands/providers/lock.html) command
to update the lock file that Terraform uses to ensure predictable runs when
using ambiguous provider version constraints.
