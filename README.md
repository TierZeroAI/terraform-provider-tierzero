# Terraform Provider for TierZero

The TierZero Terraform Provider allows you to manage TierZero Alert Responders and related resources using Infrastructure as Code.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22 (for development)

## Using the Provider

### Installation

Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    tierzero = {
      source  = "tierzero/tierzero"
      version = "~> 1.0"
    }
  }
}

provider "tierzero" {
  api_key = var.tierzero_api_key  # Or set TIERZERO_API_KEY environment variable
}
```

### Authentication

The provider requires a TierZero Organization API Key. You can provide it in two ways:

1. **Provider configuration**:
   ```hcl
   provider "tierzero" {
     api_key = "your-api-key"
   }
   ```

2. **Environment variable** (recommended):
   ```bash
   export TIERZERO_API_KEY="your-api-key"
   ```

### Quick Start

```hcl
# Create an alert responder
resource "tierzero_alert_responder" "production_critical" {
  team_name = "Production"
  name      = "Critical Error Handler"

  webhook_sources {
    type      = "PAGERDUTY"
    remote_id = "PXXXXXX"
  }

  matching_criteria {
    text_matches = ["critical", "fatal"]
  }

  runbook {
    prompt      = "Investigate and provide root cause analysis"
    fast_prompt = "Quick triage: assess severity and impact"
  }

  enabled = true
}
```

## Resources

- `tierzero_alert_responder` - Manages an alert responder that automatically investigates incoming alerts

## Data Sources

- `tierzero_webhook_subscriptions` - Lists available webhook subscriptions
- `tierzero_notification_integrations` - Lists available notification integrations

## Examples

See the [examples](./examples) directory for complete usage examples.

## Development

### Building the Provider

```bash
go build -o terraform-provider-tierzero
```

### Testing the Provider Locally

1. Build the provider:
   ```bash
   go build -o terraform-provider-tierzero
   ```

2. Create a `.terraformrc` file in your home directory:
   ```hcl
   provider_installation {
     dev_overrides {
       "tierzero/tierzero" = "/path/to/terraform-provider-tierzero"
     }
     direct {}
   }
   ```

3. Run Terraform commands in the `examples` directory

### Running Tests

```bash
# Unit tests
go test ./...

# Acceptance tests (requires TIERZERO_API_KEY)
TF_ACC=1 go test ./... -v
```

## Documentation

- [Provider Documentation](./docs/index.md)
- [Alert Responder Resource](./docs/resources/alert_responder.md)
- [Webhook Subscriptions Data Source](./docs/data-sources/webhook_subscriptions.md)
- [Notification Integrations Data Source](./docs/data-sources/notification_integrations.md)

## Contributing

Contributions are welcome! Please open an issue or pull request.

## License

This provider is licensed under the Mozilla Public License 2.0. See [LICENSE](./LICENSE) for details.

## Support

For issues and questions:
- GitHub Issues: https://github.com/tierzero/terraform-provider-tierzero/issues
- TierZero Documentation: https://docs.tierzero.com
