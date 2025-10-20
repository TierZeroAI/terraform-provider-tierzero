terraform {
  required_providers {
    tierzero = {
      source = "tierzero/tierzero"
    }
  }
}

provider "tierzero" {
  # API key from TIERZERO_API_KEY environment variable
  # base_url defaults to https://api.tierzero.ai
}

# Fetch all available webhook subscriptions
data "tierzero_webhook_subscriptions" "all" {}

# Output all webhook subscriptions
output "webhook_subscriptions" {
  value = data.tierzero_webhook_subscriptions.all.webhook_subscriptions
}

# Filter for PagerDuty webhooks
output "pagerduty_webhooks" {
  value = [
    for ws in data.tierzero_webhook_subscriptions.all.webhook_subscriptions :
    ws if ws.type == "PAGERDUTY"
  ]
}

# Use in alert responder
resource "tierzero_alert_responder" "example" {
  team_name = "Production"
  name      = "Example Alert"

  webhook_sources = [{
    type      = data.tierzero_webhook_subscriptions.all.webhook_subscriptions[0].type
    remote_id = data.tierzero_webhook_subscriptions.all.webhook_subscriptions[0].remote_id
  }]

  matching_criteria = {
    text_matches = ["error"]
  }
}
