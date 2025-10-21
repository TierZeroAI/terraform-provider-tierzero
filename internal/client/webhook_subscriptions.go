package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// WebhookSubscription represents a webhook subscription available in the organization
type WebhookSubscription struct {
	Type     string `json:"type"`      // PAGERDUTY, OPSGENIE, FIREHYDRANT, ROOTLY, SLACK
	RemoteID string `json:"remote_id"` // External webhook ID
	Name     string `json:"name"`      // Human-readable name
}

// ListWebhookSubscriptionsResponse is the response from listing webhook subscriptions
type ListWebhookSubscriptionsResponse struct {
	WebhookSubscriptions []WebhookSubscription `json:"webhook_subscriptions"`
}

// ListWebhookSubscriptions lists available webhook subscriptions for the organization
func (c *Client) ListWebhookSubscriptions(ctx context.Context) ([]WebhookSubscription, error) {
	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/webhook-subscriptions", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhook subscriptions: %w", err)
	}

	var response ListWebhookSubscriptionsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.WebhookSubscriptions, nil
}
