package deprecated_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/security"
	"github.com/nduyphuong/terraform-provider-nexus/internal/acceptance"
)

func TestAccResourceRole(t *testing.T) {
	resName := "nexus_role.acceptance"

	role := security.Role{
		ID:          acctest.RandString(10),
		Name:        acctest.RandString(10),
		Description: acctest.RandString(30),
		Privileges:  []string{"nx-all"},
		Roles:       []string{"nx-admin"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			// Creates a basic role
			{
				Config: testAccResourceRoleConfig(role),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(resName, "id", role.ID),
					resource.TestCheckResourceAttr(resName, "name", role.Name),
					resource.TestCheckResourceAttr(resName, "roleid", role.ID),
					resource.TestCheckResourceAttr(resName, "description", role.Description),
					resource.TestCheckResourceAttr(resName, "privileges.#", strconv.Itoa(len(role.Privileges))),
					resource.TestCheckResourceAttr(resName, "roles.#", strconv.Itoa(len(role.Roles))),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     role.ID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceRoleConfig(role security.Role) string {
	return fmt.Sprintf(`
resource "nexus_role" "acceptance" {
	roleid = "%s"
	name   = "%s"
	description = "%s"
	privileges = ["%s"]
	roles = ["%s"]
}
`, role.ID, role.Name, role.Description, strings.Join(role.Privileges, "\",\""), strings.Join(role.Roles, "\",\""))
}
