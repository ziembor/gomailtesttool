# Architecture Diagram - gomailtesttool Suite

## Overview

**gomailtesttool** is a unified CLI (`gomailtest`) for email infrastructure testing, with 5 protocol subcommands plus developer tooling:
- **smtp** - SMTP connectivity and TLS diagnostics
- **imap** - IMAP server testing with OAuth2
- **pop3** - POP3 server testing with OAuth2
- **jmap** - JMAP protocol testing
- **msgraph** - Microsoft Graph API (Exchange Online)
- **devtools** - Release automation and environment management

## File Structure and Dependencies

```
gomailtesttool/
├── cmd/
│   └── gomailtest/                   # Single binary entry point
│       ├── main.go                   # main() → Execute()
│       └── root.go                   # Cobra root command, registers subcommands
│
├── internal/
│   ├── common/                       # Cross-protocol utilities
│   │   ├── bootstrap/
│   │   │   └── bootstrap.go          # Signal context setup
│   │   ├── logger/                   # Structured logging
│   │   │   ├── csv.go
│   │   │   ├── json.go
│   │   │   ├── json_test.go
│   │   │   ├── logger.go
│   │   │   └── slog.go               # slog-based structured logger
│   │   ├── ratelimit/
│   │   │   ├── ratelimit.go
│   │   │   └── ratelimit_test.go
│   │   ├── retry/
│   │   │   └── retry.go              # Exponential backoff
│   │   ├── security/
│   │   │   ├── masking.go
│   │   │   └── masking_test.go
│   │   ├── validation/
│   │   │   ├── validation.go
│   │   │   ├── validation_test.go
│   │   │   └── proxy_test.go
│   │   └── version/
│   │       └── version.go            # Single source of truth for version
│   │
│   ├── devtools/                     # Developer-facing subcommand
│   │   ├── cmd.go                    # 'gomailtest devtools' root
│   │   ├── env/
│   │   │   ├── env.go                # MSGRAPH* env var management
│   │   │   └── env_cmd.go
│   │   └── release/
│   │       ├── release.go            # Release orchestration
│   │       ├── release_cmd.go
│   │       ├── version.go            # Version bump logic
│   │       ├── changelog.go          # ChangeLog/{version}.md creation
│   │       ├── git.go                # git commit/tag/push
│   │       ├── gh.go                 # GitHub PR/release via gh CLI
│   │       ├── security_scan.go      # Pre-release secret scanning
│   │       ├── editor.go             # Interactive editor prompts
│   │       └── prompt.go             # User prompts
│   │
│   ├── protocols/                    # Protocol implementations
│   │   ├── smtp/
│   │   │   ├── cmd.go                # Cobra subcommand wiring
│   │   │   ├── config.go
│   │   │   ├── config_test.go
│   │   │   ├── smtp_client.go
│   │   │   ├── smtp_client_test.go
│   │   │   ├── testconnect.go
│   │   │   ├── teststarttls.go
│   │   │   ├── testauth.go
│   │   │   ├── sendmail.go
│   │   │   ├── sendmail_test.go
│   │   │   ├── tls_display.go
│   │   │   ├── utils.go
│   │   │   └── utils_test.go
│   │   │
│   │   ├── imap/
│   │   │   ├── cmd.go
│   │   │   ├── config.go
│   │   │   ├── imap_client.go
│   │   │   ├── listfolders.go
│   │   │   ├── testconnect.go
│   │   │   ├── testauth.go
│   │   │   └── utils.go
│   │   │
│   │   ├── pop3/
│   │   │   ├── cmd.go
│   │   │   ├── config.go
│   │   │   ├── pop3_client.go
│   │   │   ├── listmail.go
│   │   │   ├── testconnect.go
│   │   │   ├── testauth.go
│   │   │   └── utils.go
│   │   │
│   │   ├── jmap/
│   │   │   ├── cmd.go
│   │   │   ├── config.go
│   │   │   ├── config_test.go
│   │   │   ├── jmap_client.go
│   │   │   ├── getmailboxes.go
│   │   │   ├── testconnect.go
│   │   │   ├── testauth.go
│   │   │   ├── utils.go
│   │   │   └── utils_test.go
│   │   │
│   │   └── msgraph/
│   │       ├── cmd.go
│   │       ├── config.go
│   │       ├── auth.go
│   │       ├── handlers.go
│   │       ├── utils.go
│   │       ├── utils_test.go
│   │       ├── cert_windows.go       # Windows cert store (build: windows)
│   │       └── cert_stub.go          # Cross-platform stub (build: !windows)
│   │
│   ├── smtp/                         # SMTP protocol primitives
│   │   ├── exchange/
│   │   │   └── detection.go
│   │   └── protocol/
│   │       ├── capabilities.go
│   │       ├── commands.go
│   │       ├── commands_test.go
│   │       ├── responses.go
│   │       └── responses_test.go
│   │
│   ├── imap/protocol/
│   │   ├── capabilities.go
│   │   └── capabilities_test.go
│   │
│   ├── pop3/protocol/
│   │   ├── capabilities.go
│   │   ├── capabilities_test.go
│   │   ├── commands.go
│   │   ├── commands_test.go
│   │   └── responses.go
│   │
│   └── jmap/protocol/
│       ├── methods.go
│       ├── methods_test.go
│       ├── session.go
│       ├── session_test.go
│       ├── types.go
│       └── types_test.go
│
├── tests/
│   ├── README.md
│   └── integration/
│       └── sendmail_test.go          # MS Graph integration tests (build: integration)
│
├── scripts/
│   └── check-integration-env.sh     # Validates MSGRAPH* env vars before integration tests
│
├── ChangeLog/                        # Per-version changelogs
├── Makefile                          # Primary build system
├── build-all.ps1                     # Windows build script
├── run-integration-tests.ps1         # Integration test runner
├── run-interactive-release.ps1       # Legacy release script (superseded by devtools)
├── go.mod
└── go.sum
```

