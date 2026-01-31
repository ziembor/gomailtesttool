# Architecture Diagram - gomailtesttool Suite

## Overview

**gomailtesttool** is a comprehensive email infrastructure testing suite with 5 specialized CLI tools:
- **msgraphtool** - Microsoft Graph API (Exchange Online)
- **smtptool** - SMTP connectivity and TLS diagnostics
- **imaptool** - IMAP server testing with OAuth2
- **pop3tool** - POP3 server testing with OAuth2
- **jmaptool** - JMAP protocol testing

## File Structure and Dependencies

```bash
gomailtesttool/
├── cmd/                              # Command-line tools (each builds to binary)
│   ├── msgraphtool/                  # Microsoft Graph tool
│   │   ├── main.go                   # Entry point
│   │   ├── config.go                 # Configuration
│   │   ├── auth.go                   # Authentication
│   │   ├── handlers.go               # Action handlers
│   │   ├── utils.go                  # Utilities
│   │   ├── cert_windows.go           # Windows cert store (build: windows)
│   │   └── cert_stub.go              # Cross-platform stub (build: !windows)
│   │
│   ├── smtptool/                     # SMTP testing tool
│   │   ├── main.go
│   │   ├── config.go
│   │   ├── handlers.go
│   │   ├── smtp_client.go            # SMTP client logic
│   │   ├── testconnect.go            # Connectivity tests
│   │   ├── teststarttls.go           # TLS diagnostics
│   │   ├── testauth.go               # Auth mechanism tests
│   │   ├── sendmail.go               # Email sending
│   │   └── *_test.go                 # Unit tests
│   │
│   ├── imaptool/                     # IMAP testing tool
│   │   ├── main.go
│   │   ├── config.go
│   │   ├── handlers.go
│   │   ├── imap_client.go            # IMAP client logic
│   │   ├── listfolders.go            # Folder operations
│   │   ├── testconnect.go            # Connectivity tests
│   │   ├── testauth.go               # Auth tests
│   │   └── *_test.go
│   │
│   ├── pop3tool/                     # POP3 testing tool
│   │   ├── main.go
│   │   ├── config.go
│   │   ├── handlers.go
│   │   ├── pop3_client.go            # POP3 client logic
│   │   ├── listmail.go               # Message retrieval
│   │   ├── testconnect.go
│   │   └── testauth.go
│   │
│   └── jmaptool/                     # JMAP testing tool
│       ├── main.go
│       ├── config.go
│       ├── handlers.go
│       ├── jmap_client.go            # JMAP client logic
│       ├── getmailboxes.go           # Mailbox operations
│       ├── testconnect.go
│       ├── testauth.go
│       └── *_test.go
│
├── internal/                         # Shared internal packages
│   ├── common/                       # Cross-tool utilities
│   │   ├── logger/                   # CSV/JSON logging
│   │   │   ├── csv.go
│   │   │   └── json_test.go
│   │   ├── ratelimit/                # Rate limiting
│   │   │   └── ratelimit_test.go
│   │   ├── retry/                    # Retry logic
│   │   ├── security/                 # Security utilities
│   │   │   └── masking_test.go
│   │   ├── validation/               # Input validation
│   │   │   ├── proxy_test.go
│   │   │   └── validation_test.go
│   │   └── version/                  # Version management
│   │
│   ├── smtp/                         # SMTP-specific packages
│   │   ├── protocol/                 # SMTP protocol logic
│   │   │   ├── commands_test.go
│   │   │   └── responses_test.go
│   │   ├── exchange/                 # Exchange detection
│   │   └── tls/                      # TLS diagnostics
│   │
│   ├── imap/                         # IMAP-specific packages
│   │   └── protocol/
│   │       └── capabilities_test.go
│   │
│   ├── pop3/                         # POP3-specific packages
│   │   └── protocol/
│   │       ├── capabilities_test.go
│   │       └── commands_test.go
│   │
│   └── jmap/                         # JMAP-specific packages
│       └── protocol/
│           ├── methods_test.go
│           ├── session_test.go
│           └── types_test.go
│
├── src/                              # Legacy msgraphtool source (being migrated to cmd/)
│   ├── msgraphtool.go                # Main entry point
│   ├── config.go                     # Configuration
│   ├── auth.go                       # Authentication
│   ├── handlers.go                   # Action dispatcher
│   ├── handler_mail.go               # Email operations
│   ├── handler_calendar.go           # Calendar operations
│   ├── handler_search.go             # Search/export operations
│   ├── logger.go                     # Logging setup
│   ├── utils.go                      # Utilities
│   ├── version.go                    # Version info
│   ├── cert_windows.go               # Windows cert store
│   ├── cert_stub.go                  # Cross-platform stub
│   │
│   ├── *_test.go                     # Unit tests
│   ├── integration_test_tool.go      # Interactive integration tests (build: integration)
│   └── msgraphtool_integration_test.go  # Automated integration tests (build: integration)
│
├── tests/                            # Test scripts and fixtures
│   ├── README.md
│   └── Test-SendMail.ps1             # Pester test example
│
├── build-all.ps1                     # Build all tools
├── run-integration-tests.ps1         # Integration test runner
├── run-interactive-release.ps1       # Release automation
├── go.mod                            # Go module definition
└── go.sum                            # Dependency checksums
```

