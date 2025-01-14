package twilio

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var descriptions map[string]string

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:         providerSchema(),
		DataSourcesMap: providerDataSourcesMap(),
		ResourcesMap:   providerResources(),
		ConfigureFunc:  providerConfigure,
	}
}

// List of supported configuration fields for your provider.
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_sid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The unique ID that identifies your Twilio account. Starts with `AC` and can be found on the Settings -> General page (https://www.twilio.com/console/project/settings).",
		},
		"auth_token": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Your secret token to access your Twilio account. Keep this safe - DO NOT check this into source control!",
		},
		"endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Allows you to change the Twilio API endpoint. Nearly everyone will leave this blank; Twilions may find use of this setting, though!",
		},
	}
}

// List of supported resources and their configuration fields.
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"twilio_phone_number": resourceTwilioPhoneNumber(),
		"twilio_messaging_service": resourceTwilioMessagingService(),
		"twilio_subaccount":   resourceTwilioSubaccount(),
		"twilio_api_key":      resourceTwilioApiKey(),
	}
}

// List of supported data sources and their configuration fields.
func providerDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
	    "twilio_messaging_service": dataTwilioMessagingService(),
		"twilio_phone_number": dataTwilioPhoneNumber(),
		"twilio_subaccount":   dataTwilioSubaccount(),
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccountSID: d.Get("account_sid").(string),
		AuthToken:  d.Get("auth_token").(string),
		Endpoint:   d.Get("endpoint").(string),
	}
	return config.Client()
}

func makeComputed(s map[string]*schema.Schema) map[string]*schema.Schema {
	for _, p := range s {
		p.Optional = false
		p.Required = false
		p.Computed = true
		p.MaxItems = 0
		p.MinItems = 0
		p.ValidateFunc = nil
		p.DefaultFunc = nil
		p.Default = nil
		if resource, ok := p.Elem.(*schema.Resource); ok {
			makeComputed(resource.Schema)
		}
	}
	return s
}
