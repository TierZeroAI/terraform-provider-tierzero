# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - TBD

### Added
- Initial release of TierZero Terraform Provider
- `tierzero_alert_responder` resource for managing alert responders
- `tierzero_webhook_subscriptions` data source for discovering webhook subscriptions
- `tierzero_notification_integrations` data source for discovering notification integrations
- Support for OpsGenie, PagerDuty, FireHydrant, Rootly, and Slack webhook sources
- Support for Slack and Discord notification integrations
- Runbook configuration with prompt and fast_prompt
- Matching criteria based on text patterns
- Enable/disable functionality for alert responders

### Documentation
- Provider configuration examples
- Resource and data source examples
- Release process documentation

[Unreleased]: https://github.com/tierzero/terraform-provider-tierzero/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/tierzero/terraform-provider-tierzero/releases/tag/v1.0.0
