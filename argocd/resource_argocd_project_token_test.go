package argocd

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccArgoCDProjectToken(t *testing.T) {
	project := "myproject"
	role := "test-role1234"
	count := 20

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// TODO: add expiry check
		// TODO: add token regeneration check
		// TODO: add token login check
		Steps: []resource.TestStep{
			{
				Config: testAccArgoCDProjectTokenSingle(project, role),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"argocd_project_token.single",
						"issued_at",
					),
				),
			},
			{
				Config: testAccArgoCDProjectTokenMultiple(project, role, count),
				Check: testCheckMultipleResourceAttrSet(
					"argocd_project_token.multiple",
					"issued_at",
					count,
				),
			},
		},
	})
}

func testAccArgoCDProjectTokenSingle(project string, role string) string {
	return fmt.Sprintf(`
resource "argocd_project_token" "single" {
  project = "%s"
  role    = "%s"
}
	`, project, role)
}
func testAccArgoCDProjectTokenMultiple(project string, role string, count int) string {
	return fmt.Sprintf(`
resource "argocd_project_token" "multiple" {
  count   = %d
  project = "%s"
  role    = "%s"
}
	`, count, project, role)
}

func testCheckMultipleResourceAttrSet(name, key string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for i := 0; i < count; i++ {
			ms := s.RootModule()
			_name := fmt.Sprintf("%s.%d", name, i)
			rs, ok := ms.Resources[_name]
			if !ok {
				return fmt.Errorf("not found: %s in %s", _name, ms.Path)
			}
			is := rs.Primary
			if is == nil {
				return fmt.Errorf("no primary instance: %s in %s", _name, ms.Path)
			}
			if val, ok := is.Attributes[key]; !ok || val == "" {
				return fmt.Errorf("%s: Attribute '%s' expected to be set", _name, key)
			}
		}
		return nil
	}
}
