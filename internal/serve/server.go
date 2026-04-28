package serve

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"msgraphtool/internal/common/version"
	"msgraphtool/internal/protocols/msgraph"
	"msgraphtool/internal/protocols/smtp"
)

// Server holds shared state for the HTTP serve command.
type Server struct {
	config      *Config
	smtpBase    *smtp.Config                     // nil when SMTP env vars are absent
	msgraphBase *msgraph.Config                  // nil when MSGRAPH env vars are absent
	graphClient *msgraphsdk.GraphServiceClient   // nil when msgraphBase is nil
	logger      *slog.Logger
}

// New creates a Server. smtpBase and msgraphBase may be nil when the
// corresponding credentials were not configured at startup.
func New(cfg *Config, smtpBase *smtp.Config, msgraphBase *msgraph.Config, graphClient *msgraphsdk.GraphServiceClient, logger *slog.Logger) *Server {
	return &Server{
		config:      cfg,
		smtpBase:    smtpBase,
		msgraphBase: msgraphBase,
		graphClient: graphClient,
		logger:      logger,
	}
}

// Run starts the HTTP server and blocks until ctx is cancelled or the server errors.
func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /smtp/sendmail", s.handleSMTPSendMail)
	mux.HandleFunc("POST /msgraph/sendmail", s.handleMsgraphSendMail)
	mux.HandleFunc("POST /ews/sendmail", s.handleEWSSendMail)

	addr := fmt.Sprintf("%s:%d", s.config.Listen, s.config.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.apiKeyMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 70 * time.Second, // handler timeout is 60s; extra 10s for response flush
	}

	s.logger.Info("HTTP server listening", "addr", addr)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

// apiKeyMiddleware enforces X-API-Key on all routes except /health.
func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" && subtle.ConstantTimeCompare([]byte(r.Header.Get("X-API-Key")), []byte(s.config.APIKey)) != 1 {
			writeJSON(w, http.StatusUnauthorized, apiResponse{Status: "error", Message: "missing or invalid X-API-Key"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": version.Get()})
}

// apiResponse is the standard JSON envelope for all endpoint responses.
type apiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
