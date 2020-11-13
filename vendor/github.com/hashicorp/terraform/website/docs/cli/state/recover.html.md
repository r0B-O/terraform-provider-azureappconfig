---
layout: "docs"
page_title: "Recovering from State Disasters - Terraform CLI"
---

# Recovering from State Disasters

If something has gone horribly wrong (possibly due to accidents when performing
other state manipulation actions), you might need to take drastic actions with
your state data.

- [The `terraform force-unlock` command](/docs/commands/force-unlock.html) can
  override the protections Terraform uses to prevent two processes from
  modifying state at the same time. You might need this if a Terraform process
  (like a normal apply) is unexpectedly terminated (like by the complete
  destruction of the VM it's running in) before it can release its lock on the
  state backend. Do not run this until you are completely certain what happened
  to the process that caused the lock to get stuck.

- [The `terraform state pull` command](/docs/commands/state/pull.html) and
  [the `terraform state push` command](/docs/commands/state/push.html) can
  directly read and write entire state files from and to the configured backend.
  You might need this for obtaining or restoring a state backup.

