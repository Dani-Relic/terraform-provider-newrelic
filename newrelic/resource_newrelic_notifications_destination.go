package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicNotificationDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNotificationDestinationCreate,
		ReadContext:   resourceNewRelicNotificationDestinationRead,
		DeleteContext: resourceNewRelicNotificationDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "(Required) The name of the destination.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidNotificationsDestinationTypes(), false),
				Description:  fmt.Sprintf("(Required) The type of the destination. One of: (%s).", strings.Join(listValidNotificationsDestinationTypes(), ", ")),
			},
			"properties": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of notification destination property types.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notification destination property key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notification destination property value.",
						},
						"label": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification destination property label.",
						},
						"display_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification destination property display key.",
						},
					},
				},
			},
			"auth": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "A set of key-value pairs to represent a Notification destination auth.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Sensitive:   true,
			},
		},
	}
}

func resourceNewRelicNotificationDestinationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	destinationInput, err := expandNotificationDestination(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic notification destinationResponse %s", destinationInput.Name)

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	destinationResponse, err := client.Notifications.AiNotificationsCreateDestinationWithContext(updatedContext, accountID, *destinationInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(destinationResponse.Destination.ID)

	return resourceNewRelicNotificationDestinationRead(updatedContext, d, meta)
}

func resourceNewRelicNotificationDestinationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic notification destinationResponse %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	filters := ai.AiNotificationsDestinationFilter{ID: d.Id()}
	sorter := notifications.AiNotificationsDestinationSorter{}
	updatedContext := updateContextWithAccountID(ctx, accountID)

	destinationResponse, err := client.Notifications.GetDestinationsWithContext(updatedContext, accountID, "", filters, sorter)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenNotificationDestination(&destinationResponse.Entities[0], d))
}

func resourceNewRelicNotificationDestinationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic notification destination %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	if _, err := client.Notifications.AiNotificationsDeleteDestinationWithContext(updatedContext, accountID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Validation function to validate allowed destination types
func listValidNotificationsDestinationTypes() []string {
	return []string{
		string(notifications.AiNotificationsDestinationTypeTypes.WEBHOOK),
		string(notifications.AiNotificationsDestinationTypeTypes.EMAIL),
		string(notifications.AiNotificationsDestinationTypeTypes.SERVICE_NOW),
		string(notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_ACCOUNT_INTEGRATION),
		string(notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_SERVICE_INTEGRATION),
	}
}
