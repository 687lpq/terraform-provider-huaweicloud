package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
)

func getSnapshotGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "evs"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	return evs.GetSnapshotGroupDetail(client, state.Primary.ID)
}

func TestAccEvsSnapshotGroup_basic(t *testing.T) {
	var (
		snapshotGroup interface{}
		rName         = fmt.Sprintf("tf-acc-sg-%s", acctest.RandString(5))
		resourceName  = "huaweicloud_evs_snapshot_group.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&snapshotGroup,
		getSnapshotGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsSnapshotGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily group backup"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "foo"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "bar"),
				),
			},
			{
				Config: testAccEvsSnapshotGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated group backup"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "foo"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.1.key", "key"),
					resource.TestCheckResourceAttr(resourceName, "tags.1.value", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"job_id", "snapshot_group_id", "name", "description", "tags", "volume_ids", "instant_access",
					"enterprise_project_id", "incremental", "server_id", "region",
				},
			},
		},
	})
}

func testAccEvsSnapshotGroup_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s-vol"
  description       = "Created by acc test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}
`, rName)
}

func testAccEvsSnapshotGroup_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_snapshot_group" "test" {
  name        = "%[2]s"
  description = "Daily group backup"
  volume_ids  = [huaweicloud_evs_volume.test.id]
  
  tags {
    key   = "foo"
    value = "bar"
  }
}
`, testAccEvsSnapshotGroup_base(rName), rName)
}

func testAccEvsSnapshotGroup_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_snapshot_group" "test" {
  name        = "%[2]s-update"
  description = "Updated group backup"
  volume_ids  = [huaweicloud_evs_volume.test.id]
  
  tags {
    key   = "foo"
    value = "bar"
  }
  
  tags {
    key   = "key"
    value = "value"
  }
}
`, testAccEvsSnapshotGroup_base(rName), rName)
}
