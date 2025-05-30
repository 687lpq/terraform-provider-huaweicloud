package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbOpengaussBackupFiles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_backup_files.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussBackupFiles_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "files.#"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.download_link"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.link_expired_time"),
					resource.TestCheckResourceAttrSet(dataSource, "bucket"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussBackupFiles_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss_egress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  flavor                = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name                  = "%[2]s"
  password              = "Huangwei!120521"
  enterprise_project_id = "%[3]s"

  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

resource "huaweicloud_gaussdb_opengauss_backup" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
  name        = "%[2]s"
  description = "test description"
}
`, common.TestBaseNetwork(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGaussdbOpengaussBackupFiles_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_opengauss_backup_files" "test" {
  backup_id = huaweicloud_gaussdb_opengauss_backup.test.id
}
`, testDataSourceGaussdbOpengaussBackupFiles_base(name))
}
