package validation

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

// CronParserOptions defines the cron format options
// Matches the operator's cron parser configuration
const CronParserOptions = cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor

// ValidateCronSchedule validates a cron schedule string
func ValidateCronSchedule(schedule string) error {
	if schedule == "" {
		return fmt.Errorf("schedule cannot be empty")
	}

	parser := cron.NewParser(CronParserOptions)
	if _, err := parser.Parse(schedule); err != nil {
		return fmt.Errorf("invalid cron schedule: %w", err)
	}

	return nil
}

// CronExamples returns common cron schedule examples
func CronExamples() map[string]string {
	return map[string]string{
		"0 18 * * *":   "Every day at 6 PM",
		"0 8 * * *":    "Every day at 8 AM",
		"0 18 * * 1-5": "Weekdays at 6 PM",
		"0 9 * * 1":    "Every Monday at 9 AM",
		"*/15 * * * *": "Every 15 minutes",
		"0 0 * * 0":    "Every Sunday at midnight",
	}
}