## Tool Architecture Overview

### Multi-Tool Build System

```bash
┌────────────────────────────────────────────────────────────────┐
│                    Build System (build-all.ps1)                 │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► go build ./cmd/msgraphtool  → bin/msgraphtool.exe
                          ├─► go build ./cmd/smtptool     → bin/smtptool.exe
                          ├─► go build ./cmd/imaptool     → bin/imaptool.exe
                          ├─► go build ./cmd/pop3tool     → bin/pop3tool.exe
                          └─► go build ./cmd/jmaptool     → bin/jmaptool.exe
```

### msgraphtool Application Flow (Microsoft Graph)

```bash
┌─────────────────────────────────────────────────────────────────┐
│                  cmd/msgraphtool/main.go                        │
│                     (Microsoft Graph Tool Entry)                 │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ├─► main()
                           │   ├─► Parse flags & environment variables (config.go)
                           │   ├─► validateConfiguration() (config.go)
                           │   ├─► setupLogger() (logger.go)
                           │   └─► Route to action handlers (handlers.go)
                           │
                           ├─► Action Handlers (dispatch based on -action flag)
                           │   │
                           │   ├─► Calendar Operations (handler_calendar.go)
                           │   │   ├─► handleGetEvents()      (-action getevents)
                           │   │   ├─► handleSendInvite()     (-action sendinvite)
                           │   │   └─► handleGetSchedule()    (-action getschedule)
                           │   │
                           │   ├─► Mail Operations (handler_mail.go)
                           │   │   ├─► handleSendMail()       (-action sendmail)
                           │   │   └─► handleGetInbox()       (-action getinbox)
                           │   │
                           │   └─► Search/Export Operations (handler_search.go)
                           │       ├─► handleExportInbox()    (-action exportinbox)
                           │       └─► handleSearchAndExport()(-action searchandexport)
                           │
                           └─► Utility Functions (utils.go)
                               ├─► showVersion()
                               ├─► generateBashCompletion()
                               └─► generatePowerShellCompletion()
```

### smtptool Application Flow

```bash
┌─────────────────────────────────────────────────────────────────┐
│                   cmd/smtptool/main.go                          │
│                      (SMTP Testing Tool Entry)                   │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ├─► main()
                           │   ├─► Parse flags (config.go)
                           │   ├─► validateConfiguration()
                           │   └─► Route to action handlers
                           │
                           └─► Action Handlers (handlers.go)
                               ├─► handleTestConnect()    (testconnect.go)
                               ├─► handleTestSTARTTLS()   (teststarttls.go)
                               ├─► handleTestAuth()       (testauth.go)
                               └─► handleSendMail()       (sendmail.go)
```

### imaptool Application Flow

```bash
┌─────────────────────────────────────────────────────────────────┐
│                   cmd/imaptool/main.go                          │
│                      (IMAP Testing Tool Entry)                   │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           └─► Action Handlers (handlers.go)
                               ├─► handleTestConnect()    (testconnect.go)
                               ├─► handleTestAuth()       (testauth.go)
                               └─► handleListFolders()    (listfolders.go)
```

