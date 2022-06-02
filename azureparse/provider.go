package azureparse

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type provider *schema.Provider

func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: providerResources(),

		Schema: map[string]*schema.Schema{
			// provider schema values taken from terraform-provider-azurerm
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
				Description: "The Subscription ID which should be used.",
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
				Description: "The Client ID which should be used.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
				Description: "The Tenant ID which should be used.",
			},

			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
				Description: "azure environment",
			},
		},
	}

	p.ConfigureContextFunc = initProvider(p)

	return p
}

func initProvider(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		builder := &authentication.Builder{
			SubscriptionID: d.Get("subscription_id").(string),
			ClientID:       d.Get("client_id").(string),
			ClientSecret:   d.Get("client_secret").(string),
			TenantID:       d.Get("tenant_id").(string),
			Environment:    d.Get("environment").(string),

			SupportsAzureCliToken:    true,
			SupportsClientSecretAuth: true,
		}

		config, err := builder.Build()
		if err != nil {
			return nil, diag.FromErr(err)
		}

		client, err := buildClient(builder, config)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		client.StopContext = ctx
		return client, nil
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azureparse_resource_group": resourceGroup(),
	}
}
