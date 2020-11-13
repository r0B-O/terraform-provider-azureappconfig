---
layout: "docs"
page_title: "Inspecting State - Terraform CLI"
---

# Inspecting State

Terraform includes some commands for reading and updating state without taking
any other actions.

- [The `terraform state list` command](/docs/commands/state/list.html)
  shows the resource addresses for every resource Terraform knows about in a
  configuration, optionally filtered by partial resource address.

- [The `terraform state show` command](/docs/commands/state/show.html)
  displays detailed state data about one resource.

- [The `terraform refresh` command](/docs/commands/refresh.html) updates
  state data to match the real-world condition of the managed resources. This is
  done automatically during plans and applies, but not when interacting with
  state directly.
