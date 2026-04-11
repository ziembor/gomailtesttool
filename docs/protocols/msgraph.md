# Microsoft Graph Protocol — gomailtest

Exchange Online mailbox operations via Microsoft Graph API: send emails, manage calendar events, export inbox.

> **Legacy name:** `msgraphtool`. The `msgraphtool` binary still works as a backward-compatibility shim. It will be removed in v3.1.

## Quick Start

```powershell
# Set authentication via environment variables
$env:MSGRAPHTENANTID = "your-tenant-id"
$env:MSGRAPHCLIENTID = "your-client-id"
$env:MSGRAPHSECRET = "your-secret"
$env:MSGRAPHMAILBOX = "user@example.com"

# Get calendar events
gomailtest msgraph getevents

# Send test email
gomailtest msgraph sendmail --to "recipient@example.com"

# Get inbox messages
gomailtest msgraph getinbox --count 10
```

## Actions

### getevents — Retrieve Calendar Events

```powershell
gomailtest msgraph getevents
gomailtest msgraph getevents --count 10
gomailtest msgraph getevents --count 5 --verbose
```

### sendmail — Send Email Messages

```powershell
# Send to self (default)
gomailtest msgraph sendmail

# Send to specific recipient
gomailtest msgraph sendmail --to "recipient@example.com"

# Custom subject and body
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "Test Email" \
    --body "This is a test message"

# Multiple recipients with CC and BCC
gomailtest msgraph sendmail \
    --to "user1@example.com,user2@example.com" \
    --cc "cc@example.com" \
    --bcc "bcc@example.com" \
    --subject "Team Update"

# HTML email
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "HTML Email" \
    --bodyHTML "<h1>Hello</h1><p>This is an <strong>HTML</strong> email.</p>"

# With attachments
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "Report Attached" \
    --attachments "C:\Reports\report.pdf,C:\Data\spreadsheet.xlsx"
```

### sendinvite — Create Calendar Invitations

```powershell
gomailtest msgraph sendinvite --subject "Team Meeting"

gomailtest msgraph sendinvite \
    --subject "Project Review" \
    --start "2026-01-15T14:00:00Z" \
    --end "2026-01-15T15:00:00Z"
```

### getinbox — Retrieve Inbox Messages

```powershell
gomailtest msgraph getinbox
gomailtest msgraph getinbox --count 20
```

### getschedule — Check Recipient Availability

```powershell
gomailtest msgraph getschedule --to "colleague@example.com"
```

### exportinbox — Export Inbox to JSON

```powershell
gomailtest msgraph exportinbox --count 50
```

Output goes to `%TEMP%\export\{date}\message_{n}_{timestamp}.json`.

### searchandexport — Search by Message ID

```powershell
gomailtest msgraph searchandexport --messageid "<message-id@example.com>"
```

## Flags

### Persistent (all subcommands)

| Flag | Description | Environment Variable | Default |
|------|-------------|---------------------|---------|
| `--tenantid` | Azure AD Tenant ID (GUID) | `MSGRAPHTENANTID` | — |
| `--clientid` | Application (Client) ID (GUID) | `MSGRAPHCLIENTID` | — |
| `--mailbox` | Target user email address | `MSGRAPHMAILBOX` | — |
| `--secret` | Client Secret | `MSGRAPHSECRET` | — |
| `--pfx` | Path to .pfx certificate file | `MSGRAPHPFX` | — |
| `--pfxpass` | Password for .pfx certificate | `MSGRAPHPFXPASS` | — |
| `--thumbprint` | Certificate thumbprint (Windows only) | `MSGRAPHTHUMBPRINT` | — |
| `--bearertoken` | Pre-obtained Bearer token | `MSGRAPHBEARERTOKEN` | — |
| `--proxy` | HTTP/HTTPS proxy URL | `MSGRAPHPROXY` | — |
| `--maxretries` | Maximum retry attempts | `MSGRAPHMAXRETRIES` | 3 |
| `--retrydelay` | Retry delay (milliseconds) | `MSGRAPHRETRYDELAY` | 2000 |
| `--count` | Number of items to retrieve | `MSGRAPHCOUNT` | 3 |
| `--verbose` | Enable verbose output | — | false |
| `--loglevel` | Log level: DEBUG, INFO, WARN, ERROR | `MSGRAPHLOGLEVEL` | INFO |
| `--logformat` | Log file format: csv, json | `MSGRAPHLOGFORMAT` | csv |

