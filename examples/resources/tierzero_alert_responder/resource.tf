terraform {
  required_providers {
    tierzero = {
      source = "tierzeroai/tierzero"
    }
  }
}

provider "tierzero" {
  # API key from TIERZERO_API_KEY environment variable
  # base_url defaults to https://api.tierzero.ai
}

# Webhook-based alert responder example (PagerDuty, OpsGenie, FireHydrant, Rootly)
resource "tierzero_alert_responder" "production_critical" {
  team_name = "Default"
  name      = "Production Critical Errors"

  webhook_sources = [{
    type      = "OPSGENIE"
    remote_id = "your-opsgenie-webhook-id"  # Replace with actual OpsGenie webhook ID
  }]

  matching_criteria = {
    text_matches = ["critical", "fatal", "emergency"]
  }

  enabled = true
}

# Slack-based alert responder example
resource "tierzero_alert_responder" "slack_database_alerts" {
  team_name = "Default"
  name      = "Slack Database Alerts"

  slack_channel_id = "C07TUN1EFFU"  # Slack channel ID (C for public, G for private)

  matching_criteria = {
    text_matches = ["database", "error", "timeout"]
  }

  enabled = true
}

# Slack alert responder with bot filter
# Runbook is optional - if not specified, uses the default investigation prompt
resource "tierzero_alert_responder" "slack_datadog_alerts" {
  team_name = "Default"
  name      = "Slack Datadog Alerts"

  slack_channel_id = "C07TUN1EFFU"

  matching_criteria = {
    text_matches          = ["alert", "warning"]
    slack_bot_app_user_id = "B01234567"  # Optional: filter by bot/sender app user ID
  }

  enabled = true
}

# Advanced example with custom runbook and notifications
# For more runbook examples, see: https://docs.tierzero.ai/prompt-library/alert-responder
resource "tierzero_alert_responder" "api_500_errors" {
  team_name = "Default"
  name      = "API 500 Error Handler"

  webhook_sources = [{
    type      = "OPSGENIE"
    remote_id = "your-opsgenie-webhook-id"  # Replace with actual OpsGenie webhook ID
  }]

  matching_criteria = {
    text_matches = ["500", "error", "api"]
  }

  runbook = {
    prompt = <<-EOT
      API requests are returning 500 errors. Investigate following these steps:
      1. Execute a spans query filtering for env:prod @http.method:<HTTP_METHOD> @http.route:* @http.status_code:500 and group by @usr.id to quantify affected users
      2. Perform separate spans aggregations to determine impacted accounts (facet on @usr.accountId) and users
      3. Collect and examine at least 5 trace IDs with 500 errors, investigate each trace with error status filtering
      4. If an error stack trace is identified with a version/git hash, investigate commits from up to 3 days prior. Flag potentially related commits as investigation leads
    EOT

    fast_prompt = "Determine how many users were affected by the 500 errors. Use spans aggregation query with filter: env:prod @http.method:<HTTP_METHOD> @http.route:* @http.status_code:500 and facet on @usr.id"
  }

  notification_integration_ids = [
    "R3JhcGhRTE5vdGlmaWNhdGlvbkludGVncmF0aW9uOjEyMw=="
  ]

  enabled = true
}

# Example with data source discovery
data "tierzero_webhook_subscriptions" "available" {}

data "tierzero_notification_integrations" "slack" {
  kind = "SLACK_ALERT"
}

locals {
  opsgenie_webhook = [
    for ws in data.tierzero_webhook_subscriptions.available.webhook_subscriptions :
    ws if ws.type == "OPSGENIE"
  ][0]
}

resource "tierzero_alert_responder" "discovered" {
  team_name = "Default"
  name      = "Alert Handler with Discovery"

  webhook_sources = [{
    type      = local.opsgenie_webhook.type
    remote_id = local.opsgenie_webhook.remote_id
  }]

  matching_criteria = {
    text_matches = ["error", "warning"]
  }

  notification_integration_ids = length(data.tierzero_notification_integrations.slack.notification_integrations) > 0 ? [
    data.tierzero_notification_integrations.slack.notification_integrations[0].id
  ] : []

  enabled = true
}
