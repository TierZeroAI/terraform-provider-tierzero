---
page_title: "Resource tierzero_alert_responder - terraform-provider-tierzero"
subcategory: ""
description: |-
  Manages a TierZero Alert Responder that automatically investigates incoming alerts.
---

# Resource: tierzero_alert_responder

Manages a TierZero Alert Responder that automatically investigates incoming alerts.

Alert responders automatically investigate alerts matching specified criteria from configured webhook sources. When an alert matches the text patterns, TierZero's AI agent analyzes the alert using the configured runbook prompts and can send results to notification channels.

## Key Concepts

- **Alert Types**: Two types of alert responders:
  - **Webhook-based**: Monitor alerts from PagerDuty, OpsGenie, FireHydrant, or Rootly via webhook integrations
  - **Slack-based**: Monitor Slack channel messages directly (requires slack_channel_id instead of webhook_sources)
- **Matching Criteria**: Define text patterns that trigger automated investigation. For Slack alerts, optionally filter by bot/sender using `slack_bot_app_user_id`
- **Runbook**: Customize investigation behavior with two types of prompts:
  - `prompt`: Main investigation directive for detailed root cause analysis. Use this to define the investigation steps, queries to run, and analysis approach. Default: "Please investigate the issue and explain the root cause to the best of your abilities!"
  - `fast_prompt`: Quick triage directive for rapid severity and impact assessment. Use this to quickly determine how many users or accounts are affected. Example fast_prompt:
    ```
    Determine how many users were affected by the 500 error.
    Use the spans aggregation query using the filter:
    env:prod @http.method:<HTTP_METHOD> @http.route:* @http.status_code:500
    and facet on @usr.id.
    ```
- **Status**: Control whether the responder is ACTIVE (`enabled = true`) or PAUSED (`enabled = false`)

For more runbook examples, see the [TierZero Prompt Library](https://docs.tierzero.ai/prompt-library/alert-responder).

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `matching_criteria` (Attributes) Criteria for matching alerts (see [below for nested schema](#nestedatt--matching_criteria))
- `name` (String) Alert responder name
- `team_name` (String) Team name

**Note**: Must specify **either** `webhook_sources` **or** `slack_channel_id` (mutually exclusive, not both).

### Optional

- `webhook_sources` (Attributes List) Webhook sources to monitor (for PagerDuty, OpsGenie, FireHydrant, Rootly). Mutually exclusive with `slack_channel_id`. (see [below for nested schema](#nestedatt--webhook_sources))
- `slack_channel_id` (String) Slack channel ID (e.g., 'C01234567' for public channels, 'G01234567' for private channels). Mutually exclusive with `webhook_sources`.
- `enabled` (Boolean) Whether the alert responder is enabled. When true, status is ACTIVE. When false, status is PAUSED. Uses enable/disable API endpoints under the hood.
- `notification_integration_ids` (List of String) Notification integration Global IDs
- `runbook` (Attributes) Investigation runbook (optional, uses default if not provided) (see [below for nested schema](#nestedatt--runbook))

### Read-Only

- `created_at` (String) Creation timestamp (ISO 8601)
- `id` (String) Alert Responder Global ID
- `updated_at` (String) Last update timestamp (ISO 8601)
- `url` (String) Link to alert responder details page (returned by create/update operations)

<a id="nestedatt--matching_criteria"></a>
### Nested Schema for `matching_criteria`

Required:

- `text_matches` (List of String) Array of text patterns to match

Optional:

- `slack_bot_app_user_id` (String) Optional Slack bot/sender app user ID to filter messages (only for Slack alerts)


<a id="nestedatt--webhook_sources"></a>
### Nested Schema for `webhook_sources`

Required:

- `remote_id` (String) External webhook ID
- `type` (String) Webhook type (PAGERDUTY, OPSGENIE, FIREHYDRANT, ROOTLY)


<a id="nestedatt--runbook"></a>
### Nested Schema for `runbook`

Optional:

- `fast_prompt` (String) Quick triage prompt
- `prompt` (String) Main investigation prompt

## Import

Import is supported using the alert responder's Global ID:

```shell
#!/bin/bash
# Import an existing alert responder by its Global ID
terraform import tierzero_alert_responder.production_critical "R3JhcGhRTEFsZXJ0OjEyMw=="
```
