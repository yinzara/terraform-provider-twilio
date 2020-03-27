package twilio

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sirupsen/logrus"
)

func dataTwilioPhoneNumber() *schema.Resource {
	s := makeComputed(resourceTwilioPhoneNumber().Schema)
	s["friendly_name"].Optional = true
	s["number"].Optional = true

	return &schema.Resource{
		Read:   dataTwilioPhoneNumberRead,
		Schema: s,
	}
}

func dataTwilioPhoneNumberRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER dataTwilioPhoneNumberRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	query := make(map[string][]string)

	friendlyName := ""
	if f, ok := d.GetOkExists("friendly_name"); ok && len(f.(string)) > 0 {
		friendlyName = f.(string)
		query["FriendlyName"] = []string{friendlyName}
	}
	number := ""
	if n, ok := d.GetOkExists("number"); ok && len(n.(string)) > 0 {
		number = n.(string)
		query["PhoneNumber"] = []string{number}
	}

	if number == "" && friendlyName == "" {
		return errors.New("'number' and/or 'friendly_name' must be specified")
	}

	log.WithFields(
		log.Fields{
			"parent_account_sid": config.AccountSID,
			"friendly_name":      friendlyName,
			"number":             number,
		},
	).Debug("START client.IncomingNumbers.GetPage")

	if page, err := client.IncomingNumbers.GetPage(context, query); err != nil {
        log.WithFields(
            log.Fields{
                "parent_account_sid": config.AccountSID,
                "friendly_name":      friendlyName,
            },
        ).Debug("END client.Accounts.GetPage")
		return fmt.Errorf("unable to find IncomingPhoneNumber with friendlyName: %s\nerror: %s", friendlyName, err.Error())
	} else {
	    for _, incNumber := range page.IncomingPhoneNumbers {
            if (friendlyName == "" || incNumber.FriendlyName == friendlyName) && (number == "" || string(incNumber.PhoneNumber) == number) {
                d.SetId(incNumber.Sid)
                log.WithFields(
                    log.Fields{
                        "parent_account_sid": config.AccountSID,
                        "sid":                d.Id(),
                    },
                ).Debug("END client.IncomingNumbers.GetPage")
                return mapTwilioPhoneNumberToTerraform(incNumber, d)
            }
        }
        if friendlyName == "" {
            err = fmt.Errorf("unable to find IncomingPhoneNumber with number: %s", number)
        } else if number == "" {
            err = fmt.Errorf("unable to find IncomingPhoneNumber with friendlyName: %s", friendlyName)
        } else {
            err = fmt.Errorf("unable to find IncomingPhoneNumber with number: %s and friendlyName: %s", number, friendlyName)
        }
        return err
	}
}
