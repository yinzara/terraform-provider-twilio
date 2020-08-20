package twilio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/kevinburke/twilio-go"
	"github.com/spf13/cast"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

func resourceTwilioPhoneNumber() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioPhoneNumberCreate,
		Read:   resourceTwilioPhoneNumberRead,
		Update: resourceTwilioPhoneNumberUpdate,
		Delete: resourceTwilioPhoneNumberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this phone number.",
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Look for this number sequence anywhere in the phone number.",
			},
			"area_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Look for a number within this area code.",
			},
			"country_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Two letter ISO country code in which you want to search for a number. See https://support.twilio.com/hc/en-us/articles/223183068-Twilio-international-phone-number-availability-and-their-capabilities for details on available countries.",
			},
			"number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full phone number, including country and area code.",
			},
			"friendly_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A friendly, human-readable name by which you can refer to this number.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the phone number was created.",
			},
			"date_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the phone number was laste updated.",
			},
			"service_sid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SID of the Service the resource is associated with.",
			},
			"address_requirements": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address requirements imposed on this number, if any.",
			},
			"is_beta": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this phone number is new to Twilio (beta status).",
			},
			"is_mms_capable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this phone number is MMS-capable.",
			},
			"is_sms_capable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this phone number is SMS-capable.",
			},
			"is_voice_capable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this phone number is voice-capable..",
			},
			"sms": {
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_sid": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "SID of the Twilio application to invoke when an SMS is sent to this number.",
						},
						"primary_http_method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "POST",
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"POST",
								"GET",
							}, false),
							Description: "The HTTP method for the primary URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
						"primary_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The URL called when an SMS is sent to this number.",
						},
						"fallback_http_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Default:     "POST",
							Description: "The HTTP method for the fallback URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
						"fallback_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The URL called if the primary URL returns a non-favorable status code.",
						},
					},
				},
			},
			"status_callback": {
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The URL called when a whenever a status change occurs on this number.",
						},
						"http_method": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Default:  "POST",
							ValidateFunc: validation.StringInSlice([]string{
								"POST",
								"GET",
							}, false),
							Description: "The HTTP method for the status callback URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
					},
				},
			},
			"voice": {
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_sid": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "SID of the Twilio application to invoke when a call is started with this number.",
						},
						"primary_http_method": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Default:  "POST",
							ValidateFunc: validation.StringInSlice([]string{
								"POST",
								"GET",
							}, false),
							Description: "The HTTP method for the primary URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
						"primary_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The URL called when a phone call starts on this number.",
						},
						"fallback_http_method": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Default:  "POST",
							ValidateFunc: validation.StringInSlice([]string{
								"POST",
								"GET",
							}, false),
							Description: "The HTTP method for the fallback URL. Can be `GET` or `POST`, defaults to `POST`.",
						},
						"fallback_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The URL called if the primary URL returns a non-favorable status code.",
						},
						"caller_id_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "If caller ID is enabled or not for this number. If enabled, incurs additional charge per call (see console for pricing). Can be `true` or `false`, defaults to `false`.",
						},
						"receive_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Determines if the line is set up for voice or fax. Can be `voice` or `fax`, defaults to `voice`.",
						},
					},
				},
			},
			"address_sid": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SID of the address associated with this phone number. May be required for certain countries.",
			},
			"trunk_sid": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SID of the voice trunk that will handle calls to this number. If set, overrides any voice URLs or applications: only the trunk will recieve the incoming call.",
			},
			"identity_sid": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SID of the identity associated with the phone number. May be required in certain countries.",
			},
			"emergency": {
				Type:     schema.TypeSet,
				MinItems: 0,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Default:      "Active",
							ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, true),
							Description:  "Status of this phone number. Either `Active` or `Inactive`.",
						},
						"address_sid": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "SID of the address used for emergency calling from this number. The address must be validated before it can be used for emergency purposes.",
						},
					},
				},
			},
		},
	}
}

func addIfNotEmpty(params url.Values, key string, value interface{}) {
	s := cast.ToString(value)

	if s != "" {
		params.Add(key, s)
	}
}

