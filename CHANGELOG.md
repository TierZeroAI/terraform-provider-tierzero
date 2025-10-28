# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.6] - 2025-10-28

### Changed
- **BREAKING**: Renamed runbook fields in `tierzero_alert_responder` resource to use more descriptive names:
  - `prompt` → `investigation_prompt`: Main investigation prompt for detailed root cause analysis
  - `fast_prompt` → `impact_and_severity_prompt`: Quick triage prompt for impact and severity assessment
- Updated all documentation and examples to use new field names
- Updated templates to reflect new field names

### Migration Guide
To migrate from 0.0.5 to 0.0.6, update your Terraform configuration:

```hcl
# Before (0.0.5)
runbook = {
  prompt      = "Investigate the issue..."
  fast_prompt = "Assess impact..."
}

# After (0.0.6)
runbook = {
  investigation_prompt        = "Investigate the issue..."
  impact_and_severity_prompt = "Assess impact..."
}
```

## [0.0.5] - 2025-10-22

### Added
- Provider now sets custom User-Agent header (`terraform-provider-tierzero/<version>`) for all API requests to identify Terraform-managed resources

### Fixed
- Fixed `created_at` field not being properly set during alert responder updates, which caused "Provider returned invalid result object after apply" errors

## [0.0.4] - 2025-10-21

### Added
- Slack-based alert responders: Added `slack_channel_id` field to `tierzero_alert_responder` resource as an alternative to `webhook_sources`
- Support for `slack_bot_app_user_id` in matching criteria to filter Slack messages by bot/sender
- Validation to ensure either `webhook_sources` or `slack_channel_id` is specified (mutually exclusive)
- Documentation and examples for Slack-based alert responders

### Changed
- **BREAKING**: `webhook_sources` is now optional instead of required in `tierzero_alert_responder` resource
- **BREAKING**: `team_name`, `webhook_sources`, and `slack_channel_id` fields in `tierzero_alert_responder` now require resource replacement when changed (API does not support switching alert types or changing team)
- Removed `SLACK` from webhook subscription types documentation (Slack channels use `slack_channel_id`, not webhook subscriptions)
- Updated runbook examples to use real-world API 500 error investigation scenario from TierZero prompt library

### Fixed
- Fixed perpetual diff issue with `team_name`, `webhook_sources`, and `slack_channel_id` changes by marking them as ForceNew (RequiresReplace)
- Fixed `matchingCriteriaChanged` helper function to check for changes to `slack_bot_app_user_id` field
- Clarified that webhook subscriptions data source only returns external providers (PagerDuty, OpsGenie, FireHydrant, Rootly)

## [0.0.3] - 2025-10-21

### Added
- Initial release of TierZero Terraform Provider
- `tierzero_alert_responder` resource for managing alert responders
- `tierzero_webhook_subscriptions` data source for discovering webhook subscriptions
- `tierzero_notification_integrations` data source for discovering notification integrations
- Support for OpsGenie, PagerDuty, FireHydrant, and Rootly webhook sources
- Support for Slack and Discord notification integrations
- Runbook configuration with prompt and fast_prompt
- Matching criteria based on text patterns
- Enable/disable functionality for alert responders

### Documentation
- Provider configuration examples
- Resource and data source examples
- Release process documentation

[Unreleased]: https://github.com/tierzero/terraform-provider-tierzero/compare/v0.0.6...HEAD
[0.0.6]: https://github.com/tierzero/terraform-provider-tierzero/compare/v0.0.5...v0.0.6
[0.0.5]: https://github.com/tierzero/terraform-provider-tierzero/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/tierzero/terraform-provider-tierzero/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/tierzero/terraform-provider-tierzero/releases/tag/v0.0.3
