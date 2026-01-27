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

Configure cooldown periods to delay updates until packages have matured. This helps avoid churn from rapid releases. Cooldown only applies to version updates, not security updates.

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    cooldown:
      default-days: 3           # Default cooldown for all updates
      semver-major-days: 7      # Major version updates wait 7 days
      semver-minor-days: 3      # Minor version updates wait 3 days
      semver-patch-days: 1      # Patch version updates wait 1 day
      include:
        - "some-package*"       # Only apply cooldown to matching packages
      exclude:
        - "critical-pkg*"       # Skip cooldown for these packages
```

**Parameters:**

| Parameter | Description |
|-----------|-------------|
| `default-days` | Default cooldown period for all dependencies |
| `semver-major-days` | Cooldown for major version updates |
| `semver-minor-days` | Cooldown for minor version updates |
| `semver-patch-days` | Cooldown for patch version updates |
| `include` | List of dependencies to apply cooldown (supports wildcards) |
| `exclude` | List of dependencies excluded from cooldown (supports wildcards) |

**Notes:**
- If semver-specific days aren't defined, `default-days` is used
- `exclude` takes precedence over `include`
- Security updates automatically bypass cooldown

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
  - dependency-type: "direct"    # Only direct dependencies
  - dependency-type: "indirect"  # Include transitive dependencies
  - dependency-type: "all"       # All dependencies
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
| NuGet (.NET) | `nuget` |

## Complete Example

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    cooldown:
      default-days: 3
      semver-major-days: 7
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
      default-days: 3
      semver-major-days: 7
    groups:
      github-actions:
        patterns:
          - "*"
```

---

## References

- [Dependabot Configuration Options](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/configuration-options-for-the-dependabot.yml-file) - Full configuration reference
- [Optimizing PR Creation](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/optimizing-pr-creation-version-updates) - Cooldown and grouping strategies
