---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_vaults_by_tags"
description: |-
  Use this data source to get available CBR vaults filtered by tags within Huaweicloud.
---

# huaweicloud_cbr_tags

Use this data source to get available CBR vaults filtered by tags within Huaweicloud.

## Example Usage

```hcl
data "huaweicloud_cbr_vaults_by_tags" "test" {
  action = "filter"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the action name. Possible values are **count** and **filter**.
  + **count**: querying count of data filtered by tags.
  + **filter**: querying details of data filtered by tags.

* `without_any_tag` - (Optional, Bool) Specifies whether ignore tags params.
  If this parameter is set to true, all resources without tags are queried.

* `tags` - (Optional, List) Specifies the instances that associate with all key/value pairs in this list.
  The [tags](#tags_struct) structure is documented below.

* `tags_any` - (Optional, List) Specifies the instances that associate with any key/value pairs in this list.
  The [tags_any](#tags_struct) structure is documented below.

* `not_tags` - (Optional, List) Specifies the instances that associate without all key/value pairs in this list.
  The [not_tags](#tags_struct) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the instances that associate without any key/value pairs in this list.
  The [not_tags_any](#tags_struct) structure is documented below.

* `sys_tags` - (Optional, List) Specifies the instances that associate with enterprise project.
  The [sys_tags](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the instances which condition matched.
  The [matches](#matches_struct) structure is documented below.

* `cloud_type` - (Optional, String) Specifies cloud type of the instances.

* `object_type` - (Optional, String) Specifies resource type of the instances.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Optional, String) The key of the resource tag.

* `values` - (Optional, List) All values corresponding to the key.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Optional, String) The key of the resource tag.

* `value` - (Optional, String) The value of the resource tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The count of the vaults.

* `resources` - The list of the security groups found.
  The [resources](#cbr_vaults_resources) structure is documented below.

<a name="cbr_vaults_resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_detail` - The detail of the security group.
  The [resource_detail](#cbr_vaults_resources_detail) structure is documented below.

* `tags` - The list of all tags for resources of the same type.
  The [tags](#cbr_vaults_tags) structure is documented below.

<a name="cbr_vaults_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

* `values` - All values corresponding to the key.

<a name="cbr_vaults_resources_detail"></a>

* `vault` - The attributes of all vault.
  The [vault](#cbr_vaults_resource_detail_vault) structure is documented below.

<a name="cbr_vaults_resource_detail_vault"></a>
The `vault` block supports:

* `id` - Vault UUID.

* `name` - Vault name.

* `description` - Nullable description field.

* `resources` - The attributes of all vault.
  The [resources](#cbr_vaults_resource_detail_vault_resources) structure is documented below.

* `provider_id` - The cloud provider identifier.

* `created_at` - Creation timestamp in ISO8601 format.

* `project_id` - The project ID that owns the vault.

* `enterprise_project_id` - The enterprise project ID.

* `auto_bind` - Whether automatic binding is enabled.

* `bind_rules` - The bind rule of vault.

* `auto_expand` - Whether automatic expansion is enabled.

* `smn_notify` - Whether SMN notifications are enabled.

* `threshold` - The usage percentage threshold.(0-100)

* `user_id` - The user ID who created the vault.

* `billing` - Capacity and billing info of the vault.
  The [billing](#cbr_vaults_resource_detail_vault_billing) structure is documented below.

* `tags` - The list of all tags for resources of the same type.
  The [tags](#cbr_vaults_resource_detail_vault_tags) structure is documented below.

* `backup_name_prefix` - The backup name prefix.

* `demand_billing` - Whether on-demand billing is enabled.

* `cbc_delete_count` - Vault deletion count.

* `frozen` - Whether the vault is frozen.

* `availability_zone` - The availability zone.

* `sys_lock_source_service` - Used to identify the SMB service. You can set it to SMB or leave it empty.

* `supplier` - Identifier of the resource provider.

* `locked` - Whether the vault is locked.

* `cross_account` - Whether cross-account access is enabled.

* `cross_account_urn` - Cross-account uniform resource name.

<a id="cbr_vaults_resource_detail_vault_resources"></a>
The `resources` block supports:

* `id` - The cloud server ID.

* `type` - The resource type.

* `name` - The resource display name.

* `auto_protect` - Whether auto-protection is enabled.

* `size` - The disk size in GB.

* `backup_size` - The backup size in GB.

* `backup_count` - The number of backups.

* `protect_status` - The protection status.

* `extra_info` - Additional resource metadata.
  The [extra_info](#cbr_vaults_resource_detail_vault_resources_extra_info) structure is documented below.

<a id="cbr_vaults_resource_detail_vault_billing"></a>
The `billing` block supports:

* `allocated` - The allocated storage in GB.

* `cloud_type` - The cloud type.

* `consistent_level` - The consistency level.

* `frozen_scene` - Scenario when an account is frozen.

* `charging_mode` - The billing mode.

* `order_id` - The order ID.

* `product_id` - The product ID.

* `protect_type` - The protection type.

* `object_type` - The object type.

* `spec_code` - The product specification code.

* `is_multi_az` - Whether multi-AZ is enabled.

* `is_double_az` - Whether dual-AZ is enabled.

* `used` - The used storage in MB.

* `storage_unit` - Name of the bucket for the vault.

* `partner_bp_id` - The bp ID of partner.

* `status` - The billing status.

* `size` - The total vault capacity in MB.

<a name="cbr_vaults_resource_detail_vault_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

* `value` - The value of the resource tag.

<a name="cbr_vaults_resource_detail_vault_resources_extra_info"></a>
The `extra_info` block supports:

* `retention_duration` - The retention duration of the extra info.

* `name` - The name of the extra info.

* `description` - The description of the extra info.

* `exclude_volumes` - All volumes of the extra info.
