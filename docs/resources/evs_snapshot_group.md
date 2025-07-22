---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot_group"
description: |-
  Manages a resource to create, update, show, and delete EVS snapshot group within HuaweiCloud.
---

# huaweicloud_evs_snapshot_group

Manages a resource to create, update, show, and delete EVS snapshot group within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the deleted snapshot
group, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "region" {}
variable "name" {}
variable "description" {}
variable "volume_ids" {}
variable "tags" {}
variable "enterprise_project_id" {}
variable "server_id" {}
variable "instant_access" {}
variable "incremental" {}

resource "huaweicloud_evs_snapshot_group" "example" {
  region                = var.region
  name                  = var.name
  description           = var.description
  volume_ids            = var.volume_ids
  tags                  = var.tags
  enterprise_project_id = var.enterprise_project_id
  server_id             = var.server_id
  instant_access        = var.instant_access
  incremental           = var.incremental
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `snapshot_groups` - (Required, List) Specifies the snapshot group list.
  The [snapshot_groups](#snapshot_groups_structure) structure is documented below.

* `tags` - (Optional, List) Specifies the key/value pairs to be associated with the snapshot group.
  The [tags](#tags_struct) structure is documented below.

<a name="snapshot_groups_structure"></a>
The `snapshot_groups` block supports:

* `server_id` - (Optional, String) Specifies the server ID to which the snapshot group are attached.

* `volume_ids` - (Optional, List) Specifies the list of volume IDs to be included in the snapshot group. Must specified
  one of server_id or volume_ids.

* `instant_access` - (Optional, Boolean) Specifies whether to enable instant access for the snapshot group. Possible
  values are **true** (enable) and **false** (disable). Default is **false**.

* `name` - (Optional, String) Specifies the snapshot group name. The maximum length is `255` bytes.

* `description` - (Optional, String) Specifies the snapshot group description. The maximum length is `255` bytes.

* `tags` - (Optional, List) Specifies the key/value pairs to be associated with the snapshot group.
  The [tags](#tags_struct) structure is documented below.

* `enterprise_project_id` - (Required, String) Specifies the enterprise project ID for the snapshot group.

* `incremental` - (Required, Boolean) Specifies whether to create an incremental snapshot. Default is **false**.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Optional, String) The key of the tag.

* `value` - (Optional, String) The value of the tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `job_id` - Indicates the task ID.

* `snapshot_group_id` - The snapshot group ID.

* `snapshot_groups` - The snapshot group list.
  The [snapshot_groups](#snapshot_groups_structure) structure is documented below.

<a name="snapshot_groups_structure"></a>
The `snapshot_groups` block supports:

* `id` - The snapshot group ID.

* `created_at` - The time when the snapshot group was created.

* `status` - The snapshot group status.

* `updated_at` - The time when the snapshot group was updated.

* `name` - The snapshot group name.

* `description` - The snapshot group description.

* `enterprise_project_id` - The ID of the enterprise project associated with the snapshot.

* `tags` - The tags of the snapshot group.

* `server_id` - (Optional, String) Specifies the server ID to which the snapshot group are attached.
