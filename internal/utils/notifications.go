package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ivyascorp-net/nagging-nancy/internal/models"
)

// NotificationMethod represents different ways to send notifications
type NotificationMethod int

const (
	// Desktop notifications via system notification daemon
	DesktopNotification NotificationMethod = iota
	// Terminal bell/beep
	TerminalBell
	// Log to file only
	LogOnly
)

// Notifier handles sending notifications across different platforms
type Notifier struct {
	method           NotificationMethod
	fallbackMethods  []NotificationMethod
	logFile          string
}

// NewNotifier creates a new notifier instance with auto-detected best method
func NewNotifier() (*Notifier, error) {
	notifier := &Notifier{
		method:          detectBestMethod(),
		fallbackMethods: []NotificationMethod{TerminalBell, LogOnly},
	}

	return notifier, nil
}

// NewNotifierWithMethod creates a notifier with a specific method
func NewNotifierWithMethod(method NotificationMethod) *Notifier {
	return &Notifier{
		method:          method,
		fallbackMethods: []NotificationMethod{TerminalBell, LogOnly},
	}
}

// detectBestMethod auto-detects the best notification method for the current system
func detectBestMethod() NotificationMethod {
	switch runtime.GOOS {
	case "linux":
		// Check for notify-send (libnotify)
		if _, err := exec.LookPath("notify-send"); err == nil {
			return DesktopNotification
		}
		// Check for dunst
		if _, err := exec.LookPath("dunstify"); err == nil {
			return DesktopNotification
		}
	case "darwin":
		// Check for osascript (built into macOS)
		if _, err := exec.LookPath("osascript"); err == nil {
			return DesktopNotification
		}
		// Check for terminal-notifier
		if _, err := exec.LookPath("terminal-notifier"); err == nil {
			return DesktopNotification
		}
	case "windows":
		// Windows has built-in notification support via PowerShell
		if _, err := exec.LookPath("powershell"); err == nil {
			return DesktopNotification
		}
	}

	// Fallback to terminal bell
	return TerminalBell
}

// Send sends a notification with the given title, message, and priority
func (n *Notifier) Send(title, message string, priority models.Priority) error {
	err := n.sendWithMethod(n.method, title, message, priority)
	if err != nil {
		// Try fallback methods
		for _, fallback := range n.fallbackMethods {
			if fallbackErr := n.sendWithMethod(fallback, title, message, priority); fallbackErr == nil {
				return nil
			}
		}
		return fmt.Errorf("all notification methods failed, last error: %w", err)
	}
	return nil
}

// sendWithMethod sends a notification using a specific method
func (n *Notifier) sendWithMethod(method NotificationMethod, title, message string, priority models.Priority) error {
	switch method {
	case DesktopNotification:
		return n.sendDesktopNotification(title, message, priority)
	case TerminalBell:
		return n.sendTerminalBell(title, message)
	case LogOnly:
		return n.logNotification(title, message)
	default:
		return fmt.Errorf("unsupported notification method: %d", method)
	}
}

// sendDesktopNotification sends a desktop notification
func (n *Notifier) sendDesktopNotification(title, message string, priority models.Priority) error {
	switch runtime.GOOS {
	case "linux":
		return n.sendLinuxDesktopNotification(title, message, priority)
	case "darwin":
		return n.sendMacOSDesktopNotification(title, message, priority)
	case "windows":
		return n.sendWindowsDesktopNotification(title, message, priority)
	default:
		return fmt.Errorf("desktop notifications not supported on %s", runtime.GOOS)
	}
}

// sendLinuxDesktopNotification sends a desktop notification on Linux
func (n *Notifier) sendLinuxDesktopNotification(title, message string, priority models.Priority) error {
	// Try notify-send first (most common)
	if _, err := exec.LookPath("notify-send"); err == nil {
		urgency := "normal"
		switch priority {
		case models.Low:
			urgency = "low"
		case models.High:
			urgency = "critical"
		}

		cmd := exec.Command("notify-send",
			"-u", urgency,
			"-a", "Nancy",
			"-i", "appointment-soon", // Standard icon
			title,
			message,
		)
		return cmd.Run()
	}

	// Try dunstify as fallback
	if _, err := exec.LookPath("dunstify"); err == nil {
		urgency := "normal"
		switch priority {
		case models.Low:
			urgency = "low"
		case models.High:
			urgency = "critical"
		}

		cmd := exec.Command("dunstify",
			"-u", urgency,
			"-a", "Nancy",
			title,
			message,
		)
		return cmd.Run()
	}

	return fmt.Errorf("no suitable notification command found (tried notify-send, dunstify)")
}

