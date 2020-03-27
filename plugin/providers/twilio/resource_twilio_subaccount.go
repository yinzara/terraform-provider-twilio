package twilio

import (
	"context"
	"errors"
	"fmt"
	"github.com/kevinburke/twilio-go"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"

	log "github.com/sirupsen/logrus"
)

func resourceTwilioSubaccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioSubaccountCreate,
		Read:   resourceTwilioSubaccountRead,
		Update: resourceTwilioSubaccountUpdate,
		Delete: resourceTwilioSubaccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"parent_account_sid": {
                Type:     schema.TypeString,
                Computed: true,
            },
			"friendly_name": {
                Type:     schema.TypeString,
                Optional: true,
            },
			"status": {
                Type:     schema.TypeString,
                Optional: true,
                Default:  "active",
            },
			"auth_token": {
                Type:     schema.TypeString,
                Computed: true,
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

func flattenSubaccountForCreate(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("FriendlyName", d.Get("friendly_name").(string))

	return v
}

func flattenSubaccountForDelete(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("status", "closed")

	return v
}

func resourceTwilioSubaccountCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioSubaccountCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration

	createParams := flattenSubaccountForCreate(d)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
		},
	).Debug("START client.AccountsCreate")

	createResult, err := client.Accounts.Create(context.TODO(), createParams)

	if err != nil {
		log.WithFields(
			log.Fields{
				"parent_account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.AccountsCreate failed")

		return err
	}

	d.SetId(createResult.Sid)
	if err = mapTwilioSubaccountToTerraform(createResult, d); err != nil {
        log.WithFields(
            log.Fields{
                "parent_account_sid": config.AccountSID,
            },
        ).WithError(err).Error("ERROR mapTwilioSubaccountToTerraform")
	    return err
    }

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     createResult.Sid,
		},
	).Debug("END client.AccountsCreate")

	return nil
}

func resourceTwilioSubaccountRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioSubaccountRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration

	sid := d.Id()

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     sid,
		},
	).Debug("START client.Accounts.Get")

	account, err := client.Accounts.Get(context.TODO(), sid)
	if err == nil {
		err = mapTwilioSubaccountToTerraform(account, d)
	}

	if err != nil {
		return fmt.Errorf("Failed to refresh account: %s", err.Error())
	}

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     d.Id(),
		},
	).Debug("END client.AccountsGet")

	return nil
}

func mapTwilioSubaccountToTerraform(account *twilio.Account, d *schema.ResourceData) error {
	err := d.Set("status", account.Status)
	if err == nil {
		err = d.Set("auth_token", account.AuthToken)
	}
	if err == nil {
		err = d.Set("friendly_name", account.FriendlyName)
	}
	if err == nil && account.DateCreated.Valid {
		err = d.Set("date_created", account.DateCreated.Time.Format("2006-01-02T15:04:05-07:00"))
	}
	if err == nil && account.DateUpdated.Valid {
		err = d.Set("date_updated", account.DateUpdated.Time.Format("2006-01-02T15:04:05-07:00"))
	}
	if err == nil {
		err = d.Set("parent_account_sid", account.OwnerAccountSid)
	}
	return err
}

func resourceTwilioSubaccountUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented")
}

func resourceTwilioSubaccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioSubaccountDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration

	sid := d.Id()

	updateData := flattenSubaccountForDelete(d)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     sid,
		},
	).Debug("START client.Accounts.Delete")

	_, err := client.Accounts.Update(context.TODO(), sid, updateData)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"subaccount_sid":     sid,
		},
	).Debug("END client.Accounts.Delete")

	if err != nil {
		return fmt.Errorf("Failed to delete account: %s", err.Error())
	}

	return nil
}