### pop3tool & jmaptool Application Flow

Similar patterns with protocol-specific handlers:
- **pop3tool**: testconnect, testauth, listmail
- **jmaptool**: testconnect, testauth, getmailboxes

## Shared Business Logic

### Internal Packages (Shared Across All Tools)

```bash
┌────────────────────────────────────────────────────────────────┐
│                    internal/ Package Structure                  │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► internal/common/              # Cross-tool utilities
                          │   ├─► logger/                   # CSV/JSON logging
                          │   │   └─► Structured logging for all tools
                          │   │
                          │   ├─► validation/               # Input validation
                          │   │   ├─► Email validation
                          │   │   ├─► GUID validation
                          │   │   ├─► Proxy URL validation
                          │   │   └─► Security checks
                          │   │
                          │   ├─► security/                 # Security utilities
                          │   │   └─► Credential masking
                          │   │
                          │   ├─► retry/                    # Retry logic
                          │   │   ├─► Exponential backoff
                          │   │   └─► Retryable error detection
                          │   │
                          │   ├─► ratelimit/               # Rate limiting
                          │   │   └─► Token bucket algorithm
                          │   │
                          │   └─► version/                  # Version management
                          │
                          ├─► internal/smtp/               # SMTP-specific
                          │   ├─► protocol/                # Protocol commands
                          │   ├─► exchange/                # Exchange detection
                          │   └─► tls/                     # TLS diagnostics
                          │
                          ├─► internal/imap/protocol/      # IMAP capabilities
                          ├─► internal/pop3/protocol/      # POP3 commands
                          └─► internal/jmap/protocol/      # JMAP session/methods
```

### msgraphtool Authentication & Client Setup (src/)

```bash
┌────────────────────────────────────────────────────────────────┐
│              msgraphtool Authentication Layer                   │
│                       (src/auth.go)                             │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► setupGraphClient()
                          │   └─► getCredential()
                          │       ├─► azidentity.NewClientSecretCredential()
                          │       │   (uses -secret flag or MSGRAPHSECRET env)
                          │       │
                          │       ├─► azidentity.NewClientCertificateCredential()
                          │       │   ├─► From PFX file (-pfx + -pfxpass)
                          │       │   │   └─► pkcs12.DecodeChain()
                          │       │   │
                          │       │   └─► From Windows Cert Store (-thumbprint)
                          │       │       └─► getCertFromStore() (cert_windows.go)
                          │       │           └─► Windows CryptoAPI (crypt32.dll)
                          │       │
                          │       ├─► azidentity.NewBearerTokenCredential()
                          │       │   (uses -accesstoken flag)
                          │       │
                          │       └─► Returns: azcore.TokenCredential
                          │
                          └─► Uses internal/common/retry package
                              ├─► Exponential backoff (50ms → 10s cap)
                              └─► Retryable error detection (429, 503, 504)
```

### msgraphtool Core Graph API Operations

```bash
┌────────────────────────────────────────────────────────────────┐
│              Microsoft Graph API Layer (src/)                   │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► Calendar Operations (handler_calendar.go)
                          │   │
                          │   ├─► handleGetEvents()
                          │   │   └─► client.Users().ByUserId().Events().Get()
                          │   │       └─► Returns: []models.Event
                          │   │
                          │   ├─► handleSendInvite()
                          │   │   ├─► parseFlexibleTime()
                          │   │   ├─► createRecipients()
                          │   │   └─► client.Users().ByUserId().Events().Post()
                          │   │
                          │   └─► handleGetSchedule()
                          │       └─► client.Users().GetSchedule().Post()
                          │           └─► Returns availability information
                          │
                          ├─► Mail Operations (handler_mail.go)
                          │   │
                          │   ├─► handleSendMail()
                          │   │   ├─► createRecipients()
                          │   │   ├─► createFileAttachments()
                          │   │   │   └─► getAttachmentContentBase64()
                          │   │   └─► client.Users().ByUserId().SendMail().Post()
                          │   │
                          │   └─► handleGetInbox()
                          │       └─► client.Users().ByUserId().Messages().Get()
                          │           └─► Returns: []models.Message
                          │
                          └─► Search/Export Operations (handler_search.go)
                              │
                              ├─► handleExportInbox()
                              │   ├─► client.Users().ByUserId().Messages().Get()
                              │   ├─► Create date-stamped dir (%TEMP%\export\{date})
                              │   └─► Export each message to individual JSON file
                              │
                              └─► handleSearchAndExport()
                                  ├─► client.Users().ByUserId().Messages().Get()
                                  │   └─► Filter by InternetMessageId
                                  │   └─► Security: Validates Message-ID format
                                  └─► Export matching message to JSON file
```

