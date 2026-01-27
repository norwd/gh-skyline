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
  - package-ecosystem: "gomod"  # or npm, pip, docker, github-actions, etc.
    directory: "/"
    schedule:
      interval: "weekly"  # daily, weekly, or monthly
```

## Cooldown Settings

Configure cooldown periods to prevent excessive PR churn. This waits for packages to mature before creating update PRs:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    cooldown:
      default: "3d"       # Wait 3 days for all updates
      types:
        patch: "1d"       # Patch updates wait 1 day
        minor: "3d"       # Minor updates wait 3 days
        major: "7d"       # Major updates wait 7 days
      exclude:
        - "critical-pkg*" # Packages to exclude from cooldown
```

**Key options:**
- `default`: Minimum age for all updates unless overridden
- `types`: Per-semver level cooldowns (`patch`, `minor`, `major`)
- `exclude`: Package names or wildcards to skip cooldown (for urgent updates)

**Note:** Security updates automatically bypass cooldown periods.

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
          - "*"  # Group all Go dependencies
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
```

## Filtering Dependencies

### Allow specific dependency types

```yaml
allow:
  - dependency-type: "direct"    # Only direct dependencies
  - dependency-type: "indirect"  # Include transitive dependencies
```

### Ignore specific dependencies

```yaml
ignore:
  - dependency-name: "lodash"
    versions: ["4.x"]  # Ignore lodash 4.x updates
  - dependency-name: "aws-sdk"
    update-types: ["version-update:semver-major"]  # Ignore major updates
```

## Common Ecosystems

| Ecosystem | `package-ecosystem` value |
|-----------|---------------------------|
| Go modules | `gomod` |
| npm/yarn | `npm` |
| Python pip | `pip` |
| Docker | `docker` |
| GitHub Actions | `github-actions` |
| Terraform | `terraform` |
| Cargo (Rust) | `cargo` |

## Complete Example

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    cooldown:
      default: "3d"
      types:
        major: "7d"
    allow:
      - dependency-type: "direct"
      - dependency-type: "indirect"
    groups:
      go-dependencies:
        patterns:
          - "*"

  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    cooldown:
      default: "3d"
      types:
        major: "7d"
    groups:
      github-actions:
        patterns:
          - "*"
```
