package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// NotificationIntegration represents a notification integration available in the organization
type NotificationIntegration struct {
	ID        string `json:"id"`         // Notification integration Global ID
	Name      string `json:"name"`       // Human-readable name
	Kind      string `json:"kind"`       // DISCORD_WEBHOOK or SLACK_ALERT
	CreatedAt string `json:"created_at"` // ISO 8601 timestamp
}

// ListNotificationIntegrationsResponse is the response from listing notification integrations
type ListNotificationIntegrationsResponse struct {
	NotificationIntegrations []NotificationIntegration `json:"notification_integrations"`
}

// ListNotificationIntegrations lists available notification integrations for the organization
// kind can be nil to list all, or one of: "DISCORD_WEBHOOK", "SLACK_ALERT"
func (c *Client) ListNotificationIntegrations(ctx context.Context, kind *string) ([]NotificationIntegration, error) {
	path := "/api/v1/notification-integrations"
	if kind != nil && *kind != "" {
		path = fmt.Sprintf("%s?kind=%s", path, *kind)
	}

	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification integrations: %w", err)
	}

	var response ListNotificationIntegrationsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.NotificationIntegrations, nil
}
