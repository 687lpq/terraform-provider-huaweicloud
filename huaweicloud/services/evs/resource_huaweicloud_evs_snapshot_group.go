package evs

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS POST /v5/{project_id}/snapshot-groups
// @API EVS PUT /v5/{project_id}/snapshot-groups/{snapshot_group_id}
// @API EVS GET /v5/{project_id}/snapshot-groups/{snapshot_group_id}
// @API EVS DELETE /v5/{project_id}/snapshot-groups/{snapshot_group_id}
// @API EVS POST /v5/{project_id}/snapshot-groups/{snapshot_group_id}/tags/create
// @API EVS POST /v5/{project_id}/snapshot-groups/{snapshot_group_id}/tags/delete
func ResourceEvsSnapshotGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsSnapshotGroupCreate,
		ReadContext:   resourceEvsSnapshotGroupRead,
		UpdateContext: resourceEvsSnapshotGroupUpdate,
		DeleteContext: resourceEvsSnapshotGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instant_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"incremental": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildCreateSnapshotGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"server_id":             utils.ValueIgnoreEmpty(d.Get("server_id")),
		"volume_ids":            d.Get("volume_ids"),
		"instant_access":        d.Get("instant_access"),
		"name":                  d.Get("name"),
		"description":           d.Get("description"),
		"enterprise_project_id": d.Get("enterprise_project_id"),
		"incremental":           d.Get("incremental"),
	}
	return map[string]interface{}{"snapshot_group": utils.RemoveNil(body)}
}

func buildUpdateSnapshotGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
	return map[string]interface{}{"snapshot_group": utils.RemoveNil(body)}
}

func GetSnapshotGroupDetail(client *golangsdk.ServiceClient, snapshotGroupID string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/snapshot-groups/{snapshot_group_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_group_id}", snapshotGroupID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForSnapshotGroupStatusAvailable(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetSnapshotGroupDetail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("snapshot_group.status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", errors.New("status is not found in API response")
			}

			if status == "available" {
				return respBody, "COMPLETED", nil
			}

			if status == "error" {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitingForSnapshotGroupDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetSnapshotGroupDetail(client, d.Id())
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "success deleted", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			status := utils.PathSearch("snapshot_group.status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", errors.New("status is not found in API response")
			}

			if status == "error_deleting" {
				return respBody, status, errors.New("an error occurred while deleting the EVS snapshot group. " +
					"Please check with your cloud admin or check the API logs " +
					"to see why this error occurred")
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func checkVolumeSnapshotsStatus(client *golangsdk.ServiceClient) (bool, error) {
	// Query all snapshots with creating status
	requestPath := client.Endpoint + "v5/{project_id}/snapshots/detail"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += "?status=creating&limit=1000"

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return false, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return false, err
	}

	snapshots := utils.PathSearch("snapshots", respBody, []interface{}{}).([]interface{})

	if len(snapshots) > 0 {
		return false, nil
	}

	return true, nil
}

func waitingForVolumeSnapshotsComplete(ctx context.Context, client *golangsdk.ServiceClient,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			completed, err := checkVolumeSnapshotsStatus(client)
			if err != nil {
				return nil, "ERROR", err
			}

			if completed {
				return "all snapshots completed", "COMPLETED", nil
			}

			return "snapshots still being created", "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEvsSnapshotGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshot-groups"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	// Check if there are any snapshots being created
	if err := waitingForVolumeSnapshotsComplete(ctx, client, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for volume snapshots to complete: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSnapshotGroupBodyParams(d)),
	}
	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating EVS snapshot group: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	snapshot_group_id := utils.PathSearch("snapshot_group_id", respBody, "").(string)
	if snapshot_group_id != "" {
		d.SetId(snapshot_group_id)
	} else {
		return diag.Errorf("error creating EVS snapshot group: ID is not found in API response")
	}

	if err := waitingForSnapshotGroupStatusAvailable(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for EVS snapshot group (%s) creation to available: %s", d.Id(), err)
	}

	// Add tags after snapshot group is created and available
	if tags := d.Get("tags").([]interface{}); len(tags) > 0 {
		if err := updateSnapshotGroupTags(client, d.Id(), tags, true); err != nil {
			return diag.Errorf("error adding tags to EVS snapshot group (%s): %s", d.Id(), err)
		}
	}

	return resourceEvsSnapshotGroupRead(ctx, d, meta)
}

func resourceEvsSnapshotGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshot-groups/{snapshot_group_id}"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EVS snapshot group")
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	group := utils.PathSearch("snapshot_group", respBody, map[string]interface{}{}).(map[string]interface{})
	mErr := multierror.Append(
		d.Set("snapshot_group", flattenEvsSnapshotGroup(group)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEvsSnapshotGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshot-groups/{snapshot_group_id}"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateSnapshotGroupBodyParams(d)),
	}
	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating EVS snapshot group: %s", err)
	}
	if d.HasChange("tags") {
		old, new := d.GetChange("tags")
		if err := updateSnapshotGroupTags(client, d.Id(), old.([]interface{}), false); err != nil {
			return diag.Errorf("error deleting old tags for EVS snapshot group (%s): %s", d.Id(), err)
		}
		if err := updateSnapshotGroupTags(client, d.Id(), new.([]interface{}), true); err != nil {
			return diag.Errorf("error creating new tags for EVS snapshot group (%s): %s", d.Id(), err)
		}
	}
	return resourceEvsSnapshotGroupRead(ctx, d, meta)
}

func resourceEvsSnapshotGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshot-groups/{snapshot_group_id}"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting EVS snapshot group")
	}

	if err := waitingForSnapshotGroupDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for EVS snapshot group (%s) deleted: %s", d.Id(), err)
	}

	return nil
}

func buildSnapshotGroupBody(tags []interface{}) map[string]interface{} {
	body := map[string]interface{}{}
	if len(tags) > 0 {
		body["tags"] = expandSnapshotGroupTags(tags)
	}
	return body
}

func expandSnapshotGroupTags(rawTags []interface{}) []interface{} {
	tags := make([]interface{}, len(rawTags))
	for i, raw := range rawTags {
		tag := raw.(map[string]interface{})
		tags[i] = map[string]interface{}{
			"key":   tag["key"].(string),
			"value": tag["value"].(string),
		}
	}
	return tags
}

func updateSnapshotGroupTags(client *golangsdk.ServiceClient, groupID string, tags []interface{}, isCreate bool) error {
	if len(tags) == 0 {
		return nil
	}
	var (
		apiPath string
	)
	if isCreate {
		apiPath = client.Endpoint + "v5/{project_id}/snapshot-groups/{snapshot_group_id}/tags/create"
	} else {
		apiPath = client.Endpoint + "v5/{project_id}/snapshot-groups/{snapshot_group_id}/tags/delete"
	}
	apiPath = strings.ReplaceAll(apiPath, "{project_id}", client.ProjectID)
	apiPath = strings.ReplaceAll(apiPath, "{snapshot_group_id}", groupID)
	requestBody := buildSnapshotGroupBody(tags)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
		OkCodes:          []int{200, 201, 202, 204},
	}
	_, err := client.Request("POST", apiPath, &requestOpt)
	return err
}

func flattenEvsSnapshotGroup(group map[string]interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, 1)

	if len(group) > 0 {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", group, nil),
			"created_at":            utils.PathSearch("created_at", group, nil),
			"status":                utils.PathSearch("status", group, nil),
			"updated_at":            utils.PathSearch("updated_at", group, nil),
			"name":                  utils.PathSearch("name", group, nil),
			"description":           utils.PathSearch("description", group, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", group, nil),
			"tags":                  utils.PathSearch("tags", group, map[string]interface{}{}).(map[string]interface{}),
			"server_id":             utils.PathSearch("server_id", group, nil),
		})
	}

	return rst
}
