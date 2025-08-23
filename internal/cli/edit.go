package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/ivyascorp-net/nagging-nancy/internal/utils"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <reminder-id>",
	Short: "Edit an existing reminder",
	Long: `Edit the title, due time, or priority of an existing reminder.

You can find reminder IDs by running 'nancy list'.

Examples:
  nancy edit a1b2c3d4 --title "New title"
  nancy edit a1b2c3d4 --time "3pm"
  nancy edit a1b2c3d4 --priority high
  nancy edit a1b2c3d4 --title "Call mom" --time "tomorrow 2pm" --priority high`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idArg := args[0]

		// Find the reminder
		reminder, err := findReminderByID(idArg)
		if err != nil {
			return fmt.Errorf("reminder not found: %w", err)
		}

		// Get flags
		title, _ := cmd.Flags().GetString("title")
		timeFlag, _ := cmd.Flags().GetString("time")
		dateFlag, _ := cmd.Flags().GetString("date")
		priorityFlag, _ := cmd.Flags().GetString("priority")
		addTags, _ := cmd.Flags().GetStringSlice("add-tags")
		removeTags, _ := cmd.Flags().GetStringSlice("remove-tags")

		// Track what changed
		var changes []string

		// Update title
		if title != "" {
			reminder.Title = title
			changes = append(changes, fmt.Sprintf("title → '%s'", title))
		}

		// Update time
		newDueTime := reminder.DueTime
		if timeFlag != "" {
			parsedTime, err := utils.ParseTimeString(timeFlag)
			if err != nil {
				return fmt.Errorf("invalid time format '%s': %w", timeFlag, err)
			}

			// If only time provided, use current date
			newDueTime = time.Date(newDueTime.Year(), newDueTime.Month(), newDueTime.Day(),
				parsedTime.Hour(), parsedTime.Minute(), 0, 0, newDueTime.Location())
			changes = append(changes, fmt.Sprintf("time → %s", parsedTime.Format("3:04 PM")))
		}

		// Update date
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
			newDueTime = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
				newDueTime.Hour(), newDueTime.Minute(), 0, 0, newDueTime.Location())
			changes = append(changes, fmt.Sprintf("date → %s", targetDate.Format("Jan 2, 2006")))
		}

		// Update due time if it changed
		if !newDueTime.Equal(reminder.DueTime) {
			reminder.DueTime = newDueTime
		}

		// Update priority
		if priorityFlag != "" {
			oldPriority := reminder.Priority
			newPriority := utils.ParsePriorityString(priorityFlag)
			if newPriority != oldPriority {
				reminder.Priority = newPriority
				changes = append(changes, fmt.Sprintf("priority → %s %s", 
					newPriority.Icon(), newPriority.String()))
			}
		}

		// Add tags
		for _, tag := range addTags {
			tag = strings.TrimSpace(tag)
			if tag != "" && !reminder.HasTag(tag) {
				reminder.AddTag(tag)
				changes = append(changes, fmt.Sprintf("added tag '%s'", tag))
			}
		}

		// Remove tags
		for _, tag := range removeTags {
			tag = strings.TrimSpace(tag)
			if tag != "" && reminder.HasTag(tag) {
				reminder.RemoveTag(tag)
				changes = append(changes, fmt.Sprintf("removed tag '%s'", tag))
			}
		}

		// Validate changes
		if len(changes) == 0 {
			fmt.Println("No changes specified. Use --title, --time, --date, --priority, --add-tags, or --remove-tags")
			return nil
		}

		// Validate the updated reminder
		if err := utils.ValidateReminderInput(reminder.Title, reminder.DueTime); err != nil {
			return err
		}

		// Save changes
		if err := getApp().GetStore().Update(reminder); err != nil {
			return fmt.Errorf("failed to update reminder: %w", err)
		}

		// Show confirmation
		fmt.Printf("✅ Updated reminder: %s\n", reminder.Title)
		fmt.Printf("   Due: %s\n", reminder.FormattedDueTime())
		fmt.Printf("   Priority: %s %s\n", reminder.Priority.Icon(), reminder.Priority.String())

		if len(reminder.Tags) > 0 {
			fmt.Printf("   Tags: %s\n", strings.Join(reminder.Tags, ", "))
		}

		fmt.Printf("   ID: %s\n\n", reminder.ID[:8])

		fmt.Println("Changes made:")
		for _, change := range changes {
			fmt.Printf("  • %s\n", change)
		}

		return nil
	},
}

func init() {
	editCmd.Flags().StringP("title", "", "", "New title for the reminder")
	editCmd.Flags().StringP("time", "t", "", "New due time (e.g., 2pm, 14:30, '3:30 PM')")
	editCmd.Flags().StringP("date", "d", "", "New due date (e.g., tomorrow, 2024-03-20, 'Mar 20')")
	editCmd.Flags().StringP("priority", "p", "", "New priority level (low, medium, high)")
	editCmd.Flags().StringSliceP("add-tags", "", []string{}, "Tags to add (e.g., work,urgent)")
	editCmd.Flags().StringSliceP("remove-tags", "", []string{}, "Tags to remove")

	editCmd.Example = `  # Edit title
  nancy edit a1b2c3d4 --title "New reminder title"

  # Edit time and priority
  nancy edit a1b2c3d4 --time "3pm" --priority high

  # Edit date
  nancy edit a1b2c3d4 --date "tomorrow"

  # Add and remove tags
  nancy edit a1b2c3d4 --add-tags "work,urgent" --remove-tags "personal"

  # Multiple changes at once
  nancy edit a1b2c3d4 --title "Call mom" --time "2pm" --priority high`
}
