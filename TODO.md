# TODO

Outstanding work items for gomailtesttool. Carried over from [CODE_REVIEW.md](CODE_REVIEW.md) (the referenced IMPROVEMENTS.md no longer exists in the repo).

## Email composition features

- [x] Attachments: File attachments with MIME type detection
  - `msgraph sendmail` (`--attachments`) and `smtp sendmail` (`--attachments`), both via shared `internal/common/email.LoadAttachments`
- [x] Inline Attachments: For embedding images in HTML emails (Content-ID/`cid:` references)
  - `msgraph sendmail` and `smtp sendmail` via `--inline-attachments`, referenced from `--bodyhtml`/`--bodyHTML` as `cid:<filename>`
- [x] Custom Headers: Add custom headers via repeatable `--header "Name: Value"`
  - `msgraph sendmail` and `smtp sendmail`; protected headers (From, To, Subject, Date, Message-ID, MIME-Version, Content-Type, etc.) cannot be overridden
- [x] Multipart Messages: Support for plain text and HTML bodies
  - `msgraph sendmail` (`--body` / `--bodyHTML`) and `smtp sendmail` (`--body` / `--bodyhtml`, `multipart/alternative` when both are set)

## Configuration

- [x] `--config <file>`: YAML config file providing default flag values for every protocol/action
  - Registered once on `rootCmd` (`cmd/gomailtest/root.go`), loaded via `bootstrap.LoadConfigFile` in every subcommand's `RunE`
  - Precedence: CLI flags > env vars > `--config` file > built-in defaults (see [docs/config-file.md](docs/config-file.md))
  - `serve` mode also supported via `bootstrap.LoadConfigFileSection`: top-level keys configure the server; nested `smtp:`/`msgraph:` sections provide defaults below `SMTP*`/`MSGRAPH*` env vars
  - [ ] add priority flags on send mail (SMTP, EWS, Graph)
