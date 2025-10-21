package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/tierzero/terraform-provider-tierzero/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &alertResponderResource{}
	_ resource.ResourceWithConfigure   = &alertResponderResource{}
	_ resource.ResourceWithImportState = &alertResponderResource{}
)

// NewAlertResponderResource is a helper function to simplify the provider implementation.
func NewAlertResponderResource() resource.Resource {
	return &alertResponderResource{}
}

// alertResponderResource is the resource implementation.
type alertResponderResource struct {
	client *client.Client
}

// alertResponderResourceModel maps the resource schema data.
type alertResponderResourceModel struct {
	ID                         types.String                   `tfsdk:"id"`
	TeamName                   types.String                   `tfsdk:"team_name"`
	Name                       types.String                   `tfsdk:"name"`
	WebhookSources             []webhookSourceModel           `tfsdk:"webhook_sources"`
	SlackChannelID             types.String                   `tfsdk:"slack_channel_id"`
	MatchingCriteria           *matchingCriteriaModel         `tfsdk:"matching_criteria"`
	Runbook                    *runbookModel                  `tfsdk:"runbook"`
	NotificationIntegrationIDs []types.String                 `tfsdk:"notification_integration_ids"`
	Enabled                    types.Bool                     `tfsdk:"enabled"`
	URL                        types.String                   `tfsdk:"url"`
	CreatedAt                  types.String                   `tfsdk:"created_at"`
	UpdatedAt                  types.String                   `tfsdk:"updated_at"`
}

type webhookSourceModel struct {
	Type     types.String `tfsdk:"type"`
	RemoteID types.String `tfsdk:"remote_id"`
}

type matchingCriteriaModel struct {
	TextMatches       []types.String `tfsdk:"text_matches"`
	SlackBotAppUserID types.String   `tfsdk:"slack_bot_app_user_id"`
}

type runbookModel struct {
	Prompt     types.String `tfsdk:"prompt"`
	FastPrompt types.String `tfsdk:"fast_prompt"`
}

// Metadata returns the resource type name.
func (r *alertResponderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_responder"
}

