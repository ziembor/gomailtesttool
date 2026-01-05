# Integration Test Results Archive

This directory contains historical integration test results for the Microsoft Graph GoLang Testing Tool.

## Purpose

This archive maintains a record of integration test runs performed on the application, primarily by Claude AI assistants. Each test report provides comprehensive verification of application functionality at a specific point in time.

## File Naming Convention

Test result files follow this naming pattern:
```
INTEGRATION_TEST_RESULTS_YYYY-MM-DD.md
```

Example: `INTEGRATION_TEST_RESULTS_2026-01-05.md`

## Test Reports

| Date | Version | Platform | Tester | Status | Report |
|------|---------|----------|--------|--------|--------|
| 2026-01-05 | 1.16.8 | Windows (Git Bash) | Claude | ✅ PASS (53/53) | [View Report](INTEGRATION_TEST_RESULTS_2026-01-05.md) |

## What's Tested

Each integration test report typically covers:

1. **Build Verification** - Binary compilation and version check
2. **API Operations** - All four actions (getinbox, getevents, sendmail, sendinvite)
3. **Authentication** - Client secret, PFX, and/or Windows Certificate Store
4. **CSV Logging** - Verification of log files and schemas
5. **Verbose Mode** - Debug output and diagnostic information
6. **Environment Variables** - MSGRAPH* variable support
7. **Input Validation** - Error handling and edge cases
8. **Unit Tests** - Go test suite execution

## For Future Testing

When performing integration tests, please:

1. Create a new test report with the current date
2. Follow the standardized format from existing reports
3. Update this README with the new test entry in the table above
4. Commit both the test report and updated README

## Test Execution Guide

For instructions on how to run integration tests, see:
- [INTEGRATION_TESTS.md](../INTEGRATION_TESTS.md) - Complete testing guide
- [run-integration-tests.ps1](../run-integration-tests.ps1) - Automated test runner

## Historical Context

This archive was established on 2026-01-05 to maintain a permanent record of integration test results, particularly those performed by AI assistants to ensure consistency and traceability of testing activities.
