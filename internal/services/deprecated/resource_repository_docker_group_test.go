package deprecated_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/repository"
	"github.com/nduyphuong/terraform-provider-nexus/internal/acceptance"
)

func testAccResourceRepositoryDockerGroup(httpPort int, httpsPort int) repository.LegacyRepository {
	repo := testAccResourceRepositoryGroup(repository.RepositoryFormatDocker)
	repo.Docker = &repository.Docker{
		ForceBasicAuth: true,
		HTTPPort:       &httpPort,
		HTTPSPort:      &httpsPort,
		V1Enabled:      false,
	}
	return repo
}

func TestAccResourceRepositoryDockerGroup(t *testing.T) {
	hostedRepo := testAccResourceRepositoryDockerHostedWithPorts(8280, 8633)
	hostedRepoResName := testAccResourceRepositoryName(hostedRepo)

	proxyRepo := testAccResourceDockerProxy()
	proxyRepoResName := testAccResourceRepositoryName(proxyRepo)

	repo := testAccResourceRepositoryDockerGroup(8180, 8533)
	repo.Group.MemberNames = []string{
		fmt.Sprintf("%s.name", hostedRepoResName),
		fmt.Sprintf("%s.name", proxyRepoResName),
	}
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryConfig(hostedRepo) + testAccResourceRepositoryConfig(proxyRepo) + testAccResourceRepositoryConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceRepositoryTestCheckFunc(repo),
					resourceRepositoryTypeGroupTestCheckFunc(repo),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "docker.#", "1"),
						resource.TestCheckResourceAttr(resName, "docker.0.force_basic_auth", strconv.FormatBool(repo.Docker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(resName, "docker.0.http_port", strconv.Itoa(*repo.Docker.HTTPPort)),
						resource.TestCheckResourceAttr(resName, "docker.0.https_port", strconv.Itoa(*repo.Docker.HTTPSPort)),
						resource.TestCheckResourceAttr(resName, "docker.0.v1enabled", strconv.FormatBool(repo.Docker.V1Enabled)),
					),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     repo.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
