# WhatsApp Worker v1.1.0

## Release Date
December 27, 2025

## Overview
This release focuses on CI/CD improvements, test coverage enhancements, and workflow automation. The service now includes comprehensive GitHub Actions workflows for continuous integration and deployment.

## What's New

### CI/CD Workflows
- **GitHub Actions CI Workflow**: Automated testing on push and pull requests
  - Unit tests with coverage reporting
  - Integration tests with RabbitMQ service container
  - Docker image building and validation
  - Markdown test reports in GitHub Actions
- **GitHub Actions CD Workflow**: Automated deployment on version tags
  - Docker image building and pushing to Docker Hub
  - Multi-tag support (vX.Y.Z, X.Y.Z, latest)
- **Go Module Caching**: Optimized dependency installation with Go module cache

### Test Improvements
- **Test Coverage**: Improved unit test coverage
  - Added comprehensive metrics tests (100% coverage)
  - Enhanced notifier tests
  - Improved queue tests
  - Fixed unused import in queue tests
- **Test Reports**: Markdown-formatted test summaries in CI
- **Coverage Threshold**: Set to 30% (packages requiring external services excluded)
- **Integration Tests**: RabbitMQ service container configuration for CI

### Docker Improvements
- **Docker Build**: Enhanced Dockerfile with multi-stage build
  - CGO support for SQLite (required by whatsmeow)
  - SQLite libraries included in final image
- **Image Testing**: Automated Docker image validation in CI

## Technical Details

### Dependencies
- No breaking changes to dependencies
- Go 1.24 (as specified in go.mod)
- CGO enabled for SQLite support (whatsmeow library)

### Configuration
- No configuration changes required
- RabbitMQ vhost configuration for integration tests

## Migration Guide
No migration required. This is a non-breaking release.

## Full Changelog

### Added
- GitHub Actions CI workflow (`.github/workflows/ci.yml`)
- GitHub Actions CD workflow (`.github/workflows/cd.yml`)
- Markdown test reports in CI
- Go module caching in CI workflows
- Docker image validation in CI
- Comprehensive metrics tests (`pkg/metrics/metrics_test.go`) with 100% coverage
- Enhanced notifier tests for phone number formats
- Improved queue tests

### Changed
- Improved test coverage reporting
- Enhanced RabbitMQ service container configuration
- Better test result parsing in CI with multiline output format
- Excluded `cmd/whatsapp-worker` from coverage (main functions are hard to test)
- Lowered coverage threshold to 30% (packages requiring external services)

### Fixed
- Fixed unused `log/slog` import in queue tests
- Fixed test result parsing to handle zero counts correctly
- Improved GitHub Actions output format to prevent "Invalid format" errors
- Enhanced coverage extraction from Go test output

## Known Issues
- Integration tests may fail due to RabbitMQ vhost permissions (to be addressed in future release)

## Contributors
- Automated CI/CD implementation
- Test coverage improvements

