package serve

// Config holds HTTP server configuration for the serve subcommand.
type Config struct {
	Port   int    // listen port (default 8080)
	Listen string // bind address; empty string means all interfaces
	APIKey string // required X-API-Key header value
}
