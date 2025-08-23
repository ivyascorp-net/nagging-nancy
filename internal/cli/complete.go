package cli

import (
	"fmt"
	"strings"

	"github.com/ivyascorp-net/nagging-nancy/internal/models"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete <reminder-id>",
	Short: "Mark a reminder as completed",
	Long: `Mark one or more reminders as completed by their ID.

You can find reminder IDs by running 'nancy list'.
You can specify multiple IDs separated by spaces.`,
	Aliases: []string{"done", "finish"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store := getApp().GetStore()
		var errors []string
		var completed []string

		for _, idArg := range args {
			// Find reminder by partial ID match
			reminder, err := findReminderByID(idArg)
			if err != nil {
				errors = append(errors, fmt.Sprintf("ID %s: %v", idArg, err))
				continue
			}

			// Check if already completed
			if reminder.Completed {
				errors = append(errors, fmt.Sprintf("ID %s: already completed", idArg))
				continue
			}

			// Mark as completed
			if err := store.CompleteReminder(reminder.ID); err != nil {
				errors = append(errors, fmt.Sprintf("ID %s: failed to complete - %v", idArg, err))
				continue
			}

			completed = append(completed, fmt.Sprintf("✅ %s", reminder.Title))
		}

		// Display results
		if len(completed) > 0 {
			fmt.Println("Completed reminders:")
			for _, item := range completed {
				fmt.Println("  " + item)
			}
		}

		if len(errors) > 0 {
			fmt.Println("\nErrors:")
			for _, err := range errors {
				fmt.Println("  ❌ " + err)
			}
			return fmt.Errorf("some reminders could not be completed")
		}

		if len(completed) == 1 {
			fmt.Println("\n🎉 Great job getting that done!")
		} else if len(completed) > 1 {
			fmt.Printf("\n🎉 Wow! You completed %d reminders. You're on fire!\n", len(completed))
		}

		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <reminder-id>",
	Short: "Delete a reminder",
	Long: `Delete one or more reminders permanently by their ID.

You can find reminder IDs by running 'nancy list'.
You can specify multiple IDs separated by spaces.

Warning: This action cannot be undone!`,
	Aliases: []string{"del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store := getApp().GetStore()
		var errors []string
		var deleted []string

		// Confirmation flag
		force, _ := cmd.Flags().GetBool("force")

		if !force && len(args) > 1 {
			fmt.Printf("⚠️  You are about to delete %d reminders. Use --force to confirm.\n", len(args))
			return nil
		}

		for _, idArg := range args {
			// Find reminder by partial ID match
			reminder, err := findReminderByID(idArg)
			if err != nil {
				errors = append(errors, fmt.Sprintf("ID %s: %v", idArg, err))
				continue
			}

			// Confirm deletion for single items (unless forced)
			if !force && len(args) == 1 {
				fmt.Printf("⚠️  Delete reminder: %s? [y/N]: ", reminder.Title)
				var response string
				fmt.Scanln(&response)

				if strings.ToLower(strings.TrimSpace(response)) != "y" &&
					strings.ToLower(strings.TrimSpace(response)) != "yes" {
					fmt.Println("❌ Deletion cancelled.")
					return nil
				}
			}

			// Delete the reminder
			if err := store.Delete(reminder.ID); err != nil {
				errors = append(errors, fmt.Sprintf("ID %s: failed to delete - %v", idArg, err))
				continue
			}

			deleted = append(deleted, fmt.Sprintf("🗑️  %s", reminder.Title))
		}

		// Display results
		if len(deleted) > 0 {
			fmt.Println("Deleted reminders:")
			for _, item := range deleted {
				fmt.Println("  " + item)
			}
		}

		if len(errors) > 0 {
			fmt.Println("\nErrors:")
			for _, err := range errors {
				fmt.Println("  ❌ " + err)
			}
			return fmt.Errorf("some reminders could not be deleted")
		}

		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompts")

	completeCmd.Example = `  # Complete a reminder by ID
  nancy complete a1b2c3d4

  # Complete multiple reminders
  nancy complete a1b2c3d4 e5f6g7h8

  # Using short ID (first 8 characters)
  nancy done a1b2c3d4`

	deleteCmd.Example = `  # Delete a reminder (with confirmation)
  nancy delete a1b2c3d4

  # Delete multiple reminders (requires --force)
  nancy delete a1b2c3d4 e5f6g7h8 --force

  # Force delete without confirmation
  nancy rm a1b2c3d4 --force`
}

// findReminderByID finds a reminder by full or partial ID
func findReminderByID(idArg string) (*models.Reminder, error) {
	store := getApp().GetStore()

	// First try exact match
	if reminder, err := store.Get(idArg); err == nil {
		return reminder, nil
	}

	// If not found and it's a short ID, try to find by prefix
	if len(idArg) >= 4 { // Minimum 4 characters for partial match
		allReminders := store.GetAll(&models.FilterOptions{ShowCompleted: true})

		var matches []*models.Reminder
		for _, reminder := range allReminders {
			if strings.HasPrefix(reminder.ID, idArg) {
				matches = append(matches, reminder)
			}
		}

		if len(matches) == 1 {
			return matches[0], nil
		} else if len(matches) > 1 {
			return nil, fmt.Errorf("ambiguous ID '%s' matches multiple reminders", idArg)
		}
	}

	return nil, fmt.Errorf("reminder not found")
}
