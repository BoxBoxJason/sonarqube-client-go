package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebhooks_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "webhooks/create") {
			t.Errorf("expected path to contain webhooks/create, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := WebhooksCreate{
			Webhook: Webhook{
				Key:       "uuid-webhook-1",
				Name:      "My Webhook",
				URL:       "https://example.com/webhook",
				HasSecret: true,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &WebhooksCreateOption{
		Name:   "My Webhook",
		URL:    "https://example.com/webhook",
		Secret: "my-secret-at-least-16",
	}

	result, resp, err := client.Webhooks.Create(opt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Webhook.Key != "uuid-webhook-1" {
		t.Errorf("expected key uuid-webhook-1, got %s", result.Webhook.Key)
	}
}

func TestWebhooks_Create_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Webhooks.Create(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Name
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		URL: "https://example.com",
	})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Test missing URL
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		Name: "My Webhook",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}

	// Test Name too long
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		Name: strings.Repeat("a", MaxWebhookNameLength+1),
		URL:  "https://example.com",
	})
	if err == nil {
		t.Error("expected error for Name exceeding max length")
	}

	// Test URL too long
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		Name: "My Webhook",
		URL:  strings.Repeat("a", MaxWebhookURLLength+1),
	})
	if err == nil {
		t.Error("expected error for URL exceeding max length")
	}

	// Test Secret too short
	_, _, err = client.Webhooks.Create(&WebhooksCreateOption{
		Name:   "My Webhook",
		URL:    "https://example.com",
		Secret: "short",
	})
	if err == nil {
		t.Error("expected error for Secret below min length")
	}
}

func TestWebhooks_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &WebhooksDeleteOption{
		Webhook: "my-webhook-key",
	}

	resp, err := client.Webhooks.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestWebhooks_Delete_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Webhooks.Delete(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Webhook
	_, err = client.Webhooks.Delete(&WebhooksDeleteOption{})
	if err == nil {
		t.Error("expected error for missing Webhook")
	}
}

func TestWebhooks_Deliveries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &WebhooksDeliveriesOption{
		Webhook: "webhook-key",
	}

	result, resp, err := client.Webhooks.Deliveries(opt)
	if err != nil {
		t.Fatalf("Deliveries failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Deliveries) != 1 {
		t.Errorf("expected 1 delivery, got %d", len(result.Deliveries))
	}
}

func TestWebhooks_Deliveries_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Webhooks.Deliveries(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}
}

func TestWebhooks_Delivery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &WebhooksDeliveryOption{
		DeliveryID: "delivery-1",
	}

	result, resp, err := client.Webhooks.Delivery(opt)
	if err != nil {
		t.Fatalf("Delivery failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Delivery.Payload == "" {
		t.Error("expected payload to be set")
	}
}

func TestWebhooks_Delivery_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Webhooks.Delivery(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing DeliveryID
	_, _, err = client.Webhooks.Delivery(&WebhooksDeliveryOption{})
	if err == nil {
		t.Error("expected error for missing DeliveryID")
	}
}

func TestWebhooks_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &WebhooksListOption{}

	result, resp, err := client.Webhooks.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Webhooks) != 1 {
		t.Errorf("expected 1 webhook, got %d", len(result.Webhooks))
	}
}

func TestWebhooks_List_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Webhooks.List(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}
}

func TestWebhooks_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &WebhooksUpdateOption{
		Webhook: "webhook-1",
		Name:    "Updated Webhook",
		URL:     "https://example.com/updated",
	}

	resp, err := client.Webhooks.Update(opt)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestWebhooks_Update_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Webhooks.Update(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Name
	_, err = client.Webhooks.Update(&WebhooksUpdateOption{
		URL:     "https://example.com",
		Webhook: "webhook-1",
	})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Test missing URL
	_, err = client.Webhooks.Update(&WebhooksUpdateOption{
		Name:    "My Webhook",
		Webhook: "webhook-1",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}

	// Test missing Webhook
	_, err = client.Webhooks.Update(&WebhooksUpdateOption{
		Name: "My Webhook",
		URL:  "https://example.com",
	})
	if err == nil {
		t.Error("expected error for missing Webhook")
	}
}
