# Example Usage — Microsoft Graph (gomailtest msgraph)

Comprehensive examples for Microsoft Graph API operations via `gomailtest msgraph`.

**Prerequisites:** Set authentication environment variables:

```powershell
$env:MSGRAPHTENANTID = "your-tenant-id"
$env:MSGRAPHCLIENTID = "your-client-id"
$env:MSGRAPHSECRET = "your-secret"  # or use --pfx/--thumbprint
$env:MSGRAPHMAILBOX = "user@example.com"
```

See [docs/protocols/msgraph.md](docs/protocols/msgraph.md) for all flags and authentication methods.

---

## 1. Get Calendar Events

```powershell
# Get default 3 upcoming events
gomailtest msgraph getevents

# Get 10 upcoming events
gomailtest msgraph getevents --count 10

# Get 5 events with verbose output
gomailtest msgraph getevents --count 5 --verbose
```

---

## 2. Send Email — Basic

```powershell
# Send to self (default behavior when no recipients specified)
gomailtest msgraph sendmail

# Send to specific recipient
gomailtest msgraph sendmail --to "recipient@example.com"

# Send with custom subject and body
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "Test Email" \
    --body "This is a test message"

# Send to multiple recipients (comma-separated)
gomailtest msgraph sendmail \
    --to "user1@example.com,user2@example.com" \
    --subject "Team Update"
```

---

## 3. Send Email — With CC/BCC

```powershell
# Send with CC recipients
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --cc "cc1@example.com,cc2@example.com" \
    --subject "Meeting Notes"

# Send with BCC recipients
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --bcc "bcc@example.com" \
    --subject "Confidential Update"

# Send with To, CC, and BCC
gomailtest msgraph sendmail \
    --to "primary@example.com" \
    --cc "cc1@example.com,cc2@example.com" \
    --bcc "bcc@example.com" \
    --subject "Quarterly Report" \
    --body "Please review the attached report."
```

---

## 4. Send Email — HTML Content

```powershell
# Send HTML email
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "HTML Email Test" \
    --bodyHTML "<h1>Hello</h1><p>This is an <strong>HTML</strong> email.</p>"

# Send both text and HTML (multipart MIME)
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "Multipart Email" \
    --body "This is the plain text version" \
    --bodyHTML "<h1>HTML Version</h1><p>This is the <em>HTML</em> version</p>"
```

---

## 5. Send Email — With Attachments

```powershell
# Send with single attachment
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "Document Attached" \
    --attachments "C:\Reports\report.pdf"

# Send with multiple attachments (comma-separated)
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "Multiple Files" \
    --attachments "C:\Files\doc1.pdf,C:\Files\spreadsheet.xlsx,C:\Files\image.png"

# HTML email with attachments
gomailtest msgraph sendmail \
    --to "recipient@example.com" \
    --subject "Report with Charts" \
    --bodyHTML "<h1>Monthly Report</h1><p>See attached files.</p>" \
    --attachments "C:\Reports\report.pdf,C:\Charts\chart.png"
```

---

## 6. Create Calendar Invites

```powershell
# Create invite with default subject and time (now + 1 hour)
gomailtest msgraph sendinvite

# Create invite with custom subject
gomailtest msgraph sendinvite --subject "Team Meeting"

# Create invite with specific start and end times
gomailtest msgraph sendinvite \
    --subject "Project Review" \
    --start "2026-01-15T14:00:00Z" \
    --end "2026-01-15T15:00:00Z"

# Create all-day event
gomailtest msgraph sendinvite \
    --subject "Conference Day" \
    --start "2026-02-01T00:00:00Z" \
    --end "2026-02-02T00:00:00Z"
```

---

## 7. Get Inbox Messages

```powershell
# Get default 3 newest messages
gomailtest msgraph getinbox

# Get 10 newest messages
gomailtest msgraph getinbox --count 10

# Get 20 newest messages with verbose output
gomailtest msgraph getinbox --count 20 --verbose
```

---

## 8. Using Proxy

```powershell
# Specify proxy on command line
gomailtest msgraph sendmail \
    --to "user@example.com" \
    --proxy "http://proxy.company.com:8080"

# Proxy via environment variable
$env:MSGRAPHPROXY = "http://proxy.company.com:8080"
gomailtest msgraph getevents
```

---

## 9. Verbose Output

```powershell
# Show detailed configuration, authentication, and API call information
gomailtest msgraph sendmail --to "user@example.com" --verbose
```

---

## 10. Complex Examples

### Send Formatted HTML Report with Multiple Attachments

```powershell
gomailtest msgraph sendmail \
    --to "team-lead@example.com,manager@example.com" \
    --cc "team@example.com" \
    --subject "Q1 2026 Performance Report" \
    --bodyHTML "<h1>Q1 Performance Report</h1><p>See attached metrics and analysis.</p>" \
    --attachments "C:\Reports\Q1-Metrics.xlsx,C:\Reports\Q1-Analysis.pdf" \
    --verbose
```

### Automated Monitoring Script

```powershell
# Log inbox and calendar to files
gomailtest msgraph getinbox --count 50 | Out-File -Append "C:\Logs\inbox-monitor.log"
gomailtest msgraph getevents --count 20 | Out-File -Append "C:\Logs\calendar-monitor.log"
```

---

## 11. Authentication Methods