### SMTP Tool Operations

```bash
┌────────────────────────────────────────────────────────────────┐
│                  SMTP Operations (cmd/smtptool/)                │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► handleTestConnect() (testconnect.go)
                          │   └─► Tests basic TCP connectivity
                          │
                          ├─► handleTestSTARTTLS() (teststarttls.go)
                          │   ├─► TLS handshake analysis
                          │   ├─► Certificate validation
                          │   ├─► Cipher strength assessment
                          │   └─► Exchange server detection
                          │
                          ├─► handleTestAuth() (testauth.go)
                          │   └─► Tests auth mechanisms (PLAIN, LOGIN, CRAM-MD5, XOAUTH2)
                          │
                          └─► handleSendMail() (sendmail.go)
                              └─► Sends test email via SMTP
```

### IMAP/POP3/JMAP Tool Operations

Each tool follows similar patterns with protocol-specific operations:

- **imaptool**: testconnect, testauth, listfolders (with XOAUTH2 support)
- **pop3tool**: testconnect, testauth, listmail (with APOP and XOAUTH2)
- **jmaptool**: testconnect, testauth, getmailboxes (HTTP-based JMAP)

### Validation & Helper Functions

```bash
┌────────────────────────────────────────────────────────────────┐
│                 Validation & Utilities (msgraphtool)            │
│                         (src/utils.go)                          │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► Data Transformation
                          │   ├─► createRecipients()         [100% coverage]
                          │   ├─► createFileAttachments()    [95.2% coverage]
                          │   └─► getAttachmentContentBase64()[100% coverage]
                          │
                          ├─► Security & Masking
                          │   ├─► maskSecret()               [100% coverage]
                          │   └─► maskGUID()                 [100% coverage]
                          │
                          ├─► Helper Functions
                          │   ├─► Int32Ptr()                 [100% coverage]
                          │   ├─► enrichGraphAPIError()
                          │   ├─► parseFlexibleTime()        [100% coverage]
                          │   └─► Shell completion generators
                          │
                          └─► Configuration & Logging (config.go, logger.go)
                              ├─► validateConfiguration()
                              ├─► setupLogger()
                              └─► parseLogLevel()
```

```bash
┌────────────────────────────────────────────────────────────────┐
│            Shared Validation (internal/common/validation/)      │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► Email Validation
                          │   ├─► validateEmail()            [100% coverage]
                          │   └─► validateEmails()           [100% coverage]
                          │
                          ├─► Identifier Validation
                          │   ├─► validateGUID()             [100% coverage]
                          │   └─► validateFilePath()         [Tested]
                          │
                          ├─► Time Validation
                          │   └─► validateRFC3339Time()      [100% coverage]
                          │
                          ├─► Proxy Validation
                          │   └─► validateProxyURL()         [Tested]
                          │
                          └─► Security Validation
                              ├─► Message-ID format validation
                              ├─► OData injection prevention
                              └─► Path traversal prevention
```

## Test Suite Architecture