## Command Structure

```
gomailtest
├── smtp
│   ├── testconnect
│   ├── teststarttls
│   ├── testauth
│   └── sendmail
├── imap
│   ├── testconnect
│   ├── testauth
│   └── listfolders
├── pop3
│   ├── testconnect
│   ├── testauth
│   └── listmail
├── jmap
│   ├── testconnect
│   ├── testauth
│   └── getmailboxes
├── msgraph
│   ├── getevents
│   ├── sendinvite
│   ├── getschedule
│   ├── sendmail
│   ├── getinbox
│   ├── exportinbox
│   └── searchandexport
└── devtools
    ├── env       (get/set/clear MSGRAPH* environment variables)
    └── release   (interactive: version bump → changelog → git tag → GitHub release)
```

## Build System

### Makefile (primary)

```
make build          → go build -ldflags="-s -w" -o bin/gomailtest ./cmd/gomailtest
make build-verbose  → same with -v flag
make test           → go test ./...
make integration-test → build + check env + go test -tags integration ./tests/integration/
make clean          → rm -f bin/gomailtest[.exe]
make help           → list targets
```

### build-all.ps1 (Windows convenience)

```
.\build-all.ps1           → build bin/gomailtest.exe
.\build-all.ps1 -Verbose  → build with verbose Go output
```

## Application Flow

```
gomailtest <subcommand> [flags]
          │
          ▼
cmd/gomailtest/root.go
  rootCmd.AddCommand(smtp.NewCmd())
  rootCmd.AddCommand(imap.NewCmd())
  rootCmd.AddCommand(pop3.NewCmd())
  rootCmd.AddCommand(jmap.NewCmd())
  rootCmd.AddCommand(msgraph.NewCmd())
  rootCmd.AddCommand(devtools.NewCmd())
          │
          ▼
internal/protocols/<protocol>/cmd.go   ← flags, validation, dispatch
          │
          ▼
internal/protocols/<protocol>/<action>.go  ← operation logic
          │
          ├─► internal/common/logger/    ← CSV/JSON/slog output
          ├─► internal/common/retry/     ← exponential backoff
          ├─► internal/common/ratelimit/ ← token bucket
          └─► internal/<protocol>/protocol/ ← protocol primitives
```

## Protocol Implementations

### smtp (internal/protocols/smtp/)

```
cmd.go
  └─► testconnect.go     — TCP connectivity test
  └─► teststarttls.go    — TLS handshake, cert validation, cipher strength, Exchange detection
  └─► testauth.go        — PLAIN, LOGIN, CRAM-MD5, XOAUTH2
  └─► sendmail.go        — send test email
  └─► smtp_client.go     — SMTP client logic
  └─► tls_display.go     — TLS diagnostic output
```

### imap (internal/protocols/imap/)

```
cmd.go
  └─► testconnect.go     — TCP/TLS connectivity
  └─► testauth.go        — PLAIN, LOGIN, XOAUTH2
  └─► listfolders.go     — list IMAP folders
  └─► imap_client.go     — IMAP client logic
```

### pop3 (internal/protocols/pop3/)

