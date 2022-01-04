package blobstore_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceBlobstoreFile(t *testing.T) {
	resourceName := "nexus_blobstore_file.acceptance"

	bs := blobstore.File{
		Name: fmt.Sprintf("test-blobstore-%s", acctest.RandString(5)),
		Path: "/nexus-data/acceptance",
		SoftQuota: &blobstore.SoftQuota{
			Limit: int64(acctest.RandIntRange(100, 300) * 1000000),
			Type:  "spaceRemainingQuota",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreFileConfig(bs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", bs.Name),
					resource.TestCheckResourceAttr(resourceName, "name", bs.Name),
					resource.TestCheckResourceAttr(resourceName, "path", bs.Path),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "soft_quota.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "soft_quota.0.limit", strconv.FormatInt(bs.SoftQuota.Limit, 10)),
						resource.TestCheckResourceAttr(resourceName, "soft_quota.0.type", bs.SoftQuota.Type),
					),
					resource.TestCheckResourceAttrSet(resourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(resourceName, "total_size_in_bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "available_space_in_bytes"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           bs.Name,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"available_space_in_bytes"},
			},
		},
	})
}

func testAccResourceBlobstoreFileConfig(bs blobstore.File) string {
	return fmt.Sprintf(`
resource "nexus_blobstore_file" "acceptance" {
	name = "%s"
	path = "%s"

	soft_quota {
		limit = %d
		type  = "%s"
	}
}`, bs.Name, bs.Path, bs.SoftQuota.Limit, bs.SoftQuota.Type)
}