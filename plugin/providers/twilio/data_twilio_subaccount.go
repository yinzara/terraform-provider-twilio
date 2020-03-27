package twilio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sirupsen/logrus"
)

func dataTwilioSubaccount() *schema.Resource {
	s := makeComputed(resourceTwilioSubaccount().Schema)
	s["friendly_name"].Required = true
	s["friendly_name"].Computed = false

	return &schema.Resource{
		Read:   dataTwilioSubaccountRead,
		Schema: s,
	}
}

func dataTwilioSubaccountRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER dataTwilioSubaccountRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	friendlyName := d.Get("friendly_name").(string)

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"friendly_name":      friendlyName,
		},
	).Debug("START client.Accounts.GetPage")

	if page, err := client.Accounts.GetPage(context, map[string][]string{
		"FriendlyName": {friendlyName},
	}); err != nil {
		log.WithFields(
			log.Fields{
				"parent_account_sid": config.AccountSID,
				"friendly_name":      friendlyName,
			},
		).WithError(err).Error("ERROR client.Accounts.GetPage")
		return fmt.Errorf("unable to find subaccount with friendlyName: %s\nerror: %s", friendlyName, err.Error())
	} else {
	    for _, account := range page.Accounts {
	        if account.FriendlyName == friendlyName {
                d.SetId(account.Sid)
                if err = mapTwilioSubaccountToTerraform(account, d); err == nil {
                    log.WithFields(
                        log.Fields{
                            "sid":                account.Sid,
                            "parent_account_sid": config.AccountSID,
                            "friendly_name":      friendlyName,
                        },
                    ).Debug("END client.Accounts.GetPage")
                    return nil
                } else {
                    log.WithFields(
                        log.Fields{
                            "parent_account_sid": config.AccountSID,
                            "sid":                account.Sid,
                            "friendly_name":      friendlyName,
                        },
                    ).WithError(err).Error("ERROR mapTwilioSubaccountToTerraform")
                    return fmt.Errorf("unable to map subaccount attributes for number with sid: %s\nerror: %s", d.Id(), err.Error())
                }
            }
        }
        return fmt.Errorf("unable to find subaccount with friendlyName: %s", friendlyName)
	}
}