```bash
┌────────────────────────────────────────────────────────────────┐
│                        Test Structure                           │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► msgraphtool Unit Tests (src/*_test.go)
                          │   ├─► auth_test.go               # Authentication tests
                          │   ├─► config_test.go             # Configuration tests
                          │   ├─► handlers_test.go           # Handler tests
                          │   ├─► logger_test.go             # Logger tests
                          │   ├─► utils_test.go              # Utility tests
                          │   ├─► csvlogger_test.go          # CSV logging tests
                          │   └─► validation_security_test.go# Security validation
                          │       ├─► Message-ID injection tests
                          │       ├─► OData injection prevention
                          │       └─► Path traversal prevention
                          │
                          ├─► msgraphtool Integration Tests (src/)
                          │   │   ├─► //go:build integration
                          │   │   │
                          │   │   ├─► integration_test_tool.go
                          │   │   │   ├─► Interactive test tool with prompts
                          │   │   │   ├─► 5 test scenarios (events, mail, invite, inbox, schedule)
                          │   │   │   ├─► Auto-confirm mode (MSGRAPH_AUTO_CONFIRM env)
                          │   │   │   └─► Pretty formatted output
                          │   │   │
                          │   │   └─► msgraphtool_integration_test.go
                          │   │       ├─► Automated Go test suite
                          │   │       ├─► 11 integration tests total:
                          │   │       │   ├─► TestIntegration_Prerequisites
                          │   │       │   ├─► TestIntegration_GraphClientCreation
                          │   │       │   ├─► TestIntegration_ListEvents
                          │   │       │   ├─► TestIntegration_ListInbox
                          │   │       │   ├─► TestIntegration_CheckAvailability
                          │   │       │   ├─► TestIntegration_SendEmail (write)
                          │   │       │   ├─► TestIntegration_CreateCalendarEvent (write)
                          │   │       │   ├─► TestIntegration_ExportInbox
                          │   │       │   ├─► TestIntegration_SearchAndExport
                          │   │       │   ├─► TestIntegration_SearchAndExport_InvalidMessageID
                          │   │       │   └─► TestIntegration_ValidateConfiguration
                          │   │       └─► Requires: MSGRAPH_INTEGRATION_WRITE=true for writes
                          │
                          ├─► smtptool Unit Tests (cmd/smtptool/*_test.go)
                          │   ├─► config_test.go             # Config validation
                          │   ├─► smtp_client_test.go        # SMTP client logic
                          │   ├─► sendmail_test.go           # Email sending
                          │   └─► utils_test.go              # Utilities
                          │
                          ├─► imaptool Unit Tests (cmd/imaptool/*_test.go)
                          │   ├─► config_test.go
                          │   └─► utils_test.go
                          │
                          ├─► jmaptool Unit Tests (cmd/jmaptool/*_test.go)
                          │   ├─► config_test.go
                          │   └─► utils_test.go
                          │
                          ├─► msgraphtool Unit Tests (cmd/msgraphtool/*_test.go)
                          │   └─► utils_test.go
                          │
                          └─► Internal Package Tests (internal/*/*_test.go)
                              ├─► common/logger/json_test.go
                              ├─► common/ratelimit/ratelimit_test.go
                              ├─► common/security/masking_test.go
                              ├─► common/validation/proxy_test.go
                              ├─► common/validation/validation_test.go
                              ├─► smtp/protocol/commands_test.go
                              ├─► smtp/protocol/responses_test.go
                              ├─► imap/protocol/capabilities_test.go
                              ├─► pop3/protocol/capabilities_test.go
                              ├─► pop3/protocol/commands_test.go
                              ├─► jmap/protocol/methods_test.go
                              ├─► jmap/protocol/session_test.go
                              └─► jmap/protocol/types_test.go
```

### Integration Test Execution

```bash
┌────────────────────────────────────────────────────────────────┐
│            Integration Test Execution Methods                   │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► Method 1: PowerShell Runner
                          │   └─► .\run-integration-tests.ps1
                          │       ├─► Validates environment variables
                          │       ├─► Builds integration_test_tool.exe
                          │       ├─► Runs interactive tests
                          │       └─► Supports: -SetEnv, -ShowEnv, -ClearEnv, -AutoConfirm
                          │
                          ├─► Method 2: Interactive Tool
                          │   └─► go run -tags=integration ./src/integration_test_tool.go
                          │       ├─► User prompts before write operations
                          │       ├─► Pretty formatted output
                          │       └─► Pass/fail summary
                          │
                          └─► Method 3: Automated Tests
                              └─► go test -tags=integration -v ./src
                                  ├─► Standard Go test output
                                  ├─► CI/CD friendly
                                  └─► Requires MSGRAPH_INTEGRATION_WRITE=true for writes
```

## Certificate Authentication Flow (Windows)

