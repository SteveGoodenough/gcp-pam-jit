package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func sendSlackMessage(link string) (error) {
	// Replace with your Slack API token
	apiToken := os.Getenv("GCP_JIT_SLACK_API_TOKEN")
	channelID := os.Getenv("GCP_JIT_SLACK_CHANNEL_ID")

	if apiToken == "" || channelID == "" {
		return fmt.Errorf("GCP_JIT_SLACK_API_TOKEN or GCP_JIT_SLACK_CHANNEL_ID environment variable not set")
	}

	api := slack.New(apiToken)

	message := fmt.Sprintf("A new PAM JIT request has been submitted. Please review and approve: %s", link)

	// send the message to Slack
	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("Error sending message: %s\n", err)
		return fmt.Errorf("error sending message to Slack: %w", err)
	}

	fmt.Println("Sent request to Slack")

	return nil

}