```
cmd.go
  └─► testconnect.go     — TCP/TLS connectivity
  └─► testauth.go        — USER/PASS, APOP, XOAUTH2
  └─► listmail.go        — retrieve message list
  └─► pop3_client.go     — POP3 client logic
```

### jmap (internal/protocols/jmap/)

```
cmd.go
  └─► testconnect.go     — JMAP session discovery
  └─► testauth.go        — Basic, Bearer
  └─► getmailboxes.go    — list JMAP mailboxes
  └─► jmap_client.go     — HTTP-based JMAP client
```

### msgraph (internal/protocols/msgraph/)

```
cmd.go → handlers.go
  ├─► Calendar
  │   ├─► handleGetEvents()     (getevents)
  │   ├─► handleSendInvite()    (sendinvite)
  │   └─► handleGetSchedule()   (getschedule)
  ├─► Mail
  │   ├─► handleSendMail()      (sendmail)
  │   └─► handleGetInbox()      (getinbox)
  └─► Export
      ├─► handleExportInbox()        (exportinbox)
      └─► handleSearchAndExport()    (searchandexport)

auth.go → getCredential()
  ├─► azidentity.NewClientSecretCredential()   (-secret / MSGRAPHSECRET)
  ├─► azidentity.NewClientCertificateCredential()
  │   ├─► From PFX file (-pfx + -pfxpass)
  │   └─► From Windows Cert Store (-thumbprint) → cert_windows.go
  └─► azidentity.NewBearerTokenCredential()    (-accesstoken)
```

## Shared Internal Packages

```
internal/common/
  ├── bootstrap/     — signal context (SIGINT/SIGTERM) wired via cobra PersistentPreRunE
  ├── logger/        — CSV action logs, JSON export, slog structured logger
  ├── ratelimit/     — token bucket algorithm
  ├── retry/         — exponential backoff (50ms → 10s cap), retryable error detection
  ├── security/      — credential masking (maskSecret, maskGUID)
  ├── validation/    — email, GUID, RFC3339, proxy URL, path, OData injection prevention
  └── version/       — single const Version = "3.1.6"
```

## devtools Subcommand

```
gomailtest devtools env
  ├── get      — print current MSGRAPH* env vars (secrets masked)
  ├── set      — persist MSGRAPH* vars to shell profile / user env
  └── clear    — remove MSGRAPH* vars

gomailtest devtools release
  ├── Step 1: Git status check (working tree must be clean)
  ├── Step 2: Security scan (Azure secrets, GUIDs, emails in source files)
  ├── Step 3: Version bump (update internal/common/version/version.go)
  ├── Step 4: Changelog creation (ChangeLog/{version}.md)
  ├── Step 5: git commit + push
  ├── Step 6: git tag v{version} + push tags
  └── Step 7: GitHub PR + Release via gh CLI
```

## Certificate Authentication Flow (Windows)

```
cert_windows.go (build: windows)
  └─► getCertFromStore(thumbprint)
      ├─► syscall.LoadDLL("crypt32.dll")
      ├─► CertOpenStore(CERT_SYSTEM_STORE_CURRENT_USER)
      ├─► CertFindCertificateInStore(by thumbprint)
      ├─► PFXExportCertStoreEx() → in-memory buffer only
      ├─► pkcs12.DecodeChain()
      └─► returns: crypto.PrivateKey + x509.Certificate
          (no temp files, automatic cleanup via defer)

cert_stub.go (build: !windows)
  └─► getCertFromStore() → always returns unsupported error
```

## Test Suite Architecture

```
Unit tests (go test ./...):
  ├── internal/protocols/smtp/          config_test.go, smtp_client_test.go,
  │                                     sendmail_test.go, utils_test.go
  ├── internal/protocols/jmap/          config_test.go, utils_test.go
  ├── internal/protocols/msgraph/       utils_test.go
  ├── internal/common/logger/           json_test.go
  ├── internal/common/ratelimit/        ratelimit_test.go
  ├── internal/common/security/         masking_test.go
  ├── internal/common/validation/       validation_test.go, proxy_test.go
  ├── internal/smtp/protocol/           commands_test.go, responses_test.go
  ├── internal/imap/protocol/           capabilities_test.go
  ├── internal/pop3/protocol/           capabilities_test.go, commands_test.go
  └── internal/jmap/protocol/           methods_test.go, session_test.go, types_test.go

Integration tests (go test -tags integration ./tests/integration/):
  └── tests/integration/sendmail_test.go
      └── Requires MSGRAPH* env vars (validated by scripts/check-integration-env.sh)
          └── make integration-test  (or: .\run-integration-tests.ps1)
```

