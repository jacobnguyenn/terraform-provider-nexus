package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nduyphuong/terraform-provider-nexus/internal/schema/common"
	"github.com/nduyphuong/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryAptHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing apt repository.",

		Read: dataSourceRepositoryAptHostedRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Hosted schemas
			"cleanup":   repository.DataSourceCleanup,
			"component": repository.DataSourceComponent,
			"storage":   repository.DataSourceHostedStorage,
			// Apt hosted schemas
			"distribution": {
				Description: "Distribution to fetch",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceRepositoryAptHostedRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceAptHostedRepositoryRead(resourceData, m)
}