### Email flags (sendmail / sendinvite)

| Flag | Description | Environment Variable |
|------|-------------|---------------------|
| `--to` | Comma-separated TO recipients | `MSGRAPHTO` |
| `--cc` | Comma-separated CC recipients | `MSGRAPHCC` |
| `--bcc` | Comma-separated BCC recipients | `MSGRAPHBCC` |
| `--subject` | Email subject | `MSGRAPHSUBJECT` |
| `--body` | Email body text | `MSGRAPHBODY` |
| `--bodyHTML` | Email body HTML | `MSGRAPHBODYHTML` |
| `--body-template` | Path to HTML template file | `MSGRAPHBODYTEMPLATE` |
| `--attachments` | Comma-separated file paths | `MSGRAPHATTACHMENTS` |
| `--start` | Start time (RFC3339) | `MSGRAPHSTART` |
| `--end` | End time (RFC3339) | `MSGRAPHEND` |
| `--messageid` | Internet Message ID | `MSGRAPHMESSAGEID` |

## Authentication Methods

Provide exactly one method. Mutually exclusive:

### Client Secret

```powershell
gomailtest msgraph getevents \
    --tenantid "..." --clientid "..." \
    --secret "your-client-secret" \
    --mailbox "user@example.com"
```

### PFX Certificate

```powershell
gomailtest msgraph getevents \
    --tenantid "..." --clientid "..." \
    --pfx "C:\Certs\app-cert.pfx" --pfxpass "certificate-password" \
    --mailbox "user@example.com"
```

### Windows Certificate Store (Thumbprint)

```powershell
gomailtest msgraph getevents \
    --tenantid "..." --clientid "..." \
    --thumbprint "CD817B3329802E692CF30D8DDF896FE811B048AB" \
    --mailbox "user@example.com"
```

### Bearer Token

```powershell
gomailtest msgraph getevents \
    --tenantid "..." --clientid "..." \
    --bearertoken "eyJ0eXAi..." \
    --mailbox "user@example.com"
```

## Environment Variables

```powershell
$env:MSGRAPHTENANTID = "tenant-id"
$env:MSGRAPHCLIENTID = "client-id"
$env:MSGRAPHSECRET = "your-secret"
$env:MSGRAPHMAILBOX = "user@example.com"

gomailtest msgraph sendmail --to "recipient@example.com"
```

## CSV Logging

Operations are logged to `%TEMP%\_msgraphtool_{action}_{date}.csv`.

## Retry Configuration

```powershell
# Custom retry settings
gomailtest msgraph getevents --maxretries 5 --retrydelay 3000

# Disable retries
gomailtest msgraph sendmail --maxretries 0
```

Retry uses exponential backoff: 2s → 4s → 8s → 16s → 30s (capped). Retries on HTTP 429, 503, 504, network timeouts. Never retries authentication failures or 4xx errors.

## Required Azure AD Permissions

| Action | Permission |
|--------|-----------|
| sendmail | `Mail.Send` |
| getevents, sendinvite | `Calendars.ReadWrite` |
| getinbox, exportinbox, searchandexport | `Mail.Read` |
| getschedule | `Calendars.Read` |

## Related Documentation

- [BUILD.md](../../BUILD.md) — Build instructions
- [EXAMPLES.md](../../EXAMPLES.md) — Extended usage examples
- [TROUBLESHOOTING.md](../../TROUBLESHOOTING.md) — Common issues
- [SECURITY.md](../../SECURITY.md) — Security policy

                          ..ooOO END OOoo..
