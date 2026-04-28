package serve

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"msgraphtool/internal/common/logger"
	"msgraphtool/internal/common/validation"
	"msgraphtool/internal/protocols/smtp"
)

// smtpSendRequest is the JSON body for POST /smtp/sendmail.
type smtpSendRequest struct {
	To      []string `json:"to"`
	From    string   `json:"from,omitempty"` // optional override for SMTPFROM
	Subject string   `json:"subject"`
	Body    string   `json:"body,omitempty"`
}

func (s *Server) handleSMTPSendMail(w http.ResponseWriter, r *http.Request) {
	if s.smtpBase == nil {
		writeJSON(w, http.StatusServiceUnavailable, apiResponse{Status: "error", Message: "SMTP not configured (set SMTPHOST and related env vars)"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit

	var req smtpSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Status: "error", Message: "invalid JSON: " + err.Error()})
		return
	}

	if len(req.To) == 0 {
		writeJSON(w, http.StatusBadRequest, apiResponse{Status: "error", Message: "to is required"})
		return
	}
	if req.Subject == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Status: "error", Message: "subject is required"})
		return
	}
	if req.From == "" && s.smtpBase.From == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Status: "error", Message: "from is required (set SMTPFROM or provide 'from' in request body)"})
		return
	}
	for _, addr := range req.To {
		if err := validation.ValidateEmail(addr); err != nil {
			writeJSON(w, http.StatusBadRequest, apiResponse{Status: "error", Message: "invalid to address " + addr + ": " + err.Error()})
			return
		}
	}

	// Clone base config and overlay request content
	cfg := *s.smtpBase
	cfg.To = req.To
	cfg.Subject = req.Subject
	cfg.Body = req.Body
	cfg.Action = smtp.ActionSendMail
	if req.From != "" {
		cfg.From = req.From
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	csvLogger, err := logger.NewLogger(logger.LogFormatCSV, "servetool", "smtp-sendmail")
	if err != nil {
		s.logger.Warn("Could not initialise CSV logger for SMTP send", "error", err)
	} else {
		defer csvLogger.Close()
	}

	if err := smtp.SendMail(ctx, &cfg, csvLogger, s.logger); err != nil {
		s.logger.Error("SMTP sendmail failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, apiResponse{Status: "error", Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Status: "ok"})
}
