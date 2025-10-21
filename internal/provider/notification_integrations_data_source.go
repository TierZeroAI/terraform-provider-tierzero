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
	_ datasource.DataSource              = &notificationIntegrationsDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationIntegrationsDataSource{}
)

// NewNotificationIntegrationsDataSource is a helper function to simplify the provider implementation.
func NewNotificationIntegrationsDataSource() datasource.DataSource {
	return &notificationIntegrationsDataSource{}
}

// notificationIntegrationsDataSource is the data source implementation.
type notificationIntegrationsDataSource struct {
	client *client.Client
}

// notificationIntegrationsDataSourceModel maps the data source schema data.
type notificationIntegrationsDataSourceModel struct {
	Kind                     types.String                     `tfsdk:"kind"`
	NotificationIntegrations []notificationIntegrationModel   `tfsdk:"notification_integrations"`
}

type notificationIntegrationModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Kind      types.String `tfsdk:"kind"`
	CreatedAt types.String `tfsdk:"created_at"`
}

// Metadata returns the data source type name.
func (d *notificationIntegrationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_integrations"
}

// Schema defines the schema for the data source.
func (d *notificationIntegrationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches available notification integrations for the organization. Use this to discover valid notification integration IDs when creating alert responders.",
		Attributes: map[string]schema.Attribute{
			"kind": schema.StringAttribute{
				Description: "Optional filter by integration kind (DISCORD_WEBHOOK or SLACK_ALERT)",
				Optional:    true,
			},
			"notification_integrations": schema.ListNestedAttribute{
				Description: "List of available notification integrations",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Notification integration Global ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Human-readable name",
							Computed:    true,
						},
						"kind": schema.StringAttribute{
							Description: "Integration kind (DISCORD_WEBHOOK or SLACK_ALERT)",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Creation timestamp (ISO 8601)",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *notificationIntegrationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *notificationIntegrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config notificationIntegrationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get kind filter if specified
	var kind *string
	if !config.Kind.IsNull() && config.Kind.ValueString() != "" {
		k := config.Kind.ValueString()
		kind = &k
	}

	// Fetch notification integrations from API
	integrations, err := d.client.ListNotificationIntegrations(ctx, kind)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Notification Integrations",
			"Could not read notification integrations: "+err.Error(),
		)
		return
	}

	// Map response to state
	config.NotificationIntegrations = make([]notificationIntegrationModel, len(integrations))
	for i, integration := range integrations {
		config.NotificationIntegrations[i] = notificationIntegrationModel{
			ID:        types.StringValue(integration.ID),
			Name:      types.StringValue(integration.Name),
			Kind:      types.StringValue(integration.Kind),
			CreatedAt: types.StringValue(integration.CreatedAt),
		}
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
