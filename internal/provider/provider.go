package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tierzero/terraform-provider-tierzero/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &TierZeroProvider{}
)

// New creates a new TierZero provider instance
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TierZeroProvider{
			version: version,
		}
	}
}

// TierZeroProvider is the provider implementation
type TierZeroProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance tests.
	version string
}

// TierZeroProviderModel describes the provider data model
type TierZeroProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
}

// Metadata returns the provider type name
func (p *TierZeroProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tierzero"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data
func (p *TierZeroProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing TierZero Alert Responders and related resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "TierZero Organization API Key. Can also be set via TIERZERO_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "TierZero API base URL. Defaults to https://api.tierzero.com",
				Optional:    true,
			},
		},
	}
}

// Configure prepares the provider for use
func (p *TierZeroProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config TierZeroProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get API key from config or environment variable
	apiKey := os.Getenv("TIERZERO_API_KEY")
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The provider requires an API key. Set the api_key attribute in the provider configuration or use the TIERZERO_API_KEY environment variable.",
		)
		return
	}

	// Get base URL from config or use default
	baseURL := "https://api.tierzero.com"
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	// Create client and make it available to resources and data sources
	apiClient := client.NewClient(baseURL, apiKey)
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

// DataSources defines the data sources implemented in the provider
func (p *TierZeroProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWebhookSubscriptionsDataSource,
		NewNotificationIntegrationsDataSource,
	}
}

// Resources defines the resources implemented in the provider
func (p *TierZeroProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAlertResponderResource,
	}
}
