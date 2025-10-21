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

# Basic alert responder example
resource "tierzero_alert_responder" "production_critical" {
  team_name = "Default"
  name      = "Production Critical Errors"

  webhook_sources = [{
    type      = "OPSGENIE"
    remote_id = "your-opsgenie-webhook-id"  # Replace with actual Opsgenie webhook ID
  }]

  matching_criteria = {
    text_matches = ["critical", "fatal", "emergency"]
  }

  enabled = true
}

# Advanced example with runbook and notifications
resource "tierzero_alert_responder" "automated_handler" {
  team_name = "Default"
  name      = "Automated Critical Alert Handler"

  webhook_sources = [{
    type      = "OPSGENIE"
    remote_id = "your-opsgenie-webhook-id"  # Replace with actual Opsgenie webhook ID
  }]

  matching_criteria = {
    text_matches = ["critical", "p1", "sev1"]
  }

  runbook = {
    prompt = <<-EOT
      Investigate this critical alert:
      1. Check recent deployments
      2. Review error rates and patterns
      3. Identify affected services
      4. Provide root cause analysis
      5. Suggest remediation steps
    EOT

    fast_prompt = "Quick impact assessment: determine severity, affected users, and business impact"
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
