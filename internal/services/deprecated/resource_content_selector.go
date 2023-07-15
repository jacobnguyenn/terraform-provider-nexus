package deprecated

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nexus "github.com/nduyphuong/go-nexus-client/nexus3"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/security"
)

func ResourceContentSelector() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated. Please use the resource nexus_security_content_selector instead.",
		Description: `!> This resource is deprecated. Please use the resource "nexus_security_content_selector" instead.

Use this resource to create a Nexus Content Selector.`,

		Create: resourceContentSelectorCreate,
		Read:   resourceContentSelectorRead,
		Update: resourceContentSelectorUpdate,
		Delete: resourceContentSelectorDelete,
		Exists: resourceContentSelectorExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Content selector name",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "A description of the content selector",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"expression": {
				Description: "The content selector expression",
				Required:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func getContentSelectorFromResourceData(d *schema.ResourceData) security.ContentSelector {
	contentSelector := security.ContentSelector{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		contentSelector.Description = description.(string)
	}

	return contentSelector
}

func setContentSelectorToResourceData(contentSelector *security.ContentSelector, d *schema.ResourceData) error {
	d.SetId(contentSelector.Name)
	d.Set("description", contentSelector.Description)
	d.Set("expression", contentSelector.Expression)
	d.Set("name", contentSelector.Name)
	return nil
}

func resourceContentSelectorCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	contentSelector := getContentSelectorFromResourceData(d)

	if err := client.Security.ContentSelector.Create(contentSelector); err != nil {
		return err
	}

	d.SetId(contentSelector.Name)

	return resourceContentSelectorRead(d, m)
}

func resourceContentSelectorRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	contentSelector, err := client.Security.ContentSelector.Get(d.Id())
	if err != nil {
		return err
	}

	if contentSelector == nil {
		d.SetId("")
		return nil
	}

	return setContentSelectorToResourceData(contentSelector, d)
}

func resourceContentSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	contentSelector := getContentSelectorFromResourceData(d)
	if err := client.Security.ContentSelector.Update(d.Id(), contentSelector); err != nil {
		return err
	}

	return resourceContentSelectorRead(d, m)
}

func resourceContentSelectorDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.ContentSelector.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceContentSelectorExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	contentSelector, err := client.Security.ContentSelector.Get(d.Id())
	return contentSelector != nil, err
}
