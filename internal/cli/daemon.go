package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/ivyascorp-net/nagging-nancy/internal/app"
	"github.com/ivyascorp-net/nagging-nancy/internal/models"
	"github.com/ivyascorp-net/nagging-nancy/internal/utils"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Daemon management commands",
	Long:  `Start, stop, and manage the Nancy daemon for background reminder notifications.`,
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Nancy daemon",
	Long:  `Start the Nancy daemon to monitor reminders and send notifications.`,
	RunE:  startDaemon,
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Nancy daemon",
	Long:  `Stop the Nancy daemon.`,
	RunE:  stopDaemon,
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check daemon status",
	Long:  `Check if the Nancy daemon is running.`,
	RunE:  daemonStatus,
}

var daemonRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the Nancy daemon",
	Long:  `Restart the Nancy daemon.`,
	RunE:  restartDaemon,
}

func init() {
	daemonCmd.AddCommand(daemonStartCmd)
	daemonCmd.AddCommand(daemonStopCmd)
	daemonCmd.AddCommand(daemonStatusCmd)
	daemonCmd.AddCommand(daemonRestartCmd)

	// Flags for daemon start
	daemonStartCmd.Flags().Duration("interval", 5*time.Minute, "Check interval for reminders")
	daemonStartCmd.Flags().Bool("foreground", false, "Run in foreground (don't daemonize)")
}

// Daemon represents the background daemon process
type Daemon struct {
	app           *app.App
	checkInterval time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
	notifier      *utils.Notifier
	lastNotified  map[string]time.Time // Track last notification time per reminder ID
}

// NewDaemon creates a new daemon instance
func NewDaemon(app *app.App, checkInterval time.Duration) (*Daemon, error) {
	notifier, err := utils.NewNotifier()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize notifier: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Daemon{
		app:           app,
		checkInterval: checkInterval,
		ctx:           ctx,
		cancel:        cancel,
		notifier:      notifier,
		lastNotified:  make(map[string]time.Time),
	}, nil
}

// Run starts the daemon monitoring loop
func (d *Daemon) Run() error {
	log.Printf("Nancy daemon started with check interval: %v", d.checkInterval)

	ticker := time.NewTicker(d.checkInterval)
	defer ticker.Stop()

	// Immediate check on startup
	d.checkReminders()

	for {
		select {
		case <-d.ctx.Done():
			log.Println("Nancy daemon stopped")
			return nil
		case <-ticker.C:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic in checkReminders: %v", r)
					}
				}()
				d.checkReminders()
			}()
		}
	}
}

// Stop gracefully stops the daemon
func (d *Daemon) Stop() {
	if d.cancel != nil {
		d.cancel()
	}
}

// checkReminders checks for due reminders and sends notifications
func (d *Daemon) checkReminders() {
	log.Printf("Checking reminders at %v", time.Now())

	// Reload reminders from storage to get any updates made by other processes
	store := d.app.GetStore()
	if err := store.Load(); err != nil {
		log.Printf("Failed to reload reminders from storage: %v", err)
		return
	}

	filter := &models.FilterOptions{
		ShowCompleted: false,
	}

	reminders := d.app.GetReminders(filter)
	now := time.Now()

	log.Printf("Found %d active reminders to check (reloaded from storage)", len(reminders))

	// Clean up notification tracking for reminders that no longer exist
	currentReminderIDs := make(map[string]bool)
	for _, reminder := range reminders {
		currentReminderIDs[reminder.ID] = true
	}

	// Remove tracking for deleted reminders
	for reminderID := range d.lastNotified {
		if !currentReminderIDs[reminderID] {
			delete(d.lastNotified, reminderID)
			log.Printf("Cleaned up notification tracking for deleted reminder: %s", reminderID)
		}
	}

	for _, reminder := range reminders {
		// Skip if already completed
		if reminder.Completed {
			continue
		}

		// Check if we should notify for this reminder
		shouldNotify := false
		notificationType := ""

		if reminder.IsOverdue() {
			// Check if we haven't notified about overdue in the last hour
			lastNotified, exists := d.lastNotified[reminder.ID]
			if !exists || now.Sub(lastNotified) > time.Hour {
				shouldNotify = true
				notificationType = "overdue"
			}
		} else if reminder.IsDueSoon() {
			// Check if we haven't notified about due soon in the last 15 minutes
			lastNotified, exists := d.lastNotified[reminder.ID]
			if !exists || now.Sub(lastNotified) > 15*time.Minute {
				shouldNotify = true
				notificationType = "due_soon"
			}
		} else if reminder.IsDueToday() {
			// Check if we haven't notified about due today
			lastNotified, exists := d.lastNotified[reminder.ID]
			todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			if !exists || lastNotified.Before(todayStart) {
				shouldNotify = true
				notificationType = "due_today"
			}
		}

		if shouldNotify {
			if err := d.sendNotification(reminder, notificationType); err != nil {
				log.Printf("Failed to send notification for reminder %s: %v", reminder.ID, err)
			} else {
				d.lastNotified[reminder.ID] = now
				log.Printf("Sent %s notification for: %s", notificationType, reminder.Title)
			}
		}
	}
}

// sendNotification sends a notification for the given reminder
func (d *Daemon) sendNotification(reminder *models.Reminder, notificationType string) error {
	var title, message string

	switch notificationType {
	case "overdue":
		title = "Overdue Reminder"
		message = fmt.Sprintf("⚠️ %s\nDue: %s", reminder.Title, reminder.FormattedDueTime())
	case "due_soon":
		title = "Reminder Due Soon"
		message = fmt.Sprintf("⏰ %s\nDue: %s", reminder.Title, reminder.FormattedDueTime())
	case "due_today":
		title = "Reminder Due Today"
		message = fmt.Sprintf("📅 %s\nDue: %s", reminder.Title, reminder.FormattedDueTime())
	default:
		title = "Nancy Reminder"
		message = reminder.Title
	}

	return d.notifier.Send(title, message, reminder.Priority)
}

// getPIDFilePath returns the path to the daemon PID file
func getPIDFilePath() (string, error) {
	app, err := app.New()
	if err != nil {
		return "", err
	}

	configDir := app.GetConfig().GetConfigDir()
	return filepath.Join(configDir, "daemon.pid"), nil
}

// writePIDFile writes the current process ID to the PID file
func writePIDFile() error {
	pidFile, err := getPIDFilePath()
	if err != nil {
		return err
	}

	pid := os.Getpid()
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

// removePIDFile removes the PID file
func removePIDFile() error {
	pidFile, err := getPIDFilePath()
	if err != nil {
		return err
	}

	return os.Remove(pidFile)
}

// isDaemonRunning checks if the daemon is currently running
func isDaemonRunning() (bool, int, error) {
	pidFile, err := getPIDFilePath()
	if err != nil {
		return false, 0, err
	}

	data, err := os.ReadFile(pidFile)
	if os.IsNotExist(err) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false, 0, err
	}

	// Check if process is running
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, pid, nil
	}

	// On Unix systems, sending signal 0 checks if process exists
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		// Process doesn't exist, clean up stale PID file
		removePIDFile()
		return false, pid, nil
	}

	return true, pid, nil
}

// startDaemon starts the Nancy daemon
func startDaemon(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	foreground, _ := cmd.Flags().GetBool("foreground")

	// Only check if daemon is already running when not in foreground mode
	// (foreground mode is used by the daemonized child process)
	if !foreground {
		if running, pid, err := isDaemonRunning(); err != nil {
			return fmt.Errorf("failed to check daemon status: %w", err)
		} else if running {
			return fmt.Errorf("daemon is already running with PID %d", pid)
		}
	}

	app := getApp()
	daemon, err := NewDaemon(app, interval)
	if err != nil {
		return fmt.Errorf("failed to create daemon: %w", err)
	}

	if !foreground {
		// Daemonize: fork and run in background
		return daemonizeProcess(interval)
	}

	// Foreground mode: run in current process (write PID file for tracking)
	if err := writePIDFile(); err != nil {
		log.Printf("Warning: failed to write PID file: %v", err)
	}
	
	// Set up cleanup on exit
	defer func() {
		if err := removePIDFile(); err != nil {
			log.Printf("Warning: failed to remove PID file: %v", err)
		}
	}()

	return runDaemonForeground(daemon, interval)
}

// daemonizeProcess forks the process and runs the daemon in background
func daemonizeProcess(interval time.Duration) error {
	// Fork the process using exec to create a true daemon
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Prepare arguments for the background process
	args := []string{
		"daemon", "start",
		"--foreground", // The child process will run in foreground mode
		"--interval", interval.String(),
	}

	// Start the process in background
	cmd := exec.Command(executable, args...)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Setctty: false, // Create new session (detach from terminal)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon process: %w", err)
	}

	// Write PID file
	pidFile, err := getPIDFilePath()
	if err != nil {
		return fmt.Errorf("failed to get PID file path: %w", err)
	}

	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	fmt.Printf("Nancy daemon started with PID %d\n", cmd.Process.Pid)
	return nil
}

// runDaemonForeground runs the daemon in the current process
func runDaemonForeground(daemon *Daemon, interval time.Duration) error {
	fmt.Println("Nancy daemon started in foreground mode")
	fmt.Printf("Check interval: %v\n", interval)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start daemon in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- daemon.Run()
	}()

	// Wait for signal or error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		daemon.Stop()
		return nil
	case err := <-errChan:
		if err != nil {
			return err
		}
		return nil
	}
}

// stopDaemon stops the Nancy daemon
func stopDaemon(cmd *cobra.Command, args []string) error {
	running, pid, err := isDaemonRunning()
	if err != nil {
		return fmt.Errorf("failed to check daemon status: %w", err)
	}

	if !running {
		fmt.Println("Daemon is not running")
		return nil
	}

	// Send TERM signal to the process
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send TERM signal to process %d: %w", pid, err)
	}

	// Wait a bit and check if process stopped
	time.Sleep(time.Second)
	if running, _, _ := isDaemonRunning(); !running {
		fmt.Println("Daemon stopped")
		return nil
	}

	// If still running, force kill
	if err := process.Signal(syscall.SIGKILL); err != nil {
		return fmt.Errorf("failed to force kill process %d: %w", pid, err)
	}

	fmt.Println("Daemon force stopped")
	return nil
}

// daemonStatus checks the daemon status
func daemonStatus(cmd *cobra.Command, args []string) error {
	running, pid, err := isDaemonRunning()
	if err != nil {
		return fmt.Errorf("failed to check daemon status: %w", err)
	}

	if running {
		fmt.Printf("Daemon is running with PID %d\n", pid)
	} else {
		fmt.Println("Daemon is not running")
	}

	return nil
}

// restartDaemon restarts the Nancy daemon
func restartDaemon(cmd *cobra.Command, args []string) error {
	// Stop if running
	if running, _, _ := isDaemonRunning(); running {
		if err := stopDaemon(cmd, args); err != nil {
			return err
		}
		// Wait a moment for cleanup
		time.Sleep(time.Second)
	}

	// Start
	return startDaemon(cmd, args)
}