```bash
┌────────────────────────────────────────────────────────────────┐
│              Windows Certificate Store Integration              │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          └─► cert_windows.go (Windows only)
                              │
                              ├─► getCertFromStore()
                              │   ├─► syscall.LoadDLL("crypt32.dll")
                              │   ├─► CertOpenStore(CERT_SYSTEM_STORE_CURRENT_USER)
                              │   ├─► CertFindCertificateInStore(by thumbprint)
                              │   ├─► PFXExportCertStoreEx() → memory buffer
                              │   ├─► pkcs12.DecodeChain()
                              │   └─► Returns: crypto.PrivateKey + x509.Certificate
                              │
                              └─► Uses Windows CryptoAPI
                                  ├─► No temporary files created
                                  ├─► Certificate extracted to memory only
                                  └─► Automatic cleanup via defer
```

## Release Automation

### Interactive Release Process (run-interactive-release.ps1)

```bash
┌────────────────────────────────────────────────────────────────┐
│              Interactive Release Script Workflow                │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► Step 1: Git Status Check
                          │   └─► Ensures working tree is clean
                          │
                          ├─► Step 2: Security Scan for Secrets
                          │   ├─► Scan patterns:
                          │   │   ├─► Azure AD Client Secrets ([a-zA-Z0-9~_-]{34,})
                          │   │   ├─► GUIDs/UUIDs (standard format)
                          │   │   ├─► Email addresses (non-example domains)
                          │   │   └─► API Keys (access_token, secret_key, etc.)
                          │   │
                          │   ├─► Scanned files:
                          │   │   ├─► ChangeLog/*.md
                          │   │   ├─► *.md (root level)
                          │   │   ├─► src/*.go
                          │   │   └─► cmd/*/*.go
                          │   │
                          │   ├─► Smart filtering:
                          │   │   ├─► Skip EXAMPLES, README, CLAUDE
                          │   │   ├─► Skip placeholders (xxx, yyy, example.com)
                          │   │   └─► Skip known safe patterns
                          │   │
                          │   └─► Blocks release if secrets detected
                          │
                          ├─► Step 3: Version Management
                          │   ├─► Validate version format (2.x.y - major locked at 2)
                          │   ├─► Prompt for new version
                          │   └─► Update version across codebase
                          │
                          ├─► Step 4: Changelog Creation
                          │   ├─► Interactive prompts for:
                          │   │   ├─► Added features
                          │   │   ├─► Changed features
                          │   │   ├─► Fixed bugs
                          │   │   └─► Security updates
                          │   └─► Create ChangeLog/{version}.md
                          │
                          ├─► Steps 5-7: Git Operations
                          │   ├─► git commit (formatted message with Claude credit)
                          │   ├─► git push origin {branch}
                          │   └─► git tag v{version} && git push origin --tags
                          │
                          └─► Steps 8-12: Optional GitHub Integration
                              ├─► Create Pull Request (via gh CLI)
                              ├─► Monitor GitHub Actions workflow
                              └─► Open releases page
```

### Build Process (build-all.ps1)

```bash
┌────────────────────────────────────────────────────────────────┐
│                   Multi-Tool Build Process                      │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► Clean previous builds
                          │   └─► Remove bin/ directory
                          │
                          ├─► Create bin/ directory
                          │
                          ├─► Build all tools in parallel
                          │   ├─► go build -ldflags="-s -w" ./cmd/msgraphtool
                          │   ├─► go build -ldflags="-s -w" ./cmd/smtptool
                          │   ├─► go build -ldflags="-s -w" ./cmd/imaptool
                          │   ├─► go build -ldflags="-s -w" ./cmd/pop3tool
                          │   └─► go build -ldflags="-s -w" ./cmd/jmaptool
                          │       └─► Output: bin/*.exe (Windows)
                          │
                          └─► Verify builds
                              └─► Display file sizes and checksums
```

### GitHub Actions CI/CD (.github/workflows/build.yml)

