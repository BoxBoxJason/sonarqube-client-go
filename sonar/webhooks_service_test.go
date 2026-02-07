package sonar

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhooks_Create(t *testing.T) {
	response := WebhooksCreate{
		Webhook: Webhook{
			Key:       "uuid-webhook-1",
			Name:      "My Webhook",
			URL:       "https://example.com/webhook",
			HasSecret: true,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/webhooks/create", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &WebhooksCreateOption{
		Name:   "My Webhook",
		URL:    "https://example.com/webhook",
		Secret: "my-secret-at-least-16",
	}

	result, resp, err := client.Webhooks.Create(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "uuid-webhook-1", result.Webhook.Key)
}

func TestWebhooks_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Webhooks.Create(nil)
	assert.Error(t, err)

	// Test missing Name
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		URL: "https://example.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		Name: "My Webhook",
	})
	assert.Error(t, err)

	// Test Name too long
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		Name: strings.Repeat("a", MaxWebhookNameLength+1),
		URL:  "https://example.com",
	})
	assert.Error(t, err)

	// Test URL too long
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		Name: "My Webhook",
		URL:  strings.Repeat("a", MaxWebhookURLLength+1),
	})
	assert.Error(t, err)

	// Test Secret too short
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		Name:   "My Webhook",
		URL:    "https://example.com",
		Secret: "short",
	})
	assert.Error(t, err)
}

func TestWebhooks_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/webhooks/delete", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &WebhooksDeleteOption{
		Webhook: "my-webhook-key",
	}

	resp, err := client.Webhooks.Delete(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestWebhooks_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Webhooks.Delete(nil)
	assert.Error(t, err)

	// Test missing Webhook
	_, err = client.Webhooks.Delete(&WebhooksDeleteOption{})
	assert.Error(t, err)
}

func TestWebhooks_Deliveries(t *testing.T) {
	response := WebhooksDeliveries{
		Deliveries: []WebhookDelivery{
			{
				ID:         "delivery-1",
				At:         "2024-01-01T00:00:00+0000",
				Success:    true,
				HTTPStatus: 200,
				DurationMs: 150,
			},
		},
		Paging: Paging{
			PageIndex: 1,
			PageSize:  10,
			Total:     1,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/webhooks/deliveries", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &WebhooksDeliveriesOption{
		Webhook: "webhook-key",
	}

	result, resp, err := client.Webhooks.Deliveries(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Deliveries, 1)
}

func TestWebhooks_Deliveries_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Webhooks.Deliveries(nil)
	assert.Error(t, err)
}

func TestWebhooks_Delivery(t *testing.T) {
	response := WebhooksDelivery{
		Delivery: WebhookDelivery{
			ID:         "delivery-1",
			At:         "2024-01-01T00:00:00+0000",
			Success:    true,
			HTTPStatus: 200,
			DurationMs: 150,
			Payload:    `{"project":"my-project"}`,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/webhooks/delivery", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &WebhooksDeliveryOption{
		DeliveryID: "delivery-1",
	}

	result, resp, err := client.Webhooks.Delivery(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, result.Delivery.Payload)
}

func TestWebhooks_Delivery_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Webhooks.Delivery(nil)
	assert.Error(t, err)

	// Test missing DeliveryID
	_, _, err = client.Webhooks.Delivery(&WebhooksDeliveryOption{})
	assert.Error(t, err)
}

func TestWebhooks_List(t *testing.T) {
	response := WebhooksList{
		Webhooks: []Webhook{
			{
				Key:       "webhook-1",
				Name:      "Global Webhook",
				URL:       "https://example.com/webhook",
				HasSecret: false,
			},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/webhooks/list", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &WebhooksListOption{}

	result, resp, err := client.Webhooks.List(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Webhooks, 1)
}

func TestWebhooks_List_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Webhooks.List(nil)
	assert.Error(t, err)
}

func TestWebhooks_Update(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/webhooks/update", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &WebhooksUpdateOption{
		Webhook: "webhook-1",
		Name:    "Updated Webhook",
		URL:     "https://example.com/updated",
	}

	resp, err := client.Webhooks.Update(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestWebhooks_Update_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Webhooks.Update(nil)
	assert.Error(t, err)

	// Test missing Name
	_, err = client.Webhooks.Update(&WebhooksUpdateOption{
		URL:     "https://example.com",
		Webhook: "webhook-1",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.Webhooks.Update(&WebhooksUpdateOption{
		Name:    "My Webhook",
		Webhook: "webhook-1",
	})
	assert.Error(t, err)

	// Test missing Webhook
	_, err = client.Webhooks.Update(&WebhooksUpdateOption{
		Name: "My Webhook",
		URL:  "https://example.com",
	})
	assert.Error(t, err)
}
