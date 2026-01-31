package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Status constants
const (
	StatusSuccess = "Success"
	StatusError   = "Error"
)

// executeAction dispatches to the appropriate action handler based on config.Action.
// Supported actions are: getevents, sendmail, sendinvite, and getinbox.
//
// For sendmail action, if no recipients are specified, the email is sent to the
// mailbox owner (self). All actions log their operations to the provided CSV logger.
//
// Returns an error if the action fails or if the action name is unknown.
func executeAction(ctx context.Context, client *msgraphsdk.GraphServiceClient, config *Config, logger *CSVLogger) error {
	switch config.Action {
	case ActionGetEvents:
		if err := listEvents(ctx, client, config.Mailbox, config.Count, config, logger); err != nil {
			return fmt.Errorf("failed to list events: %w", err)
		}
	case ActionSendMail:
		// If no recipients specified at all, default 'to' to the sender mailbox
		if len(config.To) == 0 && len(config.Cc) == 0 && len(config.Bcc) == 0 {
			config.To = []string{config.Mailbox}
		}

		sendEmail(ctx, client, config.Mailbox, config.To, config.Cc, config.Bcc, config.Subject, config.Body, config.BodyHTML, config.AttachmentFiles, config, logger)
	case ActionSendInvite:
		// Use Subject for calendar invite
		// For backward compatibility, if InviteSubject is set, use it instead
		inviteSubject := config.Subject
		if config.InviteSubject != "" {
			inviteSubject = config.InviteSubject
		}
		// If using default email subject, change to default calendar invite subject
		if inviteSubject == "Automated Tool Notification" {
			inviteSubject = "It's testing event"
		}
		createInvite(ctx, client, config.Mailbox, inviteSubject, config.StartTime, config.EndTime, config, logger)
	case ActionGetInbox:
		if err := listInbox(ctx, client, config.Mailbox, config.Count, config, logger); err != nil {
			return fmt.Errorf("failed to list inbox: %w", err)
		}
	case ActionGetSchedule:
		if err := checkAvailability(ctx, client, config.Mailbox, config.To[0], config, logger); err != nil {
			return fmt.Errorf("failed to check availability: %w", err)
		}
	case ActionExportInbox:
		if err := exportInbox(ctx, client, config.Mailbox, config.Count, config, logger); err != nil {
			return fmt.Errorf("failed to export inbox: %w", err)
		}
	case ActionSearchAndExport:
		if err := searchAndExport(ctx, client, config.Mailbox, config.MessageID, config, logger); err != nil {
			return fmt.Errorf("failed to search and export: %w", err)
		}
	default:
		return fmt.Errorf("unknown action: %s", config.Action)
	}

	return nil
}


// exportMessageToJSON serializes a message to JSON and saves it to a file
func printJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON output: %v\n", err)
	}
}

// formatEventsOutput converts a list of Eventable items to a JSON-friendly slice of maps
func formatEventsOutput(events []models.Eventable) []map[string]interface{} {
	var output []map[string]interface{}
	for _, event := range events {
		eventMap := make(map[string]interface{})
		if event.GetId() != nil {
			eventMap["id"] = *event.GetId()
		}
		if event.GetSubject() != nil {
			eventMap["subject"] = *event.GetSubject()
		}
		if event.GetStart() != nil && event.GetStart().GetDateTime() != nil {
			eventMap["start"] = *event.GetStart().GetDateTime()
		}
		if event.GetEnd() != nil && event.GetEnd().GetDateTime() != nil {
			eventMap["end"] = *event.GetEnd().GetDateTime()
		}
		if event.GetOrganizer() != nil && event.GetOrganizer().GetEmailAddress() != nil {
			eventMap["organizer"] = extractEmailAddress(event.GetOrganizer().GetEmailAddress())
		}
		output = append(output, eventMap)
	}
	return output
}

// formatMessagesOutput converts a list of Messageable items to a JSON-friendly slice of maps
func formatMessagesOutput(messages []models.Messageable) []map[string]interface{} {
	var output []map[string]interface{}
	for _, message := range messages {
		msgMap := make(map[string]interface{})
		if message.GetId() != nil {
			msgMap["id"] = *message.GetId()
		}
		if message.GetSubject() != nil {
			msgMap["subject"] = *message.GetSubject()
		}
		if message.GetReceivedDateTime() != nil {
			msgMap["receivedDateTime"] = message.GetReceivedDateTime().Format(time.RFC3339)
		}
		if message.GetFrom() != nil && message.GetFrom().GetEmailAddress() != nil {
			msgMap["from"] = extractEmailAddress(message.GetFrom().GetEmailAddress())
		}
		if message.GetToRecipients() != nil {
			msgMap["toRecipients"] = extractRecipients(message.GetToRecipients())
		}
		output = append(output, msgMap)
	}
	return output
}

// formatScheduleOutput converts a list of ScheduleInformationable items to a JSON-friendly structure
func formatScheduleOutput(schedules []models.ScheduleInformationable) []map[string]interface{} {
	var output []map[string]interface{}
	for _, schedule := range schedules {
		schMap := make(map[string]interface{})
		if schedule.GetScheduleId() != nil {
			schMap["scheduleId"] = *schedule.GetScheduleId()
		}
		if schedule.GetAvailabilityView() != nil {
			schMap["availabilityView"] = *schedule.GetAvailabilityView()
			schMap["availabilityStatus"] = interpretAvailability(*schedule.GetAvailabilityView())
		}
		// Include working hours if available
		if schedule.GetWorkingHours() != nil {
			wh := schedule.GetWorkingHours()
			whMap := make(map[string]interface{})
			if wh.GetStartTime() != nil {
				whMap["startTime"] = *wh.GetStartTime()
			}
			if wh.GetEndTime() != nil {
				whMap["endTime"] = *wh.GetEndTime()
			}
			schMap["workingHours"] = whMap
		}
		output = append(output, schMap)
	}
	return output
}