```bash
┌────────────────────────────────────────────────────────────────┐
│                   CI/CD Pipeline (GitHub Actions)               │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ├─► Test Job (All Platforms)
                          │   ├─► Matrix: ubuntu-latest, windows-latest, macos-latest
                          │   ├─► go test -v -race ./...
                          │   └─► Coverage report (Ubuntu only)
                          │
                          ├─► Lint Job (Ubuntu)
                          │   └─► golangci-lint (continue-on-error)
                          │
                          └─► Build Job (On Tag Push)
                              ├─► Matrix: Windows (amd64), Linux (amd64), macOS (arm64)
                              ├─► Build all 5 tools per platform
                              ├─► Create ZIP archives:
                              │   ├─► gomailtesttool-windows-amd64.zip
                              │   ├─► gomailtesttool-linux-amd64.zip
                              │   └─► gomailtesttool-macos-arm64.zip
                              ├─► Upload artifacts
                              └─► Create GitHub Release
                                  └─► Attach ZIP files with full documentation
```

## Data Flow Example: Send Email with Attachments

```bash
User Command:
  msgraphtool.exe -action sendmail -to "user@example.com"
    -subject "Test" -body "Hello" -attachment "file.pdf"
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. main() - Parse flags & validate configuration                │
│    └─► validateConfiguration() checks all required fields       │
└─────────────────────┬───────────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. setupGraphClient() - Authenticate                            │
│    ├─► getCredential() → azcore.TokenCredential                 │
│    └─► msgraphsdk.NewGraphServiceClientWithCredentials()        │
└─────────────────────┬───────────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. sendEmail() - Build and send message                         │
│    ├─► createRecipients(["user@example.com"]) → []Recipient    │
│    ├─► createFileAttachments(["file.pdf"])                      │
│    │   └─► getAttachmentContentBase64() → base64 string         │
│    ├─► Build models.Message object                              │
│    └─► client.Users().ByUserId().SendMail().Post()              │
└─────────────────────┬───────────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. retryWithBackoff() - Handle transient failures               │
│    ├─► isRetryableError() checks status codes (429, 503, 504)  │
│    └─► Exponential backoff: 50ms → 100ms → 200ms → ... → 10s   │
└─────────────────────┬───────────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 5. CSV Logging - Record operation result                        │
│    └─► %TEMP%\_msgraphtool_sendmail_2026-01-31.csv              │
│        Timestamp, Action, Status, Mailbox, To, CC, BCC, Subject │
└─────────────────────────────────────────────────────────────────┘
```

## Test Coverage Overview

```bash
┌──────────────────────────────────────────────────────────────────┐
│                    Test Coverage by Component                    │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  msgraphtool (src/):                                             │
│    ✅ Unit Tests: Comprehensive coverage of utilities            │
│       • Data transformation (createRecipients, attachments)      │
│       • Validation (email, GUID, time, configuration)            │
│       • Security (masking, injection prevention)                 │
│       • Helpers (Int32Ptr, parseFlexibleTime, etc.)              │
│    ✅ Integration Tests: 11 automated tests + interactive tool   │
│       • Real Graph API operations                                │
│       • Authentication testing                                   │
│       • Write protection (MSGRAPH_INTEGRATION_WRITE flag)        │
│                                                                  │
│  smtptool (cmd/smtptool/):                                       │
│    ✅ Unit Tests: Config, SMTP client, send mail logic           │
│    ✅ Protocol Tests: internal/smtp/protocol/*_test.go           │
│                                                                  │
│  imaptool (cmd/imaptool/):                                       │
│    ✅ Unit Tests: Config, utilities                              │
│    ✅ Protocol Tests: internal/imap/protocol/*_test.go           │
│                                                                  │
│  pop3tool (cmd/pop3tool/):                                       │
│    ✅ Protocol Tests: internal/pop3/protocol/*_test.go           │
│                                                                  │
│  jmaptool (cmd/jmaptool/):                                       │
│    ✅ Unit Tests: Config, utilities                              │
│    ✅ Protocol Tests: internal/jmap/protocol/*_test.go           │
│                                                                  │
│  Shared (internal/common/):                                      │
│    ✅ Logger tests (CSV/JSON)                                    │
│    ✅ Validation tests (email, GUID, proxy)                      │
│    ✅ Security tests (masking)                                   │
│    ✅ Rate limiting tests                                        │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## Key Design Patterns

### 1. Table-Driven Tests

All unit tests use the table-driven pattern for maintainability:

```go
tests := []struct {
    name     string
    input    string
    expected string
    wantErr  bool
}{ /* test cases */ }
```

### 2. Config Struct Pattern

Centralized configuration simplifies function signatures:

```go
type Config struct { /* all configuration */ }
func sendEmail(ctx context.Context, client *msgraphsdk.GraphServiceClient,
               cfg *Config) error
