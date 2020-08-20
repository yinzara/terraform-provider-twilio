package twilio

import (
    "context"
    "errors"
    "fmt"
    "github.com/hashicorp/terraform/helper/schema"
    log "github.com/sirupsen/logrus"
)

func dataTwilioMessagingService() *schema.Resource {
    s := makeComputed(resourceTwilioMessagingService().Schema)
    s["friendly_name"].Required = true
    s["friendly_name"].Computed = false

    return &schema.Resource{
        Read:   dataTwilioMessagingServiceRead,
        Schema: s,
    }
}

func dataTwilioMessagingServiceRead(d *schema.ResourceData, meta interface{}) error {
    log.Debug("ENTER dataTwilioMessagingServiceRead")

    client := meta.(*TerraformTwilioContext).client
    config := meta.(*TerraformTwilioContext).configuration
    ctx := context.TODO()

    query := make(map[string][]string)

    friendlyName := ""
    if f, ok := d.GetOkExists("friendly_name"); ok && len(f.(string)) > 0 {
        friendlyName = f.(string)
        query["FriendlyName"] = []string{friendlyName}
    }

    if friendlyName == "" {
        return errors.New("'friendly_name' must be specified")
    }

    log.WithFields(
        log.Fields{
            "account_sid": config.AccountSID,
            "friendly_name":      friendlyName,
        },
    ).Debug("START client.Message.Services.GetPage")

    if page, err := client.Message.Services.GetPage(ctx, query); err != nil {
        log.WithFields(
            log.Fields{
                "account_sid": config.AccountSID,
                "friendly_name":      friendlyName,
            },
        ).Debug("END client.Message.Services.GetPage")
        return fmt.Errorf("unable to find MessagingService with friendlyName: %s\nerror: %s", friendlyName, err.Error())
    } else {
        for _, service := range page.Services {
            d.SetId(service.Sid)
            log.WithFields(
                log.Fields{
                    "account_sid": config.AccountSID,
                    "service_sid":                d.Id(),
                },
            ).Debug("END client.Message.Services.GetPage")
            return mapTwilioMessagingServiceToTerraform(service, d)
        }
        err = fmt.Errorf("unable to find MessagingService with friendlyName: %s", friendlyName)
        return err
    }
}
