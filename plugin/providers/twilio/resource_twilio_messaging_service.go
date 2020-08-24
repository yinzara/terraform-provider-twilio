package twilio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevinburke/twilio-go"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func resourceTwilioMessagingService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioMessagingServiceCreate,
		Read:   resourceTwilioMessagingServiceRead,
		Update: resourceTwilioMessagingServiceUpdate,
		Delete: resourceTwilioMessagingServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this messaging service.",
			},
			"friendly_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A friendly, human-readable name by which you can refer to this messaging service.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the messaging service was created.",
			},
			"date_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the messaging service was last updated.",
			},
			"inbound_request_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL we call using inbound_method when a message is received by any phone number or short code in the Service. When this property is null, receiving inbound messages is disabled. All messages sent to the Twilio phone number or short code will not be logged and received on the Account.",
			},
			"inbound_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "POST",
				Description: "The HTTP method we use to call inbound_request_url. Can be GET or POST.",
			},
			"fallback_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL that we call using fallback_method if an error occurs while retrieving or executing the TwiML from the Inbound Request URL.",
			},
			"fallback_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "POST",
				Description: "The HTTP method we use to call fallback_url. Can be: GET or POST.",
			},
			"status_callback": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL we call to pass status updates about message delivery.",
			},
			"sticky_sender": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable Sticky Sender on the Service instance.",
			},
			"mms_converter": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable the MMS Converter for messages sent through the Service instance.",
			},
			"smart_encoding": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable Smart Encoding for messages sent through the Service instance.",
			},
			"fallback_to_long_code": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable Fallback to Long Code for messages sent through the Service instance.",
			},
			"area_code_geomatch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable Area Code Geomatch on the Service Instance.",
			},
			"synchronous_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Reserved.",
			},
			"validity_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "How long, in seconds, messages sent from the Service are valid. Can be an integer from 1 to 14,400.",
			},
		},
	}
}

func makeCreateServiceRequestPayload(d *schema.ResourceData) url.Values {
	createServiceRequestPayload := make(url.Values)

	addIfNotEmpty(createServiceRequestPayload, "FriendlyName", d.Get("friendly_name"))
	addIfNotEmpty(createServiceRequestPayload, "InboundRequestUrl", d.Get("inbound_request_url"))
	addIfNotEmpty(createServiceRequestPayload, "InboundMethod", d.Get("inbound_method"))
	addIfNotEmpty(createServiceRequestPayload, "FallbackUrl", d.Get("fallback_url"))
	addIfNotEmpty(createServiceRequestPayload, "FallbackMethod", d.Get("fallback_method"))
	addIfNotEmpty(createServiceRequestPayload, "StatusCallback", d.Get("status_callback"))
	addIfNotEmpty(createServiceRequestPayload, "StickySender", d.Get("sticky_sender"))
	addIfNotEmpty(createServiceRequestPayload, "MmsConverter", d.Get("mms_converter"))
	addIfNotEmpty(createServiceRequestPayload, "SmartEncoding", d.Get("smart_encoding"))
	addIfNotEmpty(createServiceRequestPayload, "FallbackToLongCode", d.Get("fallback_to_long_code"))
	addIfNotEmpty(createServiceRequestPayload, "AreaCodeGeomatch", d.Get("area_code_geomatch"))
	addIfNotEmpty(createServiceRequestPayload, "ValidityPeriod", d.Get("validity_period"))
	addIfNotEmpty(createServiceRequestPayload, "SynchronousValidation", d.Get("synchronous_validation"))

	return createServiceRequestPayload
}

func mapTwilioMessagingServiceToTerraform(ms *twilio.Service, d *schema.ResourceData) error {
	err := d.Set("sid", ms.Sid)
	if err == nil {
		err = d.Set("account_sid", ms.AccountSid)
	}
	if err == nil && ms.FriendlyName != "" {
		err = d.Set("friendly_name", ms.FriendlyName)
	}
	if err == nil && ms.DateCreated.Valid {
		err = d.Set("date_created", ms.DateCreated.Time.Format("2006-01-02T15:04:05-07:00"))
	}
	if ms.DateUpdated.Valid {
		err = d.Set("date_updated", ms.DateUpdated.Time.Format("2006-01-02T15:04:05-07:00"))
	}
	if err == nil && ms.InboundRequestURL != nil {
		err = d.Set("inbound_request_url", *ms.InboundRequestURL)
	}
	if err == nil {
		err = d.Set("inbound_method", ms.InboundMethod)
	}
	if err == nil {
		err = d.Set("fallback_url", ms.FallbackURL)
	}
	if err == nil {
		err = d.Set("fallback_method", ms.FallbackMethod)
	}
	if err == nil {
		err = d.Set("status_callback", ms.StatusCallback)
	}
	if err == nil {
		err = d.Set("mms_converter", ms.MMSConverter)
	}
	if err == nil {
		err = d.Set("smart_encoding", ms.SmartEncoding)
	}
	if err == nil {
		err = d.Set("sticky_sender", ms.StickySender)
	}
	if err == nil {
		err = d.Set("fallback_to_long_code", ms.FallbackToLongCode)
	}
	if err == nil {
		err = d.Set("area_code_geomatch", ms.AreaCodeGeomatch)
	}
	if err == nil {
		err = d.Set("validity_period", ms.ValidityPeriod)
	}
	if err == nil {
		err = d.Set("synchronous_validation", ms.SynchronousValidation)
	}
	return err
}

func resourceTwilioMessagingServiceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioMessagingServiceCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration

	params := makeCreateServiceRequestPayload(d)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.Message.Services.Create")

	result, err := client.Message.Services.Create(context.TODO(), params)

	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).Error("Caught an error when attempting to create messaging service: " + err.Error())

		return err
	}

	d.SetId(result.Sid)

	err = mapTwilioMessagingServiceToTerraform(result, d)

	if err != nil {
		return fmt.Errorf("Encountered error while reading result for create of messaging service SID %s and mapping it to TF: %s", result.Sid, err)
	}

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"service_sid": result.Sid,
		},
	).Debug("END client.Message.Services.Create")

	return nil
}

func resourceTwilioMessagingServiceRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioMessagingServiceRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration

	log.Debug("Getting SID")

	sid := d.Id()

	log.Debug("Getting messaging service")

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"service_sid": sid,
		},
	).Debug("START client.IncomingNumbers.Get")

	ph, err := client.Message.Services.Get(context.TODO(), sid)

	if err != nil {
		return fmt.Errorf("Encountered an error when getting messaging service SID %s: %s", sid, err)
	}

	err = mapTwilioMessagingServiceToTerraform(ph, d)

	if err != nil {
		return fmt.Errorf("Encountered an error while mapping Twilio API result to terraform: %s", err)
	}

	return nil
}

func resourceTwilioMessagingServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioMessagingServiceDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration

	sid := d.Id()

	updatePayload := makeCreateServiceRequestPayload(d)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"service_sid": sid,
		},
	).Debug("START client.Message.Services.Update")

	_, err := client.Message.Services.Update(context.TODO(), sid, updatePayload)

	if err != nil {
		return fmt.Errorf("Failed to update messaging service SID %s: %s", sid, err)
	}

	return nil
}

func resourceTwilioMessagingServiceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioMessagingServiceDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration

	sid := d.Id()

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"service_sid": sid,
		},
	).Debug("START client.Message.Services.Release")

	err := client.Message.Services.Delete(context.TODO(), sid)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"service_sid": sid,
		},
	).Debug("END client.Message.Services.Release")

	if err != nil {
		return fmt.Errorf("Failed to delete messaging service: %s", err.Error())
	}

	return nil
}
