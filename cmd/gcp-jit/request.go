package main

import (
	"context"
	"fmt"
	"log"
	"github.com/felixgborrego/gpc-pam-jit/pkg/pamjit"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// slackConfig represents the configuration for Slack notifications.
type slackConfig struct {
	APIToken  string `mapstructure:"api-token"`
	ChannelID string `mapstructure:"channel-id"`
}

// config represents the application's configuration.
type config struct {
	Slack slackConfig `mapstructure:"slack"`
}

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request an entitlement",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

        // load Slack configuration
        cfg, err := loadConfig(cmd)
        if err != nil {
            log.Fatalf("Error loading config: %v", err)
        }

		entitlementID := args[0]
		projectID, _ := cmd.Flags().GetString("project")
		location, _ := cmd.Flags().GetString("location")
		justification, _ := cmd.Flags().GetString("justification")
		duration, _ := cmd.Flags().GetString("duration")

		pam, err := pamjit.NewPamJitClient(context.Background(), projectID, location)
		if err != nil {
			log.Fatalf("Unable to use GCP JIT service: %v", err)
		}
		link, err := pam.RequestGrant(cmd.Context(), entitlementID, justification, duration)
		if err != nil {
			fmt.Printf("Error requesting entitlement: %v\n", err)
		} else {
			if link != "" {
				// only attempt to send to Slack if config is set
				if cfg.Slack.APIToken != "" && cfg.Slack.ChannelID != "" {
					// send the link to Slack and if it fails (e.g. env vars not set), then display the link
					err = sendSlackMessage(&cfg.Slack, link)
					if err != nil {
						fmt.Printf("Link to request: %s\n", link)
					}
				} else {
					fmt.Printf("Link to request: %s\n", link)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringP("project", "p", "", "Project ID")
	requestCmd.Flags().StringP("location", "l", "global", "Location")
	requestCmd.Flags().StringP("justification", "j", "", "Justification")
	requestCmd.Flags().StringP("duration", "d", "", "Duration (defaults to maximum)")
	requestCmd.Flags().String("slack.api-token", "", "Slack API token (if you want to send request to Slack channel)")
	requestCmd.Flags().String("slack.channel-id", "", "Slack channel ID (if you want to send request to Slack channel)")

	requestCmd.MarkFlagRequired("project")
	requestCmd.MarkFlagRequired("justification")
}

// loadConfig loads the configuration from command line arguments, environment variables or a configuration file
func loadConfig(cmd *cobra.Command) (*config, error) {
    cfg := &config{
        Slack: slackConfig{},
    }

    // process command line flags
    if err := viper.BindPFlags(cmd.Flags()); err != nil {
        return nil, fmt.Errorf("failed to bind flags: %w", err)
    }

    // process environment variables
    viper.SetEnvPrefix("GCP_JIT") // Use prefix for environment variables
    viper.AutomaticEnv()          // Read in environment variables that match

    // search for and read a configuration file.
    viper.SetConfigName("config") // name of config file (without extension)
    viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
    // viper.AddConfigPath(".")      // optionally look for config in the working directory
    viper.AddConfigPath("$HOME/.gcp-jit") // find config in the specified directory with specified file name
    viper.ReadInConfig()
    viper.Unmarshal(&cfg)

    return cfg, nil
}