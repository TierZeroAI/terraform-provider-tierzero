# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/tierzero/terraform-provider-tierzero/compare/v0.0.4...HEAD
[0.0.4]: https://github.com/tierzero/terraform-provider-tierzero/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/tierzero/terraform-provider-tierzero/releases/tag/v0.0.3