// Schema defines the schema for the resource.
func (r *alertResponderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a TierZero Alert Responder that automatically investigates incoming alerts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Alert Responder Global ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"team_name": schema.StringAttribute{
				Description: "Team name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Alert responder name",
				Required:    true,
			},
			"webhook_sources": schema.ListNestedAttribute{
				Description: "Webhook sources to monitor (for PagerDuty, OpsGenie, FireHydrant, Rootly). Mutually exclusive with slack_channel_id.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Webhook type (PAGERDUTY, OPSGENIE, FIREHYDRANT, ROOTLY)",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("PAGERDUTY", "OPSGENIE", "FIREHYDRANT", "ROOTLY"),
							},
						},
						"remote_id": schema.StringAttribute{
							Description: "External webhook ID",
							Required:    true,
						},
					},
				},
			},
			"slack_channel_id": schema.StringAttribute{
				Description: "Slack channel ID (e.g., 'C01234567' for public channels, 'G01234567' for private channels). Mutually exclusive with webhook_sources.",
				Optional:    true,
			},
			"matching_criteria": schema.SingleNestedAttribute{
				Description: "Criteria for matching alerts",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"text_matches": schema.ListAttribute{
						Description: "Array of text patterns to match",
						Required:    true,
						ElementType: types.StringType,
					},
					"slack_bot_app_user_id": schema.StringAttribute{
						Description: "Optional Slack bot/sender app user ID to filter messages (only for Slack alerts)",
						Optional:    true,
					},
				},
			},
			"runbook": schema.SingleNestedAttribute{
				Description: "Investigation runbook (optional, uses default if not provided)",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"prompt": schema.StringAttribute{
						Description: "Main investigation prompt",
						Optional:    true,
					},
					"fast_prompt": schema.StringAttribute{
						Description: "Quick triage prompt",
						Optional:    true,
					},
				},
			},
			"notification_integration_ids": schema.ListAttribute{
				Description: "Notification integration Global IDs",
				Optional:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the alert responder is enabled. When true, status is ACTIVE. When false, status is PAUSED. Uses enable/disable API endpoints under the hood.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"url": schema.StringAttribute{
				Description: "Link to alert responder details page (returned by create/update operations)",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp (ISO 8601)",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp (ISO 8601)",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *alertResponderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *alertResponderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alertResponderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either webhook_sources OR slack_channel_id is provided (not both, not neither)
	hasWebhookSources := len(plan.WebhookSources) > 0
	hasSlackChannelID := !plan.SlackChannelID.IsNull() && plan.SlackChannelID.ValueString() != ""

	if !hasWebhookSources && !hasSlackChannelID {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Must specify either webhook_sources or slack_channel_id",
		)
		return
	}

	if hasWebhookSources && hasSlackChannelID {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Cannot specify both webhook_sources and slack_channel_id. These are mutually exclusive.",
		)
		return
	}

	// Build the create request
	createReq := &client.CreateAlertResponderRequest{
		TeamName:         plan.TeamName.ValueString(),
		Name:             plan.Name.ValueString(),
		MatchingCriteria: buildMatchingCriteria(plan.MatchingCriteria),
	}

	// Set webhook_sources or slack_channel_id
	if hasWebhookSources {
		createReq.WebhookSources = buildWebhookSources(plan.WebhookSources)
	}
	if hasSlackChannelID {
		slackChannelID := plan.SlackChannelID.ValueString()
		createReq.SlackChannelID = &slackChannelID
	}

	if plan.Runbook != nil {
		createReq.Runbook = buildRunbook(plan.Runbook)
	}

	if len(plan.NotificationIntegrationIDs) > 0 {
		createReq.NotificationIntegrationIDs = buildStringList(plan.NotificationIntegrationIDs)
	}

	// Create the alert responder
	alertResponder, err := r.client.CreateAlertResponder(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Alert Responder",
			"Could not create alert responder: "+err.Error(),
		)
		return
	}

	// If enabled is false, disable the alert responder
	if !plan.Enabled.ValueBool() {
		_, err = r.client.DisableAlertResponder(ctx, alertResponder.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Disabling Alert Responder",
				"Alert responder was created but could not be disabled: "+err.Error(),
			)
			return
		}
		alertResponder.Status = "PAUSED"
	}

	// Read back to get full details
	fullAlertResponder, err := r.client.GetAlertResponder(ctx, alertResponder.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Alert Responder",
			"Could not read alert responder after creation: "+err.Error(),
		)
		return
	}

	// Preserve URL from create response (GET doesn't return it)
	fullAlertResponder.URL = alertResponder.URL

	// Map response to state
	plan.ID = types.StringValue(fullAlertResponder.ID)
	plan.URL = types.StringValue(fullAlertResponder.URL)
	plan.CreatedAt = types.StringValue(fullAlertResponder.CreatedAt)
	plan.UpdatedAt = types.StringValue(fullAlertResponder.UpdatedAt)
	plan.Enabled = types.BoolValue(fullAlertResponder.Status == "ACTIVE")

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *alertResponderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alertResponderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current alert responder
	alertResponder, err := r.client.GetAlertResponder(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			// Alert responder was deleted outside Terraform
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Alert Responder",
			"Could not read alert responder: "+err.Error(),
		)
		return
	}

	// Update state from API response
	state.TeamName = types.StringValue(alertResponder.TeamName)
	state.Name = types.StringValue(alertResponder.Name)
	state.WebhookSources = mapWebhookSources(alertResponder.WebhookSources)

	// Handle slack_channel_id
	if alertResponder.SlackChannelID != nil && *alertResponder.SlackChannelID != "" {
		state.SlackChannelID = types.StringValue(*alertResponder.SlackChannelID)
	} else {
		state.SlackChannelID = types.StringNull()
	}

	state.MatchingCriteria = mapMatchingCriteria(alertResponder.MatchingCriteria)
	state.Runbook = mapRunbook(alertResponder.Runbook)
	state.NotificationIntegrationIDs = mapStringList(alertResponder.NotificationIntegrationIDs)
	state.Enabled = types.BoolValue(alertResponder.Status == "ACTIVE")
	state.CreatedAt = types.StringValue(alertResponder.CreatedAt)
	state.UpdatedAt = types.StringValue(alertResponder.UpdatedAt)

	// Preserve URL if not returned by API (GET doesn't return it)
	if alertResponder.URL != "" {
		state.URL = types.StringValue(alertResponder.URL)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alertResponderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alertResponderResourceModel
	var state alertResponderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	// Validate that either webhook_sources OR slack_channel_id is provided (not both, not neither)
	hasWebhookSources := len(plan.WebhookSources) > 0
	hasSlackChannelID := !plan.SlackChannelID.IsNull() && plan.SlackChannelID.ValueString() != ""

	if !hasWebhookSources && !hasSlackChannelID {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Must specify either webhook_sources or slack_channel_id",
		)
		return
	}

	if hasWebhookSources && hasSlackChannelID {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Cannot specify both webhook_sources and slack_channel_id. These are mutually exclusive.",
		)
		return
	}

	// Handle enabled field changes first
	if !plan.Enabled.Equal(state.Enabled) {
		if plan.Enabled.ValueBool() {
			_, err := r.client.EnableAlertResponder(ctx, id)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Enabling Alert Responder",
					"Could not enable alert responder: "+err.Error(),
				)
				return
			}
		} else {
			_, err := r.client.DisableAlertResponder(ctx, id)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Disabling Alert Responder",
					"Could not disable alert responder: "+err.Error(),
				)
				return
			}
		}
	}

	// Check if other fields changed
	// Note: team_name is not included because it has RequiresReplace() plan modifier
	needsUpdate := !plan.Name.Equal(state.Name) ||
		webhookSourcesChanged(plan.WebhookSources, state.WebhookSources) ||
		!plan.SlackChannelID.Equal(state.SlackChannelID) ||
		matchingCriteriaChanged(plan.MatchingCriteria, state.MatchingCriteria) ||
		runbookChanged(plan.Runbook, state.Runbook) ||
		notificationIDsChanged(plan.NotificationIntegrationIDs, state.NotificationIntegrationIDs)

	if needsUpdate {
		// Build update request
		updateReq := &client.UpdateAlertResponderRequest{}

		if !plan.Name.Equal(state.Name) {
			name := plan.Name.ValueString()
			updateReq.Name = &name
		}

		if webhookSourcesChanged(plan.WebhookSources, state.WebhookSources) {
			updateReq.WebhookSources = buildWebhookSources(plan.WebhookSources)
		}

		if !plan.SlackChannelID.Equal(state.SlackChannelID) {
			if !plan.SlackChannelID.IsNull() && plan.SlackChannelID.ValueString() != "" {
				slackChannelID := plan.SlackChannelID.ValueString()
				updateReq.SlackChannelID = &slackChannelID
			}
		}

		if matchingCriteriaChanged(plan.MatchingCriteria, state.MatchingCriteria) {
			updateReq.MatchingCriteria = buildMatchingCriteria(plan.MatchingCriteria)
		}

		if runbookChanged(plan.Runbook, state.Runbook) {
			updateReq.Runbook = buildRunbook(plan.Runbook)
		}

		if notificationIDsChanged(plan.NotificationIntegrationIDs, state.NotificationIntegrationIDs) {
			updateReq.NotificationIntegrationIDs = buildStringList(plan.NotificationIntegrationIDs)
		}

		// Update the alert responder
		alertResponder, err := r.client.UpdateAlertResponder(ctx, id, updateReq)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Alert Responder",
				"Could not update alert responder: "+err.Error(),
			)
			return
		}

		// Update URL if returned
		if alertResponder.URL != "" {
			plan.URL = types.StringValue(alertResponder.URL)
		}
	}

	// Read back to get full details
	fullAlertResponder, err := r.client.GetAlertResponder(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Alert Responder",
			"Could not read alert responder after update: "+err.Error(),
		)
		return
	}

	// Update state
	plan.UpdatedAt = types.StringValue(fullAlertResponder.UpdatedAt)
	plan.Enabled = types.BoolValue(fullAlertResponder.Status == "ACTIVE")

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alertResponderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alertResponderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the alert responder
	err := r.client.DeleteAlertResponder(ctx, state.ID.ValueString())
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError(
				"Error Deleting Alert Responder",
				"Could not delete alert responder: "+err.Error(),
			)
			return
		}
	}
}

