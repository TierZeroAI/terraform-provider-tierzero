package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tierzero/terraform-provider-tierzero/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &webhookSubscriptionsDataSource{}
	_ datasource.DataSourceWithConfigure = &webhookSubscriptionsDataSource{}
)

// NewWebhookSubscriptionsDataSource is a helper function to simplify the provider implementation.
func NewWebhookSubscriptionsDataSource() datasource.DataSource {
	return &webhookSubscriptionsDataSource{}
}

// webhookSubscriptionsDataSource is the data source implementation.
type webhookSubscriptionsDataSource struct {
	client *client.Client
}

// webhookSubscriptionsDataSourceModel maps the data source schema data.
type webhookSubscriptionsDataSourceModel struct {
	WebhookSubscriptions []webhookSubscriptionModel `tfsdk:"webhook_subscriptions"`
}

type webhookSubscriptionModel struct {
	Type     types.String `tfsdk:"type"`
	RemoteID types.String `tfsdk:"remote_id"`
	Name     types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *webhookSubscriptionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_subscriptions"
}

// Schema defines the schema for the data source.
func (d *webhookSubscriptionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches available webhook subscriptions for the organization. Use this to discover valid webhook sources when creating alert responders.",
		Attributes: map[string]schema.Attribute{
			"webhook_subscriptions": schema.ListNestedAttribute{
				Description: "List of available webhook subscriptions",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Webhook type (PAGERDUTY, OPSGENIE, FIREHYDRANT, ROOTLY, SLACK)",
							Computed:    true,
						},
						"remote_id": schema.StringAttribute{
							Description: "External webhook ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Human-readable name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *webhookSubscriptionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *webhookSubscriptionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state webhookSubscriptionsDataSourceModel

	// Fetch webhook subscriptions from API
	subscriptions, err := d.client.ListWebhookSubscriptions(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Webhook Subscriptions",
			"Could not read webhook subscriptions: "+err.Error(),
		)
		return
	}

	// Map response to state
	state.WebhookSubscriptions = make([]webhookSubscriptionModel, len(subscriptions))
	for i, sub := range subscriptions {
		state.WebhookSubscriptions[i] = webhookSubscriptionModel{
			Type:     types.StringValue(sub.Type),
			RemoteID: types.StringValue(sub.RemoteID),
			Name:     types.StringValue(sub.Name),
		}
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
