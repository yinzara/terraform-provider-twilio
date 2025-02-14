package twilio

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"

	log "github.com/sirupsen/logrus"
)

func resourceTwilioApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioApiKeyCreate,
		Read:   resourceTwilioApiKeyRead,
		Update: resourceTwilioApiKeyUpdate,
		Delete: resourceTwilioApiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sid": {
                Type:     schema.TypeString,
                Computed: true,
            },
			"friendly_name": {
                Type:     schema.TypeString,
                Optional: true,
            },
			"secret": {
                Type:      schema.TypeString,
                Computed:  true,
                Sensitive: true,
            },
			"date_created": {
                Type:     schema.TypeString,
                Computed: true,
            },
			"date_updated": {
                Type:     schema.TypeString,
                Computed: true,
            },
		},
	}
}

func flattenKeyForCreate(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("FriendlyName", d.Get("friendly_name").(string))

	return v
}

func resourceTwilioApiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioApiKeyCreate")

	client := meta.(*TerraformTwilioContext).client
	context := context.TODO()

	createParams := flattenKeyForCreate(d)

	log.Debug("START client.Keys.Create")

	createResult, err := client.Keys.Create(context, createParams)

	if err != nil {
		log.WithError(err).Error("client.Keys.Create failed")

		return err
	}

	d.SetId(createResult.Sid)
	d.Set("sid", createResult.Sid)
	d.Set("secret", createResult.Secret)
	d.Set("friendly_name", createResult.FriendlyName) // In the event that the name wasn't specified, Twilio generates one for you
	d.Set("date_created", createResult.DateCreated)
	d.Set("date_updated", createResult.DateUpdated)

	log.Debug("END client.Keys.Create")

	return nil
}

func resourceTwilioApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioApiKeyRead")

	client := meta.(*TerraformTwilioContext).client
	context := context.TODO()

	sid := d.Id()

	log.Debug("START client.Keys.Get")

	key, err := client.Keys.Get(context, sid)

	d.Set("sid", key.Sid)
	// Not updating the secret as Twilio only returns it on creation, not after
	d.Set("friendly_name", key.FriendlyName) // In the event that the name wasn't specified, Twilio generates one for you
	d.Set("date_created", key.DateCreated)
	d.Set("date_updated", key.DateUpdated)

	log.Debug("END client.Keys.Get")

	if err != nil {
		return fmt.Errorf("Failed to refresh key: %s", err.Error())
	}

	return nil
}

func resourceTwilioApiKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented")
}

func resourceTwilioApiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioApiKeyDelete")

	client := meta.(*TerraformTwilioContext).client
	context := context.TODO()

	sid := d.Id()

	log.Debug("START client.Keys.Delete")

	err := client.Keys.Delete(context, sid)

	log.Debug("END client.Accounts.Delete")

	if err != nil {
		return fmt.Errorf("Failed to delete key: %s", err.Error())
	}

	return nil
}