// ImportState imports the resource into Terraform state.
func (r *alertResponderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID provided by the user
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper functions to build client types from Terraform models

func buildWebhookSources(sources []webhookSourceModel) []client.WebhookSource {
	result := make([]client.WebhookSource, len(sources))
	for i, s := range sources {
		result[i] = client.WebhookSource{
			Type:     s.Type.ValueString(),
			RemoteID: s.RemoteID.ValueString(),
		}
	}
	return result
}

func buildMatchingCriteria(mc *matchingCriteriaModel) *client.MatchingCriteria {
	if mc == nil {
		return nil
	}
	result := &client.MatchingCriteria{
		TextMatches: buildStringList(mc.TextMatches),
	}
	if !mc.SlackBotAppUserID.IsNull() && mc.SlackBotAppUserID.ValueString() != "" {
		slackBotAppUserID := mc.SlackBotAppUserID.ValueString()
		result.SlackBotAppUserID = &slackBotAppUserID
	}
	return result
}

func buildRunbook(rb *runbookModel) *client.Runbook {
	if rb == nil {
		return nil
	}
	return &client.Runbook{
		Prompt:     rb.Prompt.ValueString(),
		FastPrompt: rb.FastPrompt.ValueString(),
	}
}

func buildStringList(list []types.String) []string {
	result := make([]string, len(list))
	for i, s := range list {
		result[i] = s.ValueString()
	}
	return result
}

// Helper functions to map client types to Terraform models

func mapWebhookSources(sources []client.WebhookSource) []webhookSourceModel {
	result := make([]webhookSourceModel, len(sources))
	for i, s := range sources {
		result[i] = webhookSourceModel{
			Type:     types.StringValue(s.Type),
			RemoteID: types.StringValue(s.RemoteID),
		}
	}
	return result
}

func mapMatchingCriteria(mc *client.MatchingCriteria) *matchingCriteriaModel {
	if mc == nil {
		return nil
	}
	result := &matchingCriteriaModel{
		TextMatches: mapStringList(mc.TextMatches),
	}
	if mc.SlackBotAppUserID != nil && *mc.SlackBotAppUserID != "" {
		result.SlackBotAppUserID = types.StringValue(*mc.SlackBotAppUserID)
	} else {
		result.SlackBotAppUserID = types.StringNull()
	}
	return result
}

func mapRunbook(rb *client.Runbook) *runbookModel {
	if rb == nil {
		return nil
	}
	return &runbookModel{
		Prompt:     types.StringValue(rb.Prompt),
		FastPrompt: types.StringValue(rb.FastPrompt),
	}
}

func mapStringList(list []string) []types.String {
	result := make([]types.String, len(list))
	for i, s := range list {
		result[i] = types.StringValue(s)
	}
	return result
}

// Helper functions to detect changes

func webhookSourcesChanged(plan, state []webhookSourceModel) bool {
	if len(plan) != len(state) {
		return true
	}
	for i := range plan {
		if !plan[i].Type.Equal(state[i].Type) || !plan[i].RemoteID.Equal(state[i].RemoteID) {
			return true
		}
	}
	return false
}

func matchingCriteriaChanged(plan, state *matchingCriteriaModel) bool {
	if (plan == nil) != (state == nil) {
		return true
	}
	if plan == nil {
		return false
	}
	if len(plan.TextMatches) != len(state.TextMatches) {
		return true
	}
	for i := range plan.TextMatches {
		if !plan.TextMatches[i].Equal(state.TextMatches[i]) {
			return true
		}
	}
	return false
}

func runbookChanged(plan, state *runbookModel) bool {
	if (plan == nil) != (state == nil) {
		return true
	}
	if plan == nil {
		return false
	}
	return !plan.Prompt.Equal(state.Prompt) || !plan.FastPrompt.Equal(state.FastPrompt)
}

func notificationIDsChanged(plan, state []types.String) bool {
	if len(plan) != len(state) {
		return true
	}
	for i := range plan {
		if !plan[i].Equal(state[i]) {
			return true
		}
	}
	return false
}
