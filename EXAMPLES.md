# Changelog Examples

This document shows examples of how different commit types appear in the generated changelog.

## Example 1: Breaking Change with `!` Indicator

### Commit Message
```
feat!: redesign authentication API

The authentication flow has been completely redesigned to use JWT tokens
instead of session-based authentication.
```

### Changelog Output
```
### ⚠️  Breaking Changes

  • a1b2c3d - redesign authentication API

### Added

  • a1b2c3d - redesign authentication API
```

---

## Example 2: Breaking Change with Footer

### Commit Message
```
refactor(api): update user endpoint structure

Simplified the user API endpoints for better consistency.

BREAKING CHANGE: The /api/users endpoint now returns data in a different format. 
Update your client code to handle the new response structure.
```

### Changelog Output
```
### ⚠️  Breaking Changes

  • d4e5f6g - **api**: The /api/users endpoint now returns data in a different format. Update your client code to handle the new response structure.

### Changed

  • d4e5f6g - **api**: update user endpoint structure
```

---

## Example 3: Multiple Breaking Changes

### Commit Messages
```
feat!: remove deprecated payment methods

Removed support for legacy payment methods that were deprecated in v2.0.

BREAKING CHANGE: PayPal and Stripe v1 are no longer supported. Migrate to Stripe v3.

---

fix(auth)!: enforce password strength requirements

All passwords must now meet minimum strength requirements.
```

### Changelog Output
```
### ⚠️  Breaking Changes

  • h7i8j9k - **auth**: enforce password strength requirements
  • l0m1n2o - PayPal and Stripe v1 are no longer supported. Migrate to Stripe v3.

### Added

  • l0m1n2o - remove deprecated payment methods

### Fixed

  • h7i8j9k - **auth**: enforce password strength requirements
```

---

## Example 4: Regular Commits (No Breaking Changes)

### Commit Messages
```
feat(ui): add dark mode toggle

Users can now switch between light and dark themes in settings.

---

fix: resolve memory leak in cache layer

---

docs: update API documentation

---

chore: bump dependencies to latest versions
```

### Changelog Output
```
### Added

  • p3q4r5s - **ui**: add dark mode toggle

### Fixed

  • t6u7v8w - resolve memory leak in cache layer

### Documentation

  • x9y0z1a - update API documentation

### Maintenance

  • b2c3d4e - bump dependencies to latest versions
```

---

## Example 5: Scoped Commits

### Commit Messages
```
feat(auth): add two-factor authentication support

---

fix(database): improve connection pool handling

---

perf(api): optimize query performance for large datasets

---

test(auth): add integration tests for OAuth flow
```

### Changelog Output
```
### Added

  • f5g6h7i - **auth**: add two-factor authentication support

### Fixed

  • j8k9l0m - **database**: improve connection pool handling

### Performance

  • n1o2p3q - **api**: optimize query performance for large datasets

### Testing

  • r4s5t6u - **auth**: add integration tests for OAuth flow
```

---

## Example 6: Ignored Commits

### Commit Messages
```
ci: update deployment workflow [skip ci]

This commit will be filtered out and won't appear in the changelog.

---

feat: add new dashboard widgets

This commit will appear in the changelog.
```

### Changelog Output
```
### Added

  • v7w8x9y - add new dashboard widgets
```

_(The CI commit with `[skip ci]` is filtered out)_

---

## Example 7: Mixed Commit Types in a Release

### Full Tag Section
```
========================================
v2.0.0
========================================

### ⚠️  Breaking Changes

  • a1b2c3d - **api**: change response format to match REST standards
  • e4f5g6h - remove support for legacy authentication

### Added

  • i7j8k9l - **ui**: new admin dashboard
  • m0n1o2p - **auth**: OAuth2 provider integration
  • q3r4s5t - webhook support for events

### Fixed

  • u6v7w8x - **database**: connection timeout issues
  • y9z0a1b - memory leak in background worker

### Changed

  • c2d3e4f - **api**: simplified user endpoints

### Documentation

  • g5h6i7j - add migration guide for v2.0
  • k8l9m0n - update API reference

### Performance

  • o1p2q3r - optimize database queries

### Maintenance

  • s4t5u6v - update dependencies
  • w7x8y9z - clean up deprecated code
```

---

## Best Practices

1. **Use Breaking Changes Sparingly**: Only mark commits as breaking when they actually break backward compatibility
2. **Provide Migration Instructions**: In the `BREAKING CHANGE` footer, explain what users need to do
3. **Use Scopes**: Add scopes to help users understand which part of the codebase changed
4. **Keep Subjects Clear**: Write concise, descriptive commit subjects
5. **One Concern Per Commit**: Don't mix breaking and non-breaking changes in the same commit when possible

## Commit Message Template

```
<type>[optional scope][!]: <description>

[optional body]

[optional BREAKING CHANGE: description]
[optional other footers]
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `build`: Build system changes
- `ci`: CI/CD changes
- `chore`: Other changes (dependencies, etc.)