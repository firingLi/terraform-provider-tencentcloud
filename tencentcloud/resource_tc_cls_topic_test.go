package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClsTopic_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsTopicExists("tencentcloud_cls_topic.topic"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.topic", "topic_name", "tf-topic-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_topic.topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS topic][Exists] check: CLB topic %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS topic][Exists] check: CLB topic id is not set")
		}
		clsService := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := clsService.DescribeClsTopicById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("[CHECK][CLS topic][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsTopic = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-topic-test"
  tags        = {
    "test" = "test"
  }
}

resource "tencentcloud_cls_topic" "topic" {
  auto_split           = true
  logset_id            = tencentcloud_cls_logset.logset.id
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test"
  }
  topic_name           = "tf-topic-test"
}
`