## GitHub Actions CI/CD (.github/workflows/build.yml)

```
On: push tags (v*) | pull_request → main

test job (ubuntu / windows / macos):
  └── go test -v -race ./...
  └── coverage report (ubuntu only)

lint job (ubuntu, continue-on-error):
  └── golangci-lint

build job (on tag push, needs: test):
  Matrix: windows-latest (amd64), ubuntu-latest (amd64), macos-latest (arm64)
  ├── go build -ldflags="-s -w" -o bin/gomailtest[.exe] ./cmd/gomailtest
  ├── Verify binary exists
  ├── Create ZIP: bin/gomailtest[.exe] + README.md + TOOLS.md + EXAMPLES.md + LICENSE
  │   → gomailtesttool-{os}-{arch}.zip
  ├── Upload artifacts
  └── Create GitHub Release (softprops/action-gh-release)
```

## Data Flow Example: Send Email via msgraph

```
gomailtest msgraph sendmail -mailbox user@example.com -to dest@example.com -subject "Test"
          │
          ▼
internal/protocols/msgraph/cmd.go    — parse flags, validate config
          │
          ▼
internal/protocols/msgraph/auth.go   — getCredential() → azcore.TokenCredential
          │                             msgraphsdk.NewGraphServiceClientWithCredentials()
          ▼
internal/protocols/msgraph/handlers.go — handleSendMail()
  ├── createRecipients(["dest@example.com"])
  ├── createFileAttachments([]) → getAttachmentContentBase64()
  ├── build models.Message
  └── client.Users().ByUserId().SendMail().Post()
          │
          ▼
internal/common/retry/retry.go       — retryWithBackoff()
  ├── isRetryableError() → 429, 503, 504
  └── exponential backoff: 50ms → 100ms → 200ms → ... → 10s cap
          │
          ▼
internal/common/logger/csv.go        — append to %TEMP%\_msgraphtool_sendmail_{date}.csv
```

## Key Design Patterns

### 1. Cobra Subcommand Pattern

Each protocol registers a `NewCmd()` that returns a `*cobra.Command` with its own subcommands:

```go
func NewCmd() *cobra.Command {
    cmd := &cobra.Command{Use: "smtp", Short: "SMTP testing"}
    cmd.AddCommand(newTestConnectCmd())
    cmd.AddCommand(newTestStartTLSCmd())
    // ...
    return cmd
}
```

### 2. Table-Driven Tests

```go
tests := []struct {
    name     string
    input    string
    expected string
    wantErr  bool
}{ /* cases */ }
```

### 3. Retry with Exponential Backoff

```go
retryWithBackoff(ctx, maxRetries, baseDelay, operation func() error)
```

### 4. Platform-Specific Builds

```go
// cert_windows.go — //go:build windows  (Windows Certificate Store access)
// cert_stub.go    — //go:build !windows (returns unsupported error)
```

### 5. CSV Logging Pattern

Action-specific files prevent schema conflicts:

```
%TEMP%\_msgraphtool_sendmail_{date}.csv
%TEMP%\_msgraphtool_getevents_{date}.csv
%TEMP%\_smtptool_testconnect_{date}.csv
%TEMP%\_imaptool_listfolders_{date}.csv
%TEMP%\_pop3tool_listmail_{date}.csv
%TEMP%\_jmaptool_getmailboxes_{date}.csv
```

### 6. JSON Export Pattern

Export actions create date-stamped directories:

```
%TEMP%\export\{date}\
  message_1_{timestamp}.json
  message_2_{timestamp}.json
  message_search_{timestamp}.json
```

---

## Project Statistics

**Version:** 3.1.6 (Latest)
**Last Updated:** 2026-04-12

### Codebase Metrics
- **Binary:** 1 unified `gomailtest` (cobra CLI)
- **Protocol subcommands:** 5 (smtp, imap, pop3, jmap, msgraph)
- **Supported Platforms:** Windows (amd64), Linux (amd64), macOS (arm64)
- **Integration Tests:** MS Graph sendmail (tests/integration/)

### Architecture Evolution
- **v1.x:** Single msgraphtool binary
- **v2.0+:** Multi-tool suite (5 separate binaries) with shared internal packages
- **v3.0+:** Unified `gomailtest` binary with cobra subcommands; protocol logic in `internal/protocols/`; `devtools` subcommand replaces PS1 release scripts

                          ..ooOO END OOoo..
