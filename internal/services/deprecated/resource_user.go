package deprecated

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nexus "github.com/nduyphuong/go-nexus-client/nexus3"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/security"

	"github.com/nduyphuong/terraform-provider-nexus/internal/schema/common"
	"github.com/nduyphuong/terraform-provider-nexus/internal/tools"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated. Please use the resource nexus_security_user instead.",
		Description: `!> This resource is deprecated. Please use the resource "nexus_security_user" instead.

Use this resource to manage users.`,

		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"userid": {
				Description: "The userid which is required for login. This value cannot be changed.",
				ForceNew:    true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"firstname": {
				Description: "The first name of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"lastname": {
				Description: "The last name of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "The email address associated with the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "The password for the user.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"roles": {
				Description: "The roles which the user has been assigned within Nexus.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeSet,
			},
			"status": {
				Default:     "active",
				Description: "The user's status, e.g. active or disabled.",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"active",
					"disabled",
				}, false),
			},
		},
	}
}

func getUserFromResourceData(d *schema.ResourceData) security.User {
	return security.User{
		UserID:       d.Get("userid").(string),
		FirstName:    d.Get("firstname").(string),
		LastName:     d.Get("lastname").(string),
		EmailAddress: d.Get("email").(string),
		Password:     d.Get("password").(string),
		Status:       d.Get("status").(string),
		Roles:        tools.InterfaceSliceToStringSlice(d.Get("roles").(*schema.Set).List()),
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	user := getUserFromResourceData(d)

	if err := client.Security.User.Create(user); err != nil {
		return err
	}

	d.SetId(user.UserID)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	user, err := client.Security.User.Get(d.Id())
	if err != nil {
		return err
	}

	if user == nil {
		d.SetId("")
		return nil
	}

	d.Set("email", user.EmailAddress)
	d.Set("firstname", user.FirstName)
	d.Set("lastname", user.LastName)
	d.Set("roles", tools.StringSliceToInterfaceSlice(user.Roles))
	d.Set("status", user.Status)
	d.Set("userid", user.UserID)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if d.HasChange("password") {
		password := d.Get("password").(string)
		if err := client.Security.User.ChangePassword(d.Id(), password); err != nil {
			return err
		}
	}

	if d.HasChange("firstname") || d.HasChange("lastname") || d.HasChange("email") || d.HasChange("status") || d.HasChange("roles") {
		user := getUserFromResourceData(d)
		if err := client.Security.User.Update(d.Id(), user); err != nil {
			return err
		}
	}
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.User.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	user, err := client.Security.User.Get(d.Id())
	return user != nil, err
}
