package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/ivyascorp-net/nagging-nancy/internal/models"
	"github.com/ivyascorp-net/nagging-nancy/internal/utils"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <reminder text>",
	Short: "Add a new reminder",
	Long: `Add a new reminder with optional time, date, and priority.

Examples:
  nancy add "Call mom"
  nancy add "Meeting" --time "2pm" --priority high
  nancy add "Buy groceries tomorrow at 5pm"
  nancy add "Submit report urgent" --date "2024-03-20"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		timeFlag, _ := cmd.Flags().GetString("time")
		dateFlag, _ := cmd.Flags().GetString("date")
		priorityFlag, _ := cmd.Flags().GetString("priority")
		tagsFlag, _ := cmd.Flags().GetStringSlice("tags")

		// Join all arguments as the reminder text
		reminderText := strings.Join(args, " ")

		// Parse the reminder text for natural language time/priority
		config := getApp().GetConfig()
		defaultPriority := models.ParsePriority(config.Default.Priority)

		parsed, err := utils.ParseReminder(reminderText, defaultPriority)
		if err != nil {
			return fmt.Errorf("failed to parse reminder: %w", err)
		}

		// Override with explicit flags if provided
		dueTime := parsed.DueTime
		priority := parsed.Priority
		title := parsed.Title
		tags := parsed.Tags

		// Handle explicit time flag
		if timeFlag != "" {
			parsedTime, err := utils.ParseTimeString(timeFlag)
			if err != nil {
				return fmt.Errorf("invalid time format '%s': %w", timeFlag, err)
			}

			// If only time provided, use today's date
			now := time.Now()
			dueTime = time.Date(now.Year(), now.Month(), now.Day(),
				parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())

			// If time has passed today, schedule for tomorrow
			if dueTime.Before(now) {
				dueTime = dueTime.AddDate(0, 0, 1)
			}
		}

		// Handle explicit date flag
		if dateFlag != "" {
			var targetDate time.Time
			var err error

			// Try parsing different date formats
			dateFormats := []string{
				"2006-01-02",  // 2024-03-20
				"01/02/2006",  // 03/20/2024
				"01-02-2006",  // 03-20-2024
				"Jan 2, 2006", // Mar 20, 2024
				"Jan 2 2006",  // Mar 20 2024
				"2 Jan 2006",  // 20 Mar 2024
			}

			// Handle relative dates
			switch strings.ToLower(dateFlag) {
			case "today":
				targetDate = time.Now()
			case "tomorrow":
				targetDate = time.Now().AddDate(0, 0, 1)
			default:
				// Try parsing as explicit date
				for _, format := range dateFormats {
					if targetDate, err = time.Parse(format, dateFlag); err == nil {
						break
					}
				}
				if err != nil {
					return fmt.Errorf("invalid date format '%s'", dateFlag)
				}
			}

			// Combine date with existing time
			dueTime = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
				dueTime.Hour(), dueTime.Minute(), 0, 0, dueTime.Location())
		}

		// Handle explicit priority flag
		if priorityFlag != "" {
			priority = utils.ParsePriorityString(priorityFlag)
		}

		// Handle explicit tags flag
		if len(tagsFlag) > 0 {
			// Merge with parsed tags
			tagSet := make(map[string]bool)
			for _, tag := range tags {
				tagSet[tag] = true
			}
			for _, tag := range tagsFlag {
				tagSet[strings.TrimSpace(tag)] = true
			}

			// Convert back to slice
			tags = make([]string, 0, len(tagSet))
			for tag := range tagSet {
				if tag != "" {
					tags = append(tags, tag)
				}
			}
		}

		// Validate input
		if err := utils.ValidateReminderInput(title, dueTime); err != nil {
			return err
		}

		// Create reminder
		reminder := models.NewReminder(title, dueTime, priority)

		// Add tags
		for _, tag := range tags {
			reminder.AddTag(tag)
		}

		// Save to store
		if err := getApp().GetStore().Add(reminder); err != nil {
			return fmt.Errorf("failed to add reminder: %w", err)
		}

		// Output confirmation
		fmt.Printf("✅ Added reminder: %s\n", reminder.Title)
		fmt.Printf("   Due: %s\n", reminder.FormattedDueTime())
		fmt.Printf("   Priority: %s %s\n", priority.Icon(), priority.String())

		if len(tags) > 0 {
			fmt.Printf("   Tags: %s\n", strings.Join(tags, ", "))
		}

		// Show ID for reference
		fmt.Printf("   ID: %s\n", reminder.ID[:8])

		return nil
	},
}

func init() {
	addCmd.Flags().StringP("time", "t", "", "Due time (e.g., 2pm, 14:30, '3:30 PM')")
	addCmd.Flags().StringP("date", "d", "", "Due date (e.g., tomorrow, 2024-03-20, 'Mar 20')")
	addCmd.Flags().StringP("priority", "p", "", "Priority level (low, medium, high)")
	addCmd.Flags().StringSliceP("tags", "", []string{}, "Tags for the reminder (e.g., work,urgent)")

	// Add examples to help
	addCmd.Example = `  # Simple reminder
  nancy add "Call mom"

  # With specific time
  nancy add "Team meeting" --time "2pm"

  # With date and priority
  nancy add "Submit report" --date "tomorrow" --priority high

  # Natural language parsing
  nancy add "Doctor appointment tomorrow at 3pm urgent"

  # With tags
  nancy add "Review code" --tags "work,coding" --priority medium`
}