func makeCreateRequestPayload(d *schema.ResourceData) url.Values {
	createRequestPayload := make(url.Values)

	addIfNotEmpty(createRequestPayload, "FriendlyName", d.Get("friendly_name"))
	addIfNotEmpty(createRequestPayload, "AddressSid", d.Get("address_sid"))
	addIfNotEmpty(createRequestPayload, "TrunkSid", d.Get("trunk_sid"))
	addIfNotEmpty(createRequestPayload, "IdentitySid", d.Get("identity_sid"))

	if sms := d.Get("sms").(*schema.Set); sms.Len() > 0 {
		sms := sms.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "SmsApplicationSid", sms["application_sid"])
		addIfNotEmpty(createRequestPayload, "SmsFallbackUrl", sms["fallback_url"])
		addIfNotEmpty(createRequestPayload, "SmsFallbackMethod", sms["fallback_http_method"])
		addIfNotEmpty(createRequestPayload, "SmsMethod", sms["primary_http_method"])
		addIfNotEmpty(createRequestPayload, "SmsUrl", sms["primary_url"])
	}

	if voice := d.Get("voice").(*schema.Set); voice.Len() > 0 {
		voice := voice.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "VoiceApplicationSid", voice["application_sid"])
		addIfNotEmpty(createRequestPayload, "VoiceFallbackUrl", voice["fallback_url"])
		addIfNotEmpty(createRequestPayload, "VoiceFallbackMethod", voice["fallback_http_method"]) // TODO Map to safe values
		addIfNotEmpty(createRequestPayload, "VoiceMethod", voice["primary_http_method"])          // TODO Map to safe values
		addIfNotEmpty(createRequestPayload, "VoiceUrl", voice["primary_url"])
		addIfNotEmpty(createRequestPayload, "VoiceCallerIdLookup", voice["caller_id_enabled"])
		addIfNotEmpty(createRequestPayload, "VoiceReceiveMode", voice["receive_mode"]) // TODO Map to Twilio
	}

	if statusCallback := d.Get("status_callback").(*schema.Set); statusCallback.Len() > 0 {
		statusCallback := statusCallback.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "StatusCallbackMethod", statusCallback["http_method"]) // TODO Map to safe values
		addIfNotEmpty(createRequestPayload, "StatusCallback", statusCallback["url"])
	}

	if emergency := d.Get("emergency").(*schema.Set); emergency.Len() > 0 {
		emergency := emergency.List()[0].(map[string]interface{})

		addIfNotEmpty(createRequestPayload, "EmergencyStatus", emergency["status"]) // TODO Map to Twilio values
		addIfNotEmpty(createRequestPayload, "EmergencyAddressSid", emergency["address_sid"])
	}

	return createRequestPayload
}

func hashAnything(item interface{}) int {
	s := cast.ToString(item)
	return hashcode.String(s)
}

func hashStringKeyedMap(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		for _, value := range m {
			buf.WriteString(fmt.Sprintf("%s-", cast.ToString(value)))
		}
	}

	return hashcode.String(buf.String())
}

func mapTwilioPhoneNumberToTerraform(ph *twilio.IncomingPhoneNumber, d *schema.ResourceData) error {
	err := d.Set("sid", ph.Sid)
	if err == nil {
		err = d.Set("number", string(ph.PhoneNumber))
	}
	if err == nil {
		err = d.Set("friendly_name", ph.FriendlyName)
	}
	// d.Set("address_sid", p.AddressSid) -- address SID not in twiliogo
	// d.Set("identity_sid", p.IdentitySid) -- identity SID not in twiliogo
	if err == nil {
		err = d.Set("trunk_sid", ph.TrunkSid)
	}
	if err == nil && ph.DateCreated.Valid {
		err = d.Set("date_created", ph.DateCreated.Time.Format("2006-01-02T15:04:05-07:00"))
	}
	if ph.DateUpdated.Valid {
		err = d.Set("date_updated", ph.DateUpdated.Time.Format("2006-01-02T15:04:05-07:00"))
	}
	if err == nil {
		err = d.Set("address_requirements", ph.AddressRequirements)
	}
	if err == nil {
		err = d.Set("is_beta", ph.Beta)
	}
	if err == nil {
		err = d.Set("is_mms_capable", ph.Capabilities.MMS)
	}
	if err == nil {
		err = d.Set("is_sms_capable", ph.Capabilities.SMS)
	}
	if err == nil {
		err = d.Set("is_voice_capable", ph.Capabilities.Voice)
	}
	// d.Set("is_fax_capable", p.Capabilities.Fax) -- p.Capabilities.Fax not in twiliogo

	// Voice set
	if err == nil {
		voiceMap := make(map[string]interface{})
		voiceMap["application_sid"] = ph.VoiceApplicationSid
		voiceMap["fallback_url"] = ph.VoiceFallbackURL
		voiceMap["fallback_http_method"] = ph.VoiceFallbackMethod
		voiceMap["primary_url"] = ph.VoiceURL
		voiceMap["primary_http_method"] = ph.VoiceMethod
		voiceMap["caller_id_enabled"] = ph.VoiceCallerIDLookup
		// voiceMap["receive_mode"] = p.ReceiveMode -- receive mode not in twiliogo
		err = d.Set("voice", []map[string]interface{}{voiceMap})
	}
	if err == nil {
		// sms set
		smsMap := make(map[string]interface{})
		smsMap["application_sid"] = ph.SMSApplicationSid
		smsMap["fallback_url"] = ph.SMSFallbackURL
		smsMap["fallback_http_method"] = ph.SMSFallbackMethod
		smsMap["primary_url"] = ph.SMSURL
		smsMap["primary_http_method"] = ph.SMSMethod
		err = d.Set("sms", []map[string]interface{}{smsMap})
	}
	if err == nil {
		// status_callback
		statusCallbackMap := make(map[string]interface{})
		statusCallbackMap["url"] = ph.SMSFallbackURL
		statusCallbackMap["http_method"] = ph.SMSFallbackMethod
		err = d.Set("status_callback", []map[string]interface{}{statusCallbackMap})
	}
	if err == nil {
		// emergency
		emergencyMap := make(map[string]interface{})
		if ph.EmergencyAddressSid.Valid {
			emergencyMap["address_sid"] = ph.EmergencyAddressSid.String
		}
		emergencyMap["status"] = ph.EmergencyStatus
		err = d.Set("emergency", []map[string]interface{}{emergencyMap})
	}
	return err
}

func resourceTwilioPhoneNumberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	ctx := context.TODO()

	// Required parameter
	countryCode := d.Get("country_code").(string)

	//var searchParams url.Values
	searchParams := make(url.Values)

	serviceSid := cast.ToString(d.Get("service_sid"))
	areaCode := d.Get("area_code").(string)
	if len(areaCode) > 0 {
		searchParams.Set("AreaCode", areaCode)
	}

	search := d.Get("search").(string)
	if len(search) > 0 {
		searchParams.Set("Contains", search)
	}

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"country_code": countryCode,
		},
	).Debug("START client.Available.Numbers.Local.GetPage")

	// TODO switch based on the type of number to buy local, mobile, intl
	searchResult, err := client.AvailableNumbers.Local.GetPage(ctx, countryCode, searchParams)

	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid":  config.AccountSID,
				"country_code": countryCode,
			},
		).Error("Caught an unexpected error when searching for phone numbers")

		return err
	}

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"country_code": countryCode,
			"search":       search,
			"result_count": len(searchResult.Numbers),
		},
	).Debug("END client.Available.Nubmers.Local.GetPage")

	if searchResult != nil && len(searchResult.Numbers) == 0 {
		log.WithFields(
			log.Fields{
				"account_sid":  config.AccountSID,
				"country_code": countryCode,
			},
		).Error("No phone numbers matched the search patterns")

		return errors.New("No numbers found that match your search")
	}

	// Grab the first number that matches
	number := searchResult.Numbers[0]
	e164Number := string(number.PhoneNumber)

	buyParams := makeCreateRequestPayload(d)
	buyParams.Set("PhoneNumber", e164Number)

	log.WithFields(
		log.Fields{
			"account_sid":  config.AccountSID,
			"phone_number": e164Number,
		},
	).Debug("START client.IncomingNumbers.Create")

	buyResult, err := client.IncomingNumbers.Create(ctx, buyParams)

	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid":  config.AccountSID,
				"phone_number": e164Number,
			},
		).Error("Caught an error when attempting to purchase phone number: " + err.Error())

		return err
	}

	d.SetId(buyResult.Sid)
	d.Set("number", e164Number)

	err = mapTwilioPhoneNumberToTerraform(buyResult, d)

	if err != nil {
		return fmt.Errorf("Encountered error while reading buy result for phone number SID %s and mapping it to TF: %s", buyResult.Sid, err)
	}

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     e164Number,
			"phone_number_sid": buyResult.Sid,
		},
	).Debug("END client.IncomingNumbers.Create")

	if len(serviceSid) > 0 {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
				"phone_sid":   buyResult.Sid,
				"service_sid": serviceSid,
			},
		).Debug("START client.Message.Services.CreatePhoneNumber")
		_, err := client.Message.Services.CreatePhoneNumber(ctx, serviceSid, buyResult.Sid)
		if err != nil {
			return fmt.Errorf("Encountered error adding phone number with SID %s to messaging service with SID %s: %s", buyResult.Sid, serviceSid, err)
		}
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
				"phone_sid":   buyResult.Sid,
				"service_sid": serviceSid,
			},
		).Debug("END client.Message.Services.CreatePhoneNumber")
		err = d.Set("service_sid", serviceSid)
	}
	return err
}

func resourceTwilioPhoneNumberRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	ctx := context.TODO()

	log.Debug("Getting SID")

	sid := d.Id()

	log.Debug("Getting phone_number")

	phoneNumber := d.Get("number").(string)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("START client.IncomingNumbers.Get")

	ph, err := client.IncomingNumbers.Get(ctx, sid)

	if err != nil {
		return fmt.Errorf("Encountered an error when getting phone number SID %s: %s", sid, err)
	}

	err = mapTwilioPhoneNumberToTerraform(ph, d)

	if err != nil {
		return fmt.Errorf("Encountered an error while mapping Twilio API result to terraform: %s", err)
	}
	return err
}

func resourceTwilioPhoneNumberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	ctx := context.TODO()

	sid := d.Id()

	updatePayload := makeCreateRequestPayload(d)

	//phoneNumber := d.Get("number").(string)
	//updatePayload.Set("PhoneNumber", e164Number)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"phone_sid":   sid,
		},
	).Debug("START client.IncomingNumbers.Update")

	_, err := client.IncomingNumbers.Update(ctx, sid, updatePayload)

	if err != nil {
		return fmt.Errorf("Failed to update phone number SID %s: %s", sid, err)
	}
	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"phone_sid":   sid,
		},
	).Debug("END client.IncomingNumbers.Update")

	if d.HasChange("service_sid") {
		before, after := d.GetChange("service_sid")
		serviceIdAfter := cast.ToString(after)
		serviceIdBefore := cast.ToString(before)
		if len(serviceIdBefore) > 0 {
			log.WithFields(
				log.Fields{
					"account_sid": config.AccountSID,
					"phone_sid":   sid,
					"service_sid": serviceIdBefore,
				},
			).Debug("START client.Message.Services.DeletePhoneNumber")

			err := client.Message.Services.DeletePhoneNumber(ctx, serviceIdBefore, sid)
			if err != nil {
				return fmt.Errorf("Encountered error removing phone number with SID %s from messaging service with SID %s: %s", sid, serviceIdBefore, err)
			}

			log.WithFields(
				log.Fields{
					"account_sid": config.AccountSID,
					"phone_sid":   sid,
					"service_sid": serviceIdBefore,
				},
			).Debug("END client.Message.Services.DeletePhoneNumber")
		}
		if len(serviceIdAfter) > 0 {
			log.WithFields(
				log.Fields{
					"account_sid": config.AccountSID,
					"phone_sid":   sid,
					"service_sid": serviceIdAfter,
				},
			).Debug("START client.Message.Services.CreatePhoneNumber")

			_, err := client.Message.Services.CreatePhoneNumber(ctx, serviceIdAfter, sid)
			if err != nil && !strings.Contains(err.Error(), "already in the Messaging Service") {
				return fmt.Errorf("Encountered error adding phone number with SID %s to messaging service with SID %s: %s", sid, serviceIdAfter, err)
			}

			log.WithFields(
				log.Fields{
					"account_sid": config.AccountSID,
					"phone_sid":   sid,
					"service_sid": serviceIdAfter,
				},
			).Debug("END client.Message.Services.CreatePhoneNumber")
		}
	}

	return nil
}

func resourceTwilioPhoneNumberDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	ctx := context.TODO()

	sid := d.Id()
	phoneNumber := d.Get("number").(string)
	serviceId := cast.ToString(d.Get("service_id"))

	if len(serviceId) > 0 {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
				"phone_sid":   sid,
				"service_sid": serviceId,
			},
		).Debug("START client.Message.Services.DeletePhoneNumber")
		err := client.Message.Services.DeletePhoneNumber(ctx, serviceId, sid)
		if err != nil {
			return fmt.Errorf("Encountered error removing phone number with SID %s from messaging service with SID %s: %s", sid, serviceId, err)
		}
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
				"phone_sid":   sid,
				"service_sid": serviceId,
			},
		).Debug("END client.Message.Services.DeletePhoneNumber")
	}

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("START client.IncomingNumbers.Release")

	err := client.IncomingNumbers.Release(ctx, sid)

	log.WithFields(
		log.Fields{
			"account_sid":      config.AccountSID,
			"phone_number":     phoneNumber,
			"phone_number_sid": sid,
		},
	).Debug("END client.IncomingNumbers.Release")

	if err != nil {
		return fmt.Errorf("Failed to delete/release number: %s", err.Error())
	}

	return nil
}
