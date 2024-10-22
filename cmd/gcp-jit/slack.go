package main
import (
	"fmt"
	"context"
	"log"
	"encoding/json"
	"encoding/base64"
	"strings"

	"golang.org/x/oauth2/google"
	"github.com/slack-go/slack"
)

func sendSlackMessage(cfg *slackConfig, link string) (error) {
	api := slack.New(cfg.APIToken)
	email := getEmailAddress()

	message := fmt.Sprintf("A new PAM JIT request has been submitted by %s. Please review and approve: %s", email, link)

	// send the message to Slack
	_, _, err := api.PostMessage(cfg.ChannelID, slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("Error sending message: %s\n", err)
		return fmt.Errorf("error sending message to Slack: %w", err)
	}

	fmt.Println("Sent request to Slack")

	return nil

}

func getEmailAddress() string {
	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatalf("Failed to get default credentials: %v", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Extract email from access token
	parts := strings.Split(token.AccessToken, ".")
	if len(parts) != 3 {
		log.Fatal("Invalid access token format")
	}

	// Decode the payload (second part of the token)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Fatalf("Failed to decode token payload: %v", err)
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		log.Fatalf("Failed to parse token payload: %v", err)
	}

	fmt.Printf("Authenticated email: %s\n", claims.Email)
	return claims.Email
}