```powershell
# Client Secret (via environment variable)
$env:MSGRAPHSECRET = "your-secret"
gomailtest msgraph getevents

# PFX Certificate File
gomailtest msgraph getevents \
    --pfx "C:\Certs\app-cert.pfx" \
    --pfxpass "MyP@ssw0rd"

# Windows Certificate Store (Thumbprint)
gomailtest msgraph getevents \
    --thumbprint "CD817B3329802E692CF30D8DDF896FE811B048AB"
```

---

## 12. Export Inbox to JSON

```powershell
# Export default 3 newest messages to JSON
gomailtest msgraph exportinbox

# Export 50 messages with verbose output
gomailtest msgraph exportinbox --count 50 --verbose
```

Output directory: `%TEMP%\export\{date}\message_{n}_{timestamp}.json`

---

## 13. Search and Export by Message ID

```powershell
# Search for specific email and export
gomailtest msgraph searchandexport \
    --messageid "<message-id@example.com>"

# With verbose output to see search details
gomailtest msgraph searchandexport \
    --messageid "<CABcD123@mail.gmail.com>" \
    --verbose
```

---

## 14. Retry Configuration

```powershell
# Custom retry settings
gomailtest msgraph getevents \
    --maxretries 5 \
    --retrydelay 1000

# Disable retries
gomailtest msgraph sendmail --to "user@example.com" --maxretries 0

# Via environment variables
$env:MSGRAPHMAXRETRIES = "5"
$env:MSGRAPHRETRYDELAY = "2500"
gomailtest msgraph getevents
```

**Retry behavior:** Exponential backoff (2s → 4s → 8s → 30s max). Retries on HTTP 429, 503, 504, network timeouts. Never retries auth failures or 4xx errors.

---

## 15. CSV Log Files

Operations are logged to `%TEMP%\_msgraphtool_{action}_{date}.csv`.

Schemas:
- **getevents**: Timestamp, Action, Status, Mailbox, Event Subject, Event ID
- **sendmail**: Timestamp, Action, Status, Mailbox, To, CC, BCC, Subject, Body Type, Attachments
- **sendinvite**: Timestamp, Action, Status, Mailbox, Subject, Start Time, End Time, Event ID
- **getinbox**: Timestamp, Action, Status, Mailbox, Subject, From, To, Received DateTime

---

## 16. Environment Variables Reference

| Flag | Environment Variable | Example |
|------|---------------------|---------|
| `--tenantid` | `MSGRAPHTENANTID` | `"tenant-id"` |
| `--clientid` | `MSGRAPHCLIENTID` | `"client-id"` |
| `--secret` | `MSGRAPHSECRET` | `"secret"` |
| `--pfx` | `MSGRAPHPFX` | `"C:\\cert.pfx"` |
| `--pfxpass` | `MSGRAPHPFXPASS` | `"password"` |
| `--thumbprint` | `MSGRAPHTHUMBPRINT` | `"CD817..."` |
| `--mailbox` | `MSGRAPHMAILBOX` | `"user@example.com"` |
| `--to` | `MSGRAPHTO` | `"user1@example.com,user2@example.com"` |
| `--cc` | `MSGRAPHCC` | `"cc@example.com"` |
| `--bcc` | `MSGRAPHBCC` | `"bcc@example.com"` |
| `--subject` | `MSGRAPHSUBJECT` | `"Email Subject"` |
| `--body` | `MSGRAPHBODY` | `"Email body"` |
| `--bodyHTML` | `MSGRAPHBODYHTML` | `"<h1>HTML</h1>"` |
| `--attachments` | `MSGRAPHATTACHMENTS` | `"file1.pdf,file2.xlsx"` |
| `--start` | `MSGRAPHSTART` | `"2026-01-15T14:00:00Z"` |
| `--end` | `MSGRAPHEND` | `"2026-01-15T15:00:00Z"` |
| `--proxy` | `MSGRAPHPROXY` | `"http://proxy:8080"` |
| `--count` | `MSGRAPHCOUNT` | `"10"` |
| `--maxretries` | `MSGRAPHMAXRETRIES` | `"5"` |
| `--retrydelay` | `MSGRAPHRETRYDELAY` | `"2000"` |
| `--messageid` | `MSGRAPHMESSAGEID` | `"<msg-id@example.com>"` |

---

## Quick Reference

```powershell
# Show version
gomailtest --version

# Show all protocols
gomailtest --help

# Show msgraph subcommands
gomailtest msgraph --help

# Test authentication
gomailtest msgraph getevents --verbose

# Send quick test email
gomailtest msgraph sendmail

# View recent inbox
gomailtest msgraph getinbox --count 10
```

---

## Tips and Best Practices

1. **Security**: Use environment variables for sensitive data — avoid passing secrets as CLI flags (visible in process list)
2. **Verbose Mode**: Use `--verbose` when troubleshooting authentication or API issues
3. **CSV Logs**: Check CSV log files for historical records of all operations
4. **Graceful Shutdown**: Press Ctrl+C to interrupt long-running operations safely (CSV logger closes cleanly)
5. **Flag Precedence**: Command-line flags override environment variables
6. **Comma Separation**: Lists (`--to`, `--cc`, `--bcc`, `--attachments`) use comma-separation; spaces are trimmed
7. **Time Format**: Calendar times use RFC3339 format (e.g., `2026-01-15T14:00:00Z`)

                          ..ooOO END OOoo..
