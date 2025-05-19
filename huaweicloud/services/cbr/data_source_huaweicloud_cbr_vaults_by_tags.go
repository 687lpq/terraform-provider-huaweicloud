package cbr

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR POST /v3/{project_id}/vaults/resource_instances/action
func DataSourceVaultsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVaultsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource. If omitted, the provider-level region will be used.",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Operation identifier.",
			},
			"without_any_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether ignore tags params.",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether ignore tags params.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of included tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key of the resource tag.",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "All values corresponding to the key.",
						},
					},
				},
			},
			"tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of tags. Backups with any tags in this list will be filtered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key of the resource tag.",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "All values corresponding to the key.",
						},
					},
				},
			},
			"not_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded tags. Backups without these tags will be filtered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key of the resource tag.",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "All values corresponding to the key.",
						},
					},
				},
			},
			"not_tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of tags. Backups without any tags in this list will be filtered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key of the resource tag.",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "All values corresponding to the key.",
						},
					},
				},
			},
			"sys_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Which indicate the default enterprise project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key of the resource tag.",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "All values corresponding to the key.",
						},
					},
				},
			},
			"object_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether ignore tags params.",
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of the tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key of the tag.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the tag.",
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The count of the vaults.",
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of the vaults resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of the tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The key of the tag.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the tag.",
									},
								},
							},
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource name.",
						},
						"resource_detail": {
							Type:        schema.TypeList,
							Description: "Resource details.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vault": {
										Description: "Vault details.",
										Type:        schema.TypeList,
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Vault ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Vault name.",
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "User-defined vault description.",
												},
												"resources": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The list of the security groups found.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ID of the resource to be backed up.",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the resource to be backed up.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the resource to be backed up.",
															},
															"auto_protect": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to enable auto protect for the vault.",
															},
															"size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Allocated capacity for the associated resource, in GB.",
															},
															"backup_size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Backup size.",
															},
															"backup_count": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Number of backups.",
															},
															"protect_status": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Protection status.",
															},
															"extra_info": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Extra information of the resource.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"retention_duration": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The retention duration of the extra info.",
																		},
																		"name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The name of the extra info.",
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The description of the extra info.",
																		},
																		"exclude_volumes": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Elem:        &schema.Schema{Type: schema.TypeString},
																			Description: "All volumes of the extra info.",
																		},
																	},
																},
															},
														},
													},
												},
												"provider_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of the vault resource type.",
												},
												"created_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Creation time, for example, 2020-02-05T10:38:34.209782.",
												},
												"project_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Project ID.",
												},
												"enterprise_project_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Enterprise project ID. Its default value is 0.",
												},
												"auto_bind": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether automatic association is enabled.",
												},
												"bind_rules": {
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "Structured association rules.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"auto_expand": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to enable auto capacity expansion for the vault.",
												},
												"smn_notify": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Exception notification function.",
												},
												"threshold": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Vault capacity threshold.",
												},
												"user_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "User ID.",
												},
												"billing": {
													Type:        schema.TypeList,
													Description: "Vault details.",
													Computed:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"allocated": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Allocated capacity, in GB.",
															},
															"charging_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Billing mode, which can be post_paid or pre_paid.",
															},
															"cloud_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cloud type, which can be public or hybrid.",
															},
															"consistent_level": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Vault specification, which can be crash_consistent or app_consistent.",
															},
															"frozen_scene": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Scenario when an account is frozen.",
															},
															"order_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Order ID.",
															},
															"product_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Product ID.",
															},
															"protect_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Protection type, which can be backup or replication.",
															},
															"object_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Object type, which can be server, disk and so on.",
															},
															"spec_code": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specification codeServer backup vault.",
															},
															"is_multi_az": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Multi-AZ attribute of a vault.",
															},
															"is_double_az": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Double-AZ attribute of a vault.",
															},
															"status": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Vault status.",
															},
															"storage_unit": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the bucket for the vault.",
															},
															"used": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Used capacity, in MB.",
															},
															"size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Allocated capacity for the associated resource, in GB.",
															},
															"partner_bp_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The bp ID of partner.",
															},
														},
													},
												},
												"tags": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The list of the tags.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The key of the tag.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The value of the tag.",
															},
														},
													},
												},
												"backup_name_prefix": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The prefix of backup name.",
												},
												"demand_billing": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the vault capacity can be exceeded.",
												},
												"cbc_delete_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Vault deletion count",
												},
												"frozen": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the vault is frozen.",
												},
												"availability_zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The zone which the vault located.",
												},
												"sys_lock_source_service": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Used to identify the SMB service.",
												},
												"supplier": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The supplier of the vault.",
												},
												"locked": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the vault is locked. A locked vault cannot be unlocked.",
												},
												"cross_account": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "The account which can cross this vault.",
												},
												"cross_account_urn": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cross-account uniform resource name.",
												},
											},
										},
									},
								},
							},
						},
						"sys_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Which indicate the default enterprise project.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key of the resource tag.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the tag.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildListVaultsByTagsBody(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{}
	if v, ok := d.GetOk("action"); ok {
		body["action"] = v.(string)
	}
	if v, ok := d.GetOk("without_any_tag"); ok {
		body["without_any_tag"] = v.(bool)
	}
	if v, ok := d.GetOk("cloud_type"); ok {
		body["cloud_type"] = v.(string)
	}
	if v, ok := d.GetOk("object_type"); ok {
		body["object_type"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		body["tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("tags_any"); ok {
		body["tags_any"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("not_tags"); ok {
		body["not_tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("not_tags_any"); ok {
		body["not_tags_any"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("sys_tags"); ok {
		body["sys_tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("matches"); ok {
		body["matches"] = expandTags(v.([]interface{}))
	}
	return body
}
func expandTags(rawTags []interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, len(rawTags))
	for i, raw := range rawTags {
		tag := raw.(map[string]interface{})
		tags[i] = map[string]interface{}{
			"key":    tag["key"].(string),
			"values": tag["values"].([]interface{}),
		}
	}
	return tags
}

func flattenAllVaultsByTags(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	resources := resp.([]interface{})
	if len(resources) == 0 {
		return nil
	}
	results := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		resourceMap := map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", resource, nil),
			"resource_name": utils.PathSearch("resource_name", resource, nil),
		}

		if vaultDetail := utils.PathSearch("resource_detail.vault", resource, nil); vaultDetail != nil {
			if vaultMap, ok := vaultDetail.(map[string]interface{}); ok {
				resourceMap["resource_detail"] = []interface{}{
					map[string]interface{}{
						"vault": []interface{}{
							map[string]interface{}{
								"id":                    utils.PathSearch("id", vaultMap, nil),
								"name":                  utils.PathSearch("name", vaultMap, nil),
								"provider_id":           utils.PathSearch("provider_id", vaultMap, nil),
								"created_at":            utils.PathSearch("created_at", vaultMap, nil),
								"project_id":            utils.PathSearch("project_id", vaultMap, nil),
								"enterprise_project_id": utils.PathSearch("enterprise_project_id", vaultMap, nil),
								"auto_bind":             utils.PathSearch("auto_bind", vaultMap, nil),
								"auto_expand":           utils.PathSearch("auto_expand", vaultMap, nil),
								"smn_notify":            utils.PathSearch("smn_notify", vaultMap, nil),
								"threshold":             utils.PathSearch("threshold", vaultMap, nil),
								"user_id":               utils.PathSearch("user_id", vaultMap, nil),
								"bind_rules":            utils.PathSearch("bind_rules", vaultMap, map[string]interface{}{}),
								"resources":             flattenVaultResources(utils.PathSearch("resources", vaultMap, make([]interface{}, 0))),
								"billing":               flattenVaultBilling(utils.PathSearch("billing", vaultMap, nil)),
								"tags":                  flattenVaultTags(utils.PathSearch("tags", vaultDetail, make([]interface{}, 0))),
							},
						},
					},
				}
			}
		}
		resourceMap["tags"] = flattenVaultTags(utils.PathSearch("tags", resource, make([]interface{}, 0)))

		results = append(results, resourceMap)
	}
	return results
}

func flattenVaultResources(resources interface{}) []interface{} {
	if resources == nil {
		return nil
	}

	resourceList := resources.([]interface{})
	rst := make([]interface{}, 0, len(resourceList))

	for _, res := range resourceList {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", res, nil),
			"type":           utils.PathSearch("type", res, nil),
			"name":           utils.PathSearch("name", res, nil),
			"auto_protect":   utils.PathSearch("auto_protect", res, nil),
			"size":           utils.PathSearch("size", res, nil),
			"backup_size":    utils.PathSearch("backup_size", res, nil),
			"backup_count":   utils.PathSearch("backup_count", res, nil),
			"protect_status": utils.PathSearch("protect_status", res, nil),
			"extra_info":     flattenVaultResourcesExtraInfo(utils.PathSearch("extra_info", res, make([]interface{}, 0))),
		})
	}
	return rst
}

func flattenVaultResourcesExtraInfo(tags interface{}) []map[string]interface{} {
	if tags == nil {
		return nil
	}

	tagMap, ok := tags.(map[string]interface{})
	if !ok {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(tagMap))
	for key, value := range tagMap {
		rst = append(rst, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}
	return rst
}

func flattenVaultBilling(billing interface{}) []interface{} {
	if billing == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"allocated":        utils.PathSearch("allocated", billing, nil),
			"cloud_type":       utils.PathSearch("cloud_type", billing, nil),
			"consistent_level": utils.PathSearch("consistent_level", billing, nil),
			"charging_mode":    utils.PathSearch("charging_mode", billing, nil),
			"order_id":         utils.PathSearch("order_id", billing, nil),
			"product_id":       utils.PathSearch("product_id", billing, nil),
			"protect_type":     utils.PathSearch("protect_type", billing, nil),
			"object_type":      utils.PathSearch("object_type", billing, nil),
			"spec_code":        utils.PathSearch("spec_code", billing, nil),
			"used":             utils.PathSearch("used", billing, 0),
			"status":           utils.PathSearch("status", billing, nil),
			"size":             utils.PathSearch("size", billing, 0),
			"is_multi_az":      utils.PathSearch("is_multi_az", billing, nil),
			"is_double_az":     utils.PathSearch("is_double_az", billing, nil),
			"storage_unit":     utils.PathSearch("storage_unit", billing, nil),
			"partner_bp_id":    utils.PathSearch("partner_bp_id", billing, nil),
		},
	}
}

func flattenVaultTags(tags interface{}) []map[string]interface{} {
	if tags == nil {
		return nil
	}

	tagList := tags.([]interface{})
	rst := make([]map[string]interface{}, 0, len(tagList))

	for _, tag := range tagList {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}
	return rst
}

func dataSourceVaultsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/vault/resource_instances/action"
		product = "cbr"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	var allVaults []interface{}
	offset := 0
	totalCount := 0

	for {
		requestBody := buildListVaultsByTagsBody(d)
		if requestBody["action"] == "filter" {
			requestBody["limit"] = "1000"
			requestBody["offset"] = strconv.Itoa(offset)
		}

		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         requestBody,
		}
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying CBR vaults by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		vaults := utils.PathSearch("resources", respBody, []interface{}{}).([]interface{})
		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		if len(vaults) == 0 {
			break
		}
		allVaults = append(allVaults, vaults...)
		offset += len(vaults)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_count", totalCount),
		d.Set("resources", flattenAllVaultsByTags(allVaults)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the CBR vault by tags: %s", mErr)
	}
	return nil
}
