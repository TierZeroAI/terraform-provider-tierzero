terraform {
  required_providers {
    tierzero = {
      source = "tierzero/tierzero"
    }
  }
}

provider "tierzero" {
  # API key from TIERZERO_API_KEY environment variable
}

# Fetch all notification integrations
data "tierzero_notification_integrations" "all" {}

# Fetch only Slack integrations
data "tierzero_notification_integrations" "slack" {
  kind = "SLACK_ALERT"
}

# Fetch only Discord integrations
data "tierzero_notification_integrations" "discord" {
  kind = "DISCORD_WEBHOOK"
}

# Output integrations
output "all_integrations" {
  value = data.tierzero_notification_integrations.all.notification_integrations
}

output "slack_integrations" {
  value = data.tierzero_notification_integrations.slack.notification_integrations
}

# Use in alert responder
resource "tierzero_alert_responder" "with_notifications" {
  team_name = "Production"
  name      = "Alert with Notifications"

  webhook_sources = [{
    type      = "PAGERDUTY"
    remote_id = "PXXXXXX"
  }]

  matching_criteria = {
    text_matches = ["critical"]
  }

  notification_integration_ids = [
    data.tierzero_notification_integrations.slack.notification_integrations[0].id
  ]

  enabled = true
}
