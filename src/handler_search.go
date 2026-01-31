package main

import (
	"context"
	"fmt"
	"strings"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// searchAndExport searches for a message by Internet Message ID and exports it to JSON.
// Provides OData injection protection through input validation and escaping.
func searchAndExport(ctx context.Context, client *msgraphsdk.GraphServiceClient, mailbox string, messageID string, config *Config, logger *CSVLogger) error {
	// Configure request with filter
	// Note: We search the whole mailbox (Messages endpoint), not just Inbox
	// SECURITY: Escape single quotes for OData filter (defense-in-depth)
	// Even though validateMessageID() blocks quotes, we escape as an additional safeguard
	escapedMessageID := strings.ReplaceAll(messageID, "'", "''")
	filter := fmt.Sprintf("internetMessageId eq '%s'", escapedMessageID)
	requestConfig := &users.ItemMessagesRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.ItemMessagesRequestBuilderGetQueryParameters{
			Filter: &filter,
			Select: []string{"id", "internetMessageId", "subject", "receivedDateTime", "from", "toRecipients", "ccRecipients", "bccRecipients", "body", "hasAttachments"},
		},
	}

	logVerbose(config.VerboseMode, "Calling Graph API: GET /users/%s/messages?$filter=%s", mailbox, filter)

	// Execute API call with retry logic
	var getValueFunc func() []models.Messageable
	err := retryWithBackoff(ctx, config.MaxRetries, config.RetryDelay, func() error {
		apiResult, apiErr := client.Users().ByUserId(mailbox).Messages().Get(ctx, requestConfig)
		if apiErr == nil {
			getValueFunc = apiResult.GetValue
		}
		return apiErr
	})

	if err != nil {
		enrichedErr := enrichGraphAPIError(err, logger, "searchAndExport")
		return fmt.Errorf("error searching message for %s: %w", mailbox, enrichedErr)
	}

	messages := getValueFunc()
	messageCount := len(messages)

	logVerbose(config.VerboseMode, "API response received: %d messages", messageCount)

	if messageCount == 0 {
		if config.OutputFormat == "json" {
			printJSON([]interface{}{}) // Empty array
		} else {
			fmt.Printf("No message found with Internet Message ID: %s\n", messageID)
		}
		if logger != nil {
			logger.WriteRow([]string{ActionSearchAndExport, StatusSuccess, mailbox, "Message not found", messageID})
		}
		return nil
	}

	// Print JSON output if requested
	if config.OutputFormat == "json" {
		printJSON(formatMessagesOutput(messages))
	}

	// Create export directory
	exportDir, err := createExportDir()
	if err != nil {
		return err
	}

	if config.OutputFormat != "json" {
		fmt.Printf("Export directory: %s\n", exportDir)
	}

	// Export found messages (usually 1, but duplicates technically possible in some scenarios)
	for _, message := range messages {
		if err := exportMessageToJSON(message, exportDir, config); err != nil {
			return fmt.Errorf("failed to export message: %w", err)
		}
		if config.OutputFormat != "json" {
			fmt.Printf("Successfully exported message: %s\n", *message.GetSubject())
		}
		if logger != nil {
			logger.WriteRow([]string{ActionSearchAndExport, StatusSuccess, mailbox, "Exported successfully", *message.GetId()})
		}
	}

	return nil
}
