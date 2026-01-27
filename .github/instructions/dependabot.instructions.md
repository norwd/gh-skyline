---
applyTo: "**/*dependabot.yml"
description: Dependabot configuration patterns and best practices
---

# Dependabot Configuration

Guidelines for configuring Dependabot version updates.

## Basic Structure

```yaml
version: 2
updates:
  - package-ecosystem: "gomod" # or npm, pip, docker, github-actions, etc.
    directory: "/"
    schedule:
      interval: "weekly" # daily, weekly, or monthly
```

## Cooldown Settings

Configure cooldown periods to delay updates until packages have matured. This helps avoid churn from rapid releases. Cooldown only applies to version updates, not security updates.

```yaml
cooldown:
  default-days: 7 # Default cooldown for all updates
  semver-major-days: 7 # Major version updates wait 7 days
  semver-minor-days: 3 # Minor version updates wait 3 days
  semver-patch-days: 1 # Patch version updates wait 1 day
  include:
    - "some-package*" # Only apply cooldown to matching packages
  exclude:
    - "critical-pkg*" # Skip cooldown for these packages
```

**Parameters:**

| Parameter           | Description                                                      |
| ------------------- | ---------------------------------------------------------------- |
| `default-days`      | Default cooldown period for all dependencies                     |
| `semver-major-days` | Cooldown for major version updates                               |
| `semver-minor-days` | Cooldown for minor version updates                               |
| `semver-patch-days` | Cooldown for patch version updates                               |
| `include`           | List of dependencies to apply cooldown (supports wildcards)      |
| `exclude`           | List of dependencies excluded from cooldown (supports wildcards) |

**Notes:**

- If semver-specific days aren't defined, `default-days` is used
- `exclude` takes precedence over `include`
- Security updates automatically bypass cooldown

### SemVer Cooldown Support by Ecosystem

**IMPORTANT:** The `semver-major-days`, `semver-minor-days`, and `semver-patch-days` options are NOT supported by all package ecosystems. For unsupported ecosystems, use only `default-days`.

| Ecosystem        | SemVer Cooldown Supported       |
| ---------------- | ------------------------------- |
| `gomod`          | ✅ Yes                          |
| `npm`            | ✅ Yes                          |
| `pip`            | ✅ Yes                          |
| `bundler`        | ✅ Yes                          |
| `cargo`          | ✅ Yes                          |
| `composer`       | ✅ Yes                          |
| `maven`          | ✅ Yes                          |
| `gradle`         | ✅ Yes                          |
| `nuget`          | ✅ Yes                          |
| `docker`         | ✅ Yes                          |
| `github-actions` | ❌ No - use `default-days` only |
| `gitsubmodule`   | ❌ No - use `default-days` only |
| `terraform`      | ❌ No - use `default-days` only |

**Example for github-actions (no semver support):**

```yaml
- package-ecosystem: "github-actions"
  directory: "/"
  schedule:
    interval: "weekly"
  cooldown:
    default-days: 7 # Only default-days is supported
```

## Grouping Updates

Group related dependencies into single PRs to reduce noise:

```yaml
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    groups:
      go-dependencies:
        patterns:
          - "*" # Group all Go dependencies
```

You can create multiple groups with specific patterns:

```yaml
groups:
  aws-sdk:
    patterns:
      - "github.com/aws/*"
  testing:
    patterns:
      - "*test*"
      - "*mock*"
    dependency-type: "development"
```

Group parameters:

- `patterns`: Include dependencies matching these patterns
- `exclude-patterns`: Exclude dependencies matching these patterns
- `dependency-type`: Limit to `development` or `production`
- `update-types`: Limit to `minor`, `patch`, or `major`

## Filtering Dependencies

### Allow specific dependency types

```yaml
allow:
  - dependency-type: "direct" # Only direct dependencies
  - dependency-type: "indirect" # Include transitive dependencies
  - dependency-type: "all" # All dependencies
```

### Ignore specific dependencies

```yaml
ignore:
  - dependency-name: "lodash"
    versions: ["4.x"] # Ignore lodash 4.x updates
  - dependency-name: "aws-sdk"
    update-types: ["version-update:semver-major"] # Ignore major updates
```

## Common Ecosystems

| Ecosystem      | `package-ecosystem` value |
| -------------- | ------------------------- |
| Go modules     | `gomod`                   |
| npm/Yarn       | `npm`                     |
| Python pip     | `pip`                     |
| Docker         | `docker`                  |
| GitHub Actions | `github-actions`          |
| Terraform      | `terraform`               |
| Cargo (Rust)   | `cargo`                   |
| NuGet (.NET)   | `nuget`                   |

## Complete Example

```yaml
version: 2
updates:
  # Go modules - supports semver cooldown
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    cooldown:
      default-days: 7
      semver-major-days: 7
    allow:
      - dependency-type: "direct"
      - dependency-type: "indirect"
    groups:
      go-dependencies:
        patterns:
          - "*"

  # GitHub Actions - does NOT support semver cooldown
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    cooldown:
      default-days: 7
    groups:
      github-actions:
        patterns:
          - "*"
```

---

## References

- [Dependabot Configuration Options](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/configuration-options-for-the-dependabot.yml-file) - Full configuration reference
- [Supported Ecosystems](https://docs.github.com/en/code-security/dependabot/ecosystems-supported-by-dependabot/supported-ecosystems-and-repositories) - List of supported package ecosystems
- [Optimizing PR Creation](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/optimizing-pr-creation-version-updates) - Cooldown and grouping strategies
