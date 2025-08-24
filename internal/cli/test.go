package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ivyascorp-net/nagging-nancy/internal/utils"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test system functionality",
	Long:  `Test various system components like notifications.`,
}

var testNotificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Test notification system",
	Long:  `Send a test notification to verify the notification system is working.`,
	RunE:  testNotification,
}

func init() {
	testCmd.AddCommand(testNotificationCmd)
}

// testNotification sends a test notification
func testNotification(cmd *cobra.Command, args []string) error {
	notifier, err := utils.NewNotifier()
	if err != nil {
		return fmt.Errorf("failed to create notifier: %w", err)
	}

	fmt.Printf("Using notification method: %s\n", utils.GetMethodName(notifier.GetMethod()))
	fmt.Println("Sending test notification...")

	if err := notifier.TestNotification(); err != nil {
		return fmt.Errorf("failed to send test notification: %w", err)
	}

	fmt.Println("Test notification sent successfully!")
	
	// Show available methods
	methods := utils.GetAvailableMethods()
	fmt.Println("\nAvailable notification methods:")
	for _, method := range methods {
		fmt.Printf("  - %s\n", utils.GetMethodName(method))
	}

	return nil
}