// sendMacOSDesktopNotification sends a desktop notification on macOS
func (n *Notifier) sendMacOSDesktopNotification(title, message string, priority models.Priority) error {
	// Try terminal-notifier first (if installed)
	if _, err := exec.LookPath("terminal-notifier"); err == nil {
		args := []string{
			"-title", title,
			"-message", message,
			"-sender", "com.apple.Terminal", // Use Terminal as sender
		}

		// Add sound for high priority
		if priority == models.High {
			args = append(args, "-sound", "default")
		}

		cmd := exec.Command("terminal-notifier", args...)
		return cmd.Run()
	}

	// Use built-in osascript as fallback
	if _, err := exec.LookPath("osascript"); err == nil {
		script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
		if priority == models.High {
			script = fmt.Sprintf(`display notification "%s" with title "%s" sound name "default"`, message, title)
		}

		cmd := exec.Command("osascript", "-e", script)
		return cmd.Run()
	}

	return fmt.Errorf("no suitable notification command found (tried terminal-notifier, osascript)")
}

// sendWindowsDesktopNotification sends a desktop notification on Windows
func (n *Notifier) sendWindowsDesktopNotification(title, message string, priority models.Priority) error {
	// Use PowerShell to show Windows Toast notification
	script := fmt.Sprintf(`
		[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null;
		[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null;
		$template = @"
<toast>
	<visual>
		<binding template="ToastGeneric">
			<text>%s</text>
			<text>%s</text>
		</binding>
	</visual>
</toast>
"@;
		$xml = New-Object Windows.Data.Xml.Dom.XmlDocument;
		$xml.LoadXml($template);
		$toast = New-Object Windows.UI.Notifications.ToastNotification $xml;
		[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Nancy").Show($toast);
	`, title, message)

	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

// sendTerminalBell sends a terminal bell notification
func (n *Notifier) sendTerminalBell(title, message string) error {
	// Print notification to stderr with bell character
	fmt.Fprintf(os.Stderr, "\a🔔 %s: %s\n", title, message)
	return nil
}

// logNotification logs the notification to a file or stderr
func (n *Notifier) logNotification(title, message string) error {
	logMessage := fmt.Sprintf("[NOTIFICATION] %s: %s", title, message)
	
	if n.logFile != "" {
		// TODO: Implement file logging
		// For now, just print to stderr
		fmt.Fprintln(os.Stderr, logMessage)
	} else {
		fmt.Fprintln(os.Stderr, logMessage)
	}
	
	return nil
}

// SetLogFile sets the file path for log-only notifications
func (n *Notifier) SetLogFile(path string) {
	n.logFile = path
}

// TestNotification sends a test notification to verify the system works
func (n *Notifier) TestNotification() error {
	return n.Send(
		"Nancy Test Notification",
		"If you see this, notifications are working correctly! 🎉",
		models.Medium,
	)
}

// GetMethod returns the current notification method
func (n *Notifier) GetMethod() NotificationMethod {
	return n.method
}

// SetMethod sets the notification method
func (n *Notifier) SetMethod(method NotificationMethod) {
	n.method = method
}

// GetMethodName returns a human-readable name for a notification method
func GetMethodName(method NotificationMethod) string {
	switch method {
	case DesktopNotification:
		return "Desktop Notification"
	case TerminalBell:
		return "Terminal Bell"
	case LogOnly:
		return "Log Only"
	default:
		return "Unknown"
	}
}

// GetAvailableMethods returns a list of available notification methods for the current system
func GetAvailableMethods() []NotificationMethod {
	var methods []NotificationMethod

	// Check if desktop notifications are available
	switch runtime.GOOS {
	case "linux":
		if _, err := exec.LookPath("notify-send"); err == nil {
			methods = append(methods, DesktopNotification)
		} else if _, err := exec.LookPath("dunstify"); err == nil {
			methods = append(methods, DesktopNotification)
		}
	case "darwin":
		if _, err := exec.LookPath("osascript"); err == nil {
			methods = append(methods, DesktopNotification)
		}
	case "windows":
		if _, err := exec.LookPath("powershell"); err == nil {
			methods = append(methods, DesktopNotification)
		}
	}

	// Terminal bell is always available
	methods = append(methods, TerminalBell)
	methods = append(methods, LogOnly)

	return methods
}
