package deprecated

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nexus "github.com/nduyphuong/go-nexus-client/nexus3"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/security"
)

func ResourceAnonymous() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated. Please use the resource nexus_security_anonymous instead.",
		Description: `!> This resource is deprecated. Please use the resource "nexus_security_anonymous" instead.

Use this resource to change the anonymous configuration of the nexus repository manager.`,

		Create: resourceAnonymousUpdate,
		Read:   resourceAnonymousRead,
		Update: resourceAnonymousUpdate,
		Delete: resourceAnonymousDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"enabled": {
				Description: "Activate the anonymous access to the repository manager, defaults to `false`  if unset",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"user_id": {
				Description: "The user id used by anonymous access, defaults to `anonymous` if unset",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "anonymous",
			},
			"realm_name": {
				Description: "The name of the used realm, defaults to `NexusAuthorizingRealm`  if unset",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NexusAuthorizingRealm",
			},
		},
	}
}

func getAnonymousFromResourceData(d *schema.ResourceData) security.AnonymousAccessSettings {
	return security.AnonymousAccessSettings{
		Enabled:   d.Get("enabled").(bool),
		UserID:    d.Get("user_id").(string),
		RealmName: d.Get("realm_name").(string),
	}
}

func setAnonymousToResourceData(anonymous *security.AnonymousAccessSettings, d *schema.ResourceData) error {
	d.SetId("anonymous")
	d.Set("enabled", anonymous.Enabled)
	d.Set("user_id", anonymous.UserID)
	d.Set("realm_name", anonymous.RealmName)
	return nil
}

func resourceAnonymousRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	anonymous, err := client.Security.Anonymous.Read()
	if err != nil {
		return err
	}

	return setAnonymousToResourceData(anonymous, d)
}

func resourceAnonymousUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	anonymous := getAnonymousFromResourceData(d)
	if err := client.Security.Anonymous.Update(anonymous); err != nil {
		return err
	}

	return resourceAnonymousRead(d, m)
}

func resourceAnonymousDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
