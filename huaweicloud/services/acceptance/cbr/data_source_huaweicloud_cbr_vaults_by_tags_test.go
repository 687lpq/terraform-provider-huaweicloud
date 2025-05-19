package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVaultsByTags_basic(t *testing.T) {
	var (
		all      = "data.huaweicloud_cbr_vaults_by_tags.test"
		dc       = acceptance.InitDataSourceCheck(all)
		testName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTagsVaultsByTags_basic1(testName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "resources.0.resource_name", fmt.Sprintf("%s_0", testName)),
				),
			},
			{
				Config: testAccDataTagsVaultsByTags_basic2(testName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "total_count", "2"),
				),
			},
		},
	})
}

func testAccDataVaultsByTags_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  count = 2
  
  name                  = format("%[1]s_%%d", count.index)
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "0"
  
  tags = {
    foo = format("bar%%d", count.index)
  }
}
`, name)
}

func testAccDataTagsVaultsByTags_basic1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cbr_vaults_by_tags" "test" {
  depends_on = [huaweicloud_cbr_vault.test]
  action     = "filter"
  tags {
    key    = "foo"
    values = ["bar0","bar1"]
  }
}

output "generated_names" {
  value = [for vault in huaweicloud_cbr_vault.test[*] : vault.name]
}
output "retrieved_names" {
  value = [for vault in data.huaweicloud_cbr_vaults_by_tags.test.resources[*] : vault.resource_name]
}
`, testAccDataVaultsByTags_base(name))
}

func testAccDataTagsVaultsByTags_basic2(name string) string {
	return fmt.Sprintf(`
%[1]s
data "huaweicloud_cbr_vaults_by_tags" "test" {
  depends_on = [huaweicloud_cbr_vault.test]
  action     = "count"
  tags {
    key    = "foo"
    values = ["bar0","bar1"]
  }
}
output "tags_validation" {
  value = data.huaweicloud_cbr_vaults_by_tags.test.total_count
}
`, testAccDataVaultsByTags_base(name))
}