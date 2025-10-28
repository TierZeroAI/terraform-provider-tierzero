package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AlertResponder represents an alert responder resource
// Note: The Terraform schema uses 'enabled' (bool) instead of 'status' (string)
type AlertResponder struct {
	ID                         string                 `json:"id"`
	OrganizationName           string                 `json:"organization_name,omitempty"`
	TeamName                   string                 `json:"team_name"`
	Name                       string                 `json:"name"`
	Runbook                    *Runbook               `json:"runbook,omitempty"`
	MatchingCriteria           *MatchingCriteria      `json:"matching_criteria"`
	WebhookSources             []WebhookSource        `json:"webhook_sources,omitempty"`
	SlackChannelID             *string                `json:"slack_channel_id,omitempty"`
	NotificationIntegrationIDs []string               `json:"notification_integration_ids,omitempty"`
	Status                     string                 `json:"status"`           // API field: "ACTIVE" or "PAUSED" (not exposed in Terraform schema)
	CreatedAt                  string                 `json:"created_at,omitempty"`
	UpdatedAt                  string                 `json:"updated_at,omitempty"`
	URL                        string                 `json:"url,omitempty"`    // Returned by: Create, Update, List (not by Get, Enable, Disable)
}

// Runbook contains investigation prompts
type Runbook struct {
	InvestigationPrompt      string `json:"investigation_prompt,omitempty"`
	ImpactAndSeverityPrompt  string `json:"impact_and_severity_prompt,omitempty"`
}

// MatchingCriteria defines how alerts are matched
type MatchingCriteria struct {
	TextMatches         []string `json:"text_matches"`
	SlackBotAppUserID   *string  `json:"slack_bot_app_user_id,omitempty"`
}

// WebhookSource represents a webhook configuration
type WebhookSource struct {
	Type     string `json:"type"`     // PAGERDUTY, OPSGENIE, FIREHYDRANT, ROOTLY, SLACK
	RemoteID string `json:"remote_id"`
}

// CreateAlertResponderRequest is the request body for creating an alert responder
type CreateAlertResponderRequest struct {
	TeamName                   string                `json:"team_name"`
	Name                       string                `json:"name"`
	WebhookSources             []WebhookSource       `json:"webhook_sources,omitempty"`
	SlackChannelID             *string               `json:"slack_channel_id,omitempty"`
	MatchingCriteria           *MatchingCriteria     `json:"matching_criteria"`
	Runbook                    *Runbook              `json:"runbook,omitempty"`
	NotificationIntegrationIDs []string              `json:"notification_integration_ids,omitempty"`
}

// UpdateAlertResponderRequest is the request body for updating an alert responder
type UpdateAlertResponderRequest struct {
	Name                       *string               `json:"name,omitempty"`
	MatchingCriteria           *MatchingCriteria     `json:"matching_criteria,omitempty"`
	WebhookSources             []WebhookSource       `json:"webhook_sources,omitempty"`
	SlackChannelID             *string               `json:"slack_channel_id,omitempty"`
	Runbook                    *Runbook              `json:"runbook,omitempty"`
	NotificationIntegrationIDs []string              `json:"notification_integration_ids,omitempty"`
}

// ListAlertRespondersResponse is the response from listing alert responders
type ListAlertRespondersResponse struct {
	AlertResponders []AlertResponder `json:"alert_responders"`
}

// CreateAlertResponder creates a new alert responder
func (c *Client) CreateAlertResponder(ctx context.Context, req *CreateAlertResponderRequest) (*AlertResponder, error) {
	respBody, err := c.doRequest(ctx, http.MethodPost, "/api/v1/alert-responders", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert responder: %w", err)
	}

	var alertResponder AlertResponder
	if err := json.Unmarshal(respBody, &alertResponder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &alertResponder, nil
}

// GetAlertResponder retrieves an alert responder by ID
func (c *Client) GetAlertResponder(ctx context.Context, id string) (*AlertResponder, error) {
	path := fmt.Sprintf("/api/v1/alert-responders/%s", id)
	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert responder: %w", err)
	}

	var alertResponder AlertResponder
	if err := json.Unmarshal(respBody, &alertResponder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &alertResponder, nil
}

// ListAlertResponders lists all alert responders, optionally filtered by team
func (c *Client) ListAlertResponders(ctx context.Context, teamName *string) ([]AlertResponder, error) {
	path := "/api/v1/alert-responders"
	if teamName != nil && *teamName != "" {
		path = fmt.Sprintf("%s?team_name=%s", path, *teamName)
	}

	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list alert responders: %w", err)
	}

	var response ListAlertRespondersResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.AlertResponders, nil
}

// UpdateAlertResponder updates an existing alert responder
func (c *Client) UpdateAlertResponder(ctx context.Context, id string, req *UpdateAlertResponderRequest) (*AlertResponder, error) {
	path := fmt.Sprintf("/api/v1/alert-responders/%s", id)
	respBody, err := c.doRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update alert responder: %w", err)
	}

	var alertResponder AlertResponder
	if err := json.Unmarshal(respBody, &alertResponder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &alertResponder, nil
}

// DeleteAlertResponder deletes an alert responder (soft delete)
func (c *Client) DeleteAlertResponder(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/v1/alert-responders/%s", id)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete alert responder: %w", err)
	}

	return nil
}

// EnableAlertResponder enables an alert responder
func (c *Client) EnableAlertResponder(ctx context.Context, id string) (*AlertResponder, error) {
	path := fmt.Sprintf("/api/v1/alert-responders/%s/enable", id)
	respBody, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to enable alert responder: %w", err)
	}

	var alertResponder AlertResponder
	if err := json.Unmarshal(respBody, &alertResponder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &alertResponder, nil
}

// DisableAlertResponder disables an alert responder
func (c *Client) DisableAlertResponder(ctx context.Context, id string) (*AlertResponder, error) {
	path := fmt.Sprintf("/api/v1/alert-responders/%s/disable", id)
	respBody, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to disable alert responder: %w", err)
	}

	var alertResponder AlertResponder
	if err := json.Unmarshal(respBody, &alertResponder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &alertResponder, nil
}