```

### 3. Retry with Exponential Backoff

Handles transient failures gracefully:

```go
retryWithBackoff(ctx, maxRetries, baseDelay, operation func() error)
```

### 4. Platform-Specific Builds

Build tags enable Windows-specific features while maintaining cross-platform support for **Windows, Linux, and macOS**:

```go
// msgraphtool only:
// cert_windows.go - //go:build windows (Windows Certificate Store access)
// cert_stub.go    - //go:build !windows (Linux/macOS stub)

// Integration tests:
// integration_test_tool.go        - //go:build integration
// msgraphtool_integration_test.go - //go:build integration
```

**GitHub Actions Workflow** builds all 5 tools for all three platforms:

- `gomailtesttool-windows-amd64.zip` - All tools for Windows (.exe)
- `gomailtesttool-linux-amd64.zip` - All tools for Linux (ELF)
- `gomailtesttool-macos-arm64.zip` - All tools for macOS (Mach-O, Apple Silicon)

**Each ZIP contains:**
- msgraphtool.exe (or msgraphtool on Linux/macOS)
- smtptool.exe
- imaptool.exe
- pop3tool.exe
- jmaptool.exe
- Complete documentation (README.md, tool-specific READMEs, EXAMPLES.md, LICENSE)

**Note:** The `-thumbprint` authentication method (Windows Certificate Store) is only available in msgraphtool on Windows. Linux and macOS users should use `-secret` or `-pfx` authentication.

### 5. CSV Logging Pattern

Action-specific CSV files prevent schema conflicts (all tools use `internal/common/logger`):

```bash
# msgraphtool logs
_msgraphtool_sendmail_2026-01-31.csv
_msgraphtool_getevents_2026-01-31.csv
_msgraphtool_getinbox_2026-01-31.csv
_msgraphtool_exportinbox_2026-01-31.csv
_msgraphtool_searchandexport_2026-01-31.csv

# smtptool logs
_smtptool_testconnect_2026-01-31.csv
_smtptool_teststarttls_2026-01-31.csv
_smtptool_sendmail_2026-01-31.csv

# imaptool, pop3tool, jmaptool logs
_imaptool_listfolders_2026-01-31.csv
_pop3tool_listmail_2026-01-31.csv
_jmaptool_getmailboxes_2026-01-31.csv
```

### 6. JSON Export Pattern (msgraphtool v1.21.0+)

msgraphtool export actions create date-stamped directories with individual JSON files:

```bash
%TEMP%\export\2026-01-31\
├── message_1_2026-01-31T10-30-45.json
├── message_2_2026-01-31T10-25-12.json
├── message_3_2026-01-31T09-58-03.json
└── message_search_2026-01-31T11-15-30.json (from searchandexport)
```

Each JSON file contains the complete message structure including:
- Headers (from, to, cc, bcc, subject, date)
- Body (HTML and text versions)
- Metadata (message ID, conversation ID, importance)
- Flags (isRead, isDraft, hasAttachments)

---

## Project Statistics

**Version:** 2.6.11 (Latest)
**Last Updated:** 2026-01-31

### Codebase Metrics
- **Total Lines of Go Code:** ~24,920
- **Test Code:** ~10,182 lines (41% of codebase)
- **Tools:** 5 (msgraphtool, smtptool, imaptool, pop3tool, jmaptool)
- **Integration Tests:** 11 (msgraphtool only)
- **Supported Platforms:** Windows, Linux, macOS

### Test Coverage
- **msgraphtool (src/):** Unit tests + Integration tests
- **Unit Tests:** Data transformation, validation, security, helpers
- **Integration Tests:** Real API calls to Microsoft Graph
- **Protocol Tests:** SMTP, IMAP, POP3, JMAP protocol logic

### Architecture Evolution
- **v1.x:** Single msgraphtool binary
- **v2.0+:** Multi-tool suite with shared internal packages
- **v2.6+:** Organized build output (bin/), comprehensive CI/CD

                          ..ooOO END OOoo..
