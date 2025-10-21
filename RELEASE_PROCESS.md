# Release Process for TierZero Terraform Provider

This document describes how to create and publish a new release of the TierZero Terraform Provider.

## Overview

Releases are fully automated using GitHub Actions and GoReleaser. When you push a version tag, the system automatically:
1. Builds binaries for all supported platforms (macOS, Linux, Windows)
2. Signs the release with GPG
3. Creates a GitHub release
4. Publishes to the Terraform Registry

## Prerequisites

- ✅ GPG key configured in GitHub Secrets (`GPG_PRIVATE_KEY`, `GPG_PASSPHRASE`)
- ✅ Repository connected to [Terraform Registry](https://registry.terraform.io)
- ✅ All tests passing
- ✅ CHANGELOG.md updated with new version

## Release Steps

### 1. Prepare the Release

```bash
# Ensure you're on main branch with latest changes
git checkout main
git pull origin main

# Update CHANGELOG.md with new version and changes
# Example:
# ## [1.2.0] - 2025-10-21
# ### Added
# - New data source for alert templates
# ### Fixed
# - Bug in webhook source validation

# Commit the changelog
git add CHANGELOG.md
git commit -m "Prepare v1.2.0 release"
git push origin main
```

### 2. Create and Push the Version Tag

```bash
# Create a tag following semantic versioning (MAJOR.MINOR.PATCH)
git tag v1.2.0

# Push the tag to trigger the release workflow
git push origin v1.2.0
```

**Important:** Tag format must be `vX.Y.Z` (e.g., `v1.0.0`, `v1.2.3`, `v2.0.0`)

### 3. Monitor the Release

1. Go to [GitHub Actions](https://github.com/tierzero/terraform-provider-tierzero/actions)
2. Watch the "Release" workflow (takes ~5-10 minutes)
3. If successful, a new release appears at [Releases](https://github.com/tierzero/terraform-provider-tierzero/releases)

### 4. Verify the Release

Check that the release includes:
- ✅ Binaries for all platforms (6+ ZIP files)
- ✅ `terraform-provider-tierzero_X.Y.Z_SHA256SUMS` (checksums)
- ✅ `terraform-provider-tierzero_X.Y.Z_SHA256SUMS.sig` (GPG signature)
- ✅ Release notes with changelog

### 5. Terraform Registry Auto-Updates

The [Terraform Registry](https://registry.terraform.io/providers/tierzero/tierzero) automatically detects the new release within ~15 minutes. No manual action needed!

## Version Numbers (Semantic Versioning)

Follow [SemVer](https://semver.org/):

- **MAJOR** (v2.0.0): Breaking changes, incompatible API changes
- **MINOR** (v1.2.0): New features, backward-compatible
- **PATCH** (v1.0.1): Bug fixes, backward-compatible

Examples:
```bash
git tag v1.0.0  # Initial release
git tag v1.1.0  # Added new data source (minor)
git tag v1.1.1  # Fixed a bug (patch)
git tag v2.0.0  # Renamed resource attribute (major - breaking)
```

## Pre-releases

For beta/RC versions, use a suffix:

```bash
git tag v1.2.0-rc1      # Release candidate 1
git tag v1.2.0-beta.1   # Beta 1
git push origin v1.2.0-rc1
```

These are marked as "Pre-release" on GitHub and Terraform Registry.

## Troubleshooting

### Release workflow fails

1. Check the [Actions logs](https://github.com/tierzero/terraform-provider-tierzero/actions)
2. Common issues:
   - Tests failing: Fix tests and create a new tag
   - GPG signing error: Verify `GPG_PRIVATE_KEY` secret is set correctly
   - Build errors: Check Go version compatibility

### Need to redo a release

If you need to fix a release:

```bash
# Delete the tag locally and remotely
git tag -d v1.2.0
git push origin :refs/tags/v1.2.0

# Delete the GitHub release manually in the web UI
# Fix the issues, then create the tag again
git tag v1.2.0
git push origin v1.2.0
```

### Terraform Registry not updating

- Wait 15-30 minutes (registry polls GitHub periodically)
- Verify the release has proper artifacts and signature
- Check [Terraform Registry provider settings](https://registry.terraform.io/settings/providers)

## Quick Reference

```bash
# Standard release workflow
git checkout main
git pull
# Update CHANGELOG.md
git add CHANGELOG.md
git commit -m "Prepare vX.Y.Z release"
git push
git tag vX.Y.Z
git push origin vX.Y.Z
# Wait for GitHub Actions to complete
# Verify at https://github.com/tierzero/terraform-provider-tierzero/releases
```

## Questions?

- GoReleaser docs: https://goreleaser.com
- Terraform Registry docs: https://www.terraform.io/registry/providers/publishing
- GitHub Actions docs: https://docs.github.com/en/actions
