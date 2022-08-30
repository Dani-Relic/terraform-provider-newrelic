//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/notifications"

	"github.com/stretchr/testify/assert"
)

func TestExpandNotificationChannel(t *testing.T) {
	property := map[string]interface{}{
		"key":   "payload",
		"value": "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
		"label": "Payload Template",
	}

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Expanded     *notifications.AiNotificationsChannel
	}{
		"valid channel": {
			Data: map[string]interface{}{
				"name":           "testing123",
				"type":           "WEBHOOK",
				"properties":     []interface{}{property},
				"product":        "IINT",
				"destination_id": "b1e90a32-23b7-4028-b2c7-ffbdfe103852",
			},
			Expanded: &notifications.AiNotificationsChannel{
				Name: "testing123",
				Type: notifications.AiNotificationsChannelTypeTypes.WEBHOOK,
				Properties: []notifications.AiNotificationsProperty{
					{
						Key:   "payload",
						Value: "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
						Label: "Payload Template",
					},
				},
			},
		},
	}

	r := resourceNewRelicNotificationChannel()

	for _, tc := range cases {
		d := r.TestResourceData()

		for k, v := range tc.Data {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("err: %s", err)
			}
		}

		expanded, err := expandNotificationChannel(d)

		if tc.ExpectErr {
			assert.NotNil(t, err)
			assert.Equal(t, err.Error(), tc.ExpectReason)
		} else {
			assert.Nil(t, err)
		}

		if tc.Expanded != nil {
			assert.Equal(t, tc.Expanded.Name, expanded.Name)
		}
	}
}

func TestFlattenNotificationChannel(t *testing.T) {
	r := resourceNewRelicNotificationChannel()

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Flattened    *notifications.AiNotificationsChannel
	}{
		"minimal": {
			Data: map[string]interface{}{
				"name": "testing123",
				"type": "WEBHOOK",
			},
			Flattened: &notifications.AiNotificationsChannel{
				Name:          "testing123",
				Type:          "WEBHOOK",
				Product:       "IINT",
				DestinationId: "b1e90a32-23b7-4028-b2c7-ffbdfe103852",
			},
		},
	}

	for _, tc := range cases {
		if tc.Flattened != nil {
			d := r.TestResourceData()
			err := flattenNotificationChannel(tc.Flattened, d)
			assert.NoError(t, err)

			for k, v := range tc.Data {
				var x interface{}
				var ok bool
				if x, ok = d.GetOk(k); !ok {
					t.Fatalf("err: %s", err)
				}

				assert.Equal(t, x, v)
			}
		}
	}
}
