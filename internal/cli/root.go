package cli

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ivyascorp-net/nagging-nancy/internal/app"
	"github.com/ivyascorp-net/nagging-nancy/internal/models"
	"github.com/ivyascorp-net/nagging-nancy/internal/tui"
)

var (
	appInstance *app.App
	rootCmd     = &cobra.Command{
		Use:   "nancy",
		Short: "Nagging Nancy - Your friendly terminal reminders app",
		Long: `Nancy is a fast, lightweight terminal application that helps you
manage reminders and tasks without leaving your command line.

Built with Go and Bubble Tea for a smooth, responsive experience.`,
		Version: app.GetVersion(),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default action - launch TUI
			return runTUI()
		},
	}
)

func init() {
	// Initialize the app instance
	var err error
	appInstance, err = app.New()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Add subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(testCmd)
	// rootCmd.AddCommand(tuiCmd)
	// rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)

	// Global flags
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// runTUI launches the terminal user interface
func runTUI() error {
	// Create TUI model
	model := tui.NewModel(appInstance.GetStore(), appInstance.GetConfig())

	// Create Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Start the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}

	return nil
}

// Version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(app.GetVersion())
	},
}

// getApp returns the global app instance
func getApp() *app.App {
	return appInstance
}

// checkError is a helper function for error handling in commands
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// formatReminder formats a reminder for display in CLI
func formatReminder(reminder *models.Reminder, index int) string {
	status := "●"
	if reminder.Completed {
		status = "✓"
	}

	priorityIcon := reminder.Priority.Icon()
	timeStr := reminder.FormattedDueTime()

	return fmt.Sprintf("%d. %s %s %s - %s",
		index+1, status, priorityIcon, reminder.Title, timeStr)
}
