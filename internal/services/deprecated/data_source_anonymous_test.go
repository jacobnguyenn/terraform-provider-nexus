package deprecated_test

import (
	"strconv"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/security"
)

func TestAccDataSourceAnonymous(t *testing.T) {
	dataSourceName := "data.nexus_anonymous.acceptance"

	anonym := security.AnonymousAccessSettings{
		Enabled:   true,
		UserID:    acctest.RandString(20),
		RealmName: "NexusAuthenticatingRealm",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAnonymousConfig(anonym),
				Check:  nil,
			},
			{
				Config: testAccResourceAnonymousConfig(anonym) + testAccDataSourceAnonymousConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "enabled", strconv.FormatBool(anonym.Enabled)),
					resource.TestCheckResourceAttr(dataSourceName, "user_id", anonym.UserID),
					resource.TestCheckResourceAttr(dataSourceName, "realm_name", anonym.RealmName),
				),
			},
		},
	})
}

func testAccDataSourceAnonymousConfig() string {
	return `
data "nexus_anonymous" "acceptance" {}
`
}
