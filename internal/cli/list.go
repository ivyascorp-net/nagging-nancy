package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/ivyascorp-net/nagging-nancy/internal/models"
	"github.com/ivyascorp-net/nagging-nancy/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List reminders",
	Long: `List reminders with optional filtering.

Examples:
  nancy list                    # All active reminders
  nancy list --today           # Today's reminders only
  nancy list --priority high   # High priority only
  nancy list --completed       # Completed reminders
  nancy list --all             # All reminders including completed`,
	Aliases: []string{"ls", "show"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		showToday, _ := cmd.Flags().GetBool("today")
		showWeek, _ := cmd.Flags().GetBool("week")
		showCompleted, _ := cmd.Flags().GetBool("completed")
		showOverdue, _ := cmd.Flags().GetBool("overdue")
		showAll, _ := cmd.Flags().GetBool("all")
		priorityFlag, _ := cmd.Flags().GetString("priority")
		tagsFlag, _ := cmd.Flags().GetStringSlice("tags")
		limit, _ := cmd.Flags().GetInt("limit")

		// Build filter options
		filter := &models.FilterOptions{
			ShowCompleted: showCompleted || showAll,
			DueToday:      showToday,
			Overdue:       showOverdue,
			Limit:         limit,
		}

		// Handle priority filter
		if priorityFlag != "" {
			priority := utils.ParsePriorityString(priorityFlag)
			filter.Priority = &priority
		}

		// Handle tags filter
		if len(tagsFlag) > 0 {
			filter.Tags = tagsFlag
		}

		// Get reminders from store
		store := getApp().GetStore()
		reminders := store.GetAll(filter)

		// Handle week filter (not in FilterOptions, so filter manually)
		if showWeek {
			weekReminders := make([]*models.Reminder, 0)
			for _, reminder := range reminders {
				if isThisWeek(reminder.DueTime) {
					weekReminders = append(weekReminders, reminder)
				}
			}
			reminders = weekReminders
		}

		// Display results
		if len(reminders) == 0 {
			if showCompleted {
				fmt.Println("📝 No completed reminders found.")
			} else if showToday {
				fmt.Println("📅 No reminders due today.")
			} else if showOverdue {
				fmt.Println("⏰ No overdue reminders.")
			} else {
				fmt.Println("🎉 All caught up! No active reminders.")
			}
			fmt.Println("\nAdd a new reminder with: nancy add \"Your reminder\"")
			return nil
		}

		// Display header
		if showCompleted {
			fmt.Println("📝 Completed Reminders")
		} else if showToday {
			fmt.Println("📅 Today's Reminders")
		} else if showOverdue {
			fmt.Println("⚠️  Overdue Reminders")
		} else if showWeek {
			fmt.Println("📆 This Week's Reminders")
		} else {
			fmt.Println("📋 Reminders")
		}

		fmt.Println(strings.Repeat("─", 50))

		// Display reminders
		for i, reminder := range reminders {
			displayReminder(reminder, i+1)
		}

		// Display summary
		fmt.Println(strings.Repeat("─", 50))

		// Get counts
		total, active, completed, overdue := store.Count()

		if showAll {
			fmt.Printf("📊 Total: %d | Active: %d | Completed: %d | Overdue: %d\n",
				total, active, completed, overdue)
		} else if showCompleted {
			fmt.Printf("📊 Showing %d completed reminders\n", len(reminders))
		} else {
			fmt.Printf("📊 Showing %d reminders | Active: %d | Overdue: %d\n",
				len(reminders), active, overdue)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().Bool("today", false, "Show only today's reminders")
	listCmd.Flags().Bool("week", false, "Show this week's reminders")
	listCmd.Flags().Bool("completed", false, "Show completed reminders")
	listCmd.Flags().Bool("overdue", false, "Show overdue reminders")
	listCmd.Flags().Bool("all", false, "Show all reminders (including completed)")
	listCmd.Flags().StringP("priority", "p", "", "Filter by priority (low, medium, high)")
	listCmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags")
	listCmd.Flags().IntP("limit", "l", 0, "Limit number of results (0 = no limit)")

	// Add examples
	listCmd.Example = `  # List active reminders
  nancy list

  # Today's reminders only
  nancy list --today

  # High priority reminders
  nancy list --priority high

  # Overdue reminders
  nancy list --overdue

  # Completed reminders
  nancy list --completed

  # All reminders with tags
  nancy list --tags work,urgent --all`
}

// displayReminder formats and displays a single reminder
func displayReminder(reminder *models.Reminder, index int) {
	// Status icon
	status := "●"
	if reminder.Completed {
		status = "✓"
	}

	// Priority icon and color would go here in a real TUI
	priorityIcon := reminder.Priority.Icon()

	// Time information
	timeStr := reminder.FormattedDueTime()

	// Status information
	statusInfo := ""
	if reminder.IsOverdue() {
		statusInfo = " ⚠️ OVERDUE"
	} else if reminder.IsDueSoon() {
		statusInfo = " ⏰ DUE SOON"
	}

	// Build the line
	fmt.Printf("%2d. %s %s %s%s\n", index, status, priorityIcon, reminder.Title, statusInfo)

	// Show due time and additional info
	fmt.Printf("    📅 %s", timeStr)

	if len(reminder.Tags) > 0 {
		fmt.Printf(" | 🏷️  %s", strings.Join(reminder.Tags, ", "))
	}

	// Show time until due for active reminders
	if !reminder.Completed {
		timeUntil := reminder.TimeUntilDue()
		if timeUntil > 0 {
			fmt.Printf(" | ⏳ %s", utils.FormatDuration(timeUntil))
		}
	}

	fmt.Printf(" | 🆔 %s\n", reminder.ID[:8])
	fmt.Println()
}

// isThisWeek checks if a time falls within the current week
func isThisWeek(t time.Time) bool {
	now := time.Now()

	// Find the start of this week (Sunday)
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	// Find the end of this week (Saturday)
	weekEnd := weekStart.AddDate(0, 0, 7)

	return t.After(weekStart) && t.Before(weekEnd)
}
