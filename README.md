# 📝 Nagging Nancy

> Your friendly, persistent terminal-based reminders app

Nancy is a fast, lightweight terminal application that helps you manage reminders and tasks without leaving your command line. Built with Go and Bubble Tea for a smooth, responsive experience.

## ✨ Features

- 🚀 **Lightning fast** - Single binary, instant startup
- 🎨 **Beautiful TUI** - Clean, interactive terminal interface
- ⚡ **Quick CLI commands** - Add reminders without opening the interface
- 🔔 **Smart notifications** - Cross-platform desktop notifications with intelligent timing
- 📅 **Natural language** - "Call mom tomorrow at 2pm"
- 🎯 **Priority system** - High, medium, low priority tasks
- 💾 **Local storage** - Your data stays on your machine
- 🌙 **Background daemon** - Persistent reminder monitoring with configurable intervals
- ⌨️ **Keyboard-first** - Navigate entirely with keyboard shortcuts

## 📦 Installation

### Download Binary (Recommended)
```bash
# Download latest release for your platform
curl -L https://github.com/yourusername/nagging-nancy/releases/latest/download/nancy-linux-amd64 -o nancy
chmod +x nancy
sudo mv nancy /usr/local/bin/
```

### Build from Source
```bash
git clone https://github.com/yourusername/nagging-nancy.git
cd nagging-nancy
go build -o nancy cmd/nancy/main.go
sudo mv nancy /usr/local/bin/
```

### Package Managers
```bash
# Homebrew (macOS/Linux)
brew install nagging-nancy

# Go install
go install github.com/yourusername/nagging-nancy/cmd/nancy@latest
```

## 🚀 Quick Start

### Launch Interactive Interface
```bash
nancy
# or
nancy tui
```

### CLI Commands
```bash
# Add reminders
nancy add "Call mom"
nancy add "Buy groceries" --time "5pm" --priority high
nancy add "Meeting tomorrow at 10am"

# List reminders
nancy list                    # All active reminders
nancy list --today           # Today's reminders only
nancy list --priority high   # High priority only

# Complete tasks
nancy complete 1             # Complete reminder with ID 1
nancy done 1                 # Alias for complete

# Delete reminders
nancy delete 2               # Delete reminder with ID 2
nancy rm 2                   # Alias for delete

# Background daemon
nancy daemon start           # Start background notifications
nancy daemon start --foreground  # Run in foreground for debugging
nancy daemon stop            # Stop background notifications
nancy daemon status          # Check daemon status
nancy daemon restart         # Restart daemon

# Test notifications
nancy test notification      # Send test notification
```

## ⌨️ TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `a` / `n` | Add new reminder |
| `space` | Toggle complete |
| `d` | Delete reminder |
| `e` | Edit reminder |
| `f` | Filter reminders |
| `h` / `?` | Help screen |
| `q` / `ctrl+c` | Quit |
| `tab` | Switch between sections |

## 📋 Usage Examples

### Adding Reminders
```bash
# Simple reminder
nancy add "Water the plants"

# With specific time
nancy add "Team standup" --time "9:30am"

# With priority and date
nancy add "Submit report" --date "2024-03-20" --priority high

# Natural language parsing
nancy add "Doctor appointment tomorrow at 2pm"
nancy add "Call John today at 3:30"
nancy add "Grocery shopping this Friday at 6pm"
```

### Listing and Filtering
```bash
# View different sets of reminders
nancy list --today           # Today's tasks
nancy list --week            # This week
nancy list --overdue         # Overdue items
nancy list --completed       # Completed tasks
nancy list --priority high   # High priority only

# Combine filters
nancy list --today --priority high
```

### Configuration
```bash
# Set default priority
nancy config set default.priority medium

# Configure notification timing
nancy config set notifications.advance_minutes 15

# Set working hours
nancy config set workhours.start "09:00"
nancy config set workhours.end "17:00"

# View current config
nancy config show
```

## 🔧 Configuration

Nancy stores its configuration in:
- **Linux/macOS**: `~/.config/nancy/config.yaml`
- **Windows**: `%APPDATA%/nancy/config.yaml`

Example configuration:
```yaml
default:
  priority: medium
  advance_minutes: 10

notifications:
  enabled: true
  sound: true
  advance_minutes: 15

appearance:
  theme: auto  # light, dark, auto
  show_completed: false

workhours:
  enabled: true
  start: "09:00"
  end: "17:00"
  quiet_outside: true
```

## 📁 Data Storage

Your reminders are stored locally in:
- **Linux/macOS**: `~/.config/nancy/reminders.json`
- **Windows**: `%APPDATA%/nancy/reminders.json`

Nancy never sends your data anywhere - everything stays on your machine.

## 🔔 Notification System

Nancy includes a comprehensive cross-platform notification system:

### Supported Platforms
- **Linux**: notify-send (libnotify) or dunstify
- **macOS**: osascript (built-in) or terminal-notifier  
- **Windows**: PowerShell toast notifications

### Notification Types
```bash
# Test your notification system
nancy test notification

# The daemon sends different types of notifications:
# - 📅 Due Today: Sent once per day for today's reminders
# - ⏰ Due Soon: Sent 15 minutes before due time
# - ⚠️  Overdue: Sent hourly until reminder is completed
```

### Fallback Methods
If desktop notifications aren't available, Nancy automatically falls back to:
1. **Terminal Bell** - Audible bell with message in terminal
2. **Log Only** - Messages logged to stderr

### Priority Levels
Notifications respect reminder priorities:
- **High Priority** 🔴: Critical/urgent notification with sound
- **Medium Priority** 🟡: Normal notification  
- **Low Priority** 🟢: Low urgency notification

## 🔄 Background Daemon

Start the background daemon to receive notifications even when Nancy isn't running:

```bash
# Start daemon (runs in background)
nancy daemon start

# Start daemon in foreground (for debugging)
nancy daemon start --foreground

# Start daemon with custom check interval
nancy daemon start --interval 2m

# Check daemon status
nancy daemon status

# Restart daemon
nancy daemon restart

# Stop daemon
nancy daemon stop

# Test notification system
nancy test notification
```

The daemon will:
- Monitor reminders every 5 minutes (configurable)
- Send desktop notifications for:
  - **Overdue reminders** - Every hour until completed
  - **Due soon** - 15 minutes before due time 
  - **Due today** - Once per day for today's reminders
- Use PID file to prevent multiple instances
- Handle graceful shutdown via signals
- Fall back to terminal notifications if desktop unavailable

## 🎨 Screenshots

```
┌─ Nagging Nancy ─────────────────────────────────┐
│ Today - March 15, 2024                  📅      │
├────────────────────────────────────────────────┤
│ ● Call mom                           2:00 PM 🔴 │
│ ● Buy groceries                      4:30 PM 🟡 │
│ ● Team standup                       9:30 AM 🟢 │
│ ✓ Morning jog                        7:00 AM    │
│                                                 │
│ 3 active • 1 completed • 0 overdue             │
│ [a]dd [d]elete [q]uit [?]help                   │
└────────────────────────────────────────────────┘
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🔧 Troubleshooting

### Daemon Issues
```bash
# Check if daemon is running
nancy daemon status

# Stop any stuck processes
nancy daemon stop

# Start daemon in foreground to see logs
nancy daemon start --foreground

# Check PID file location
# Linux/macOS: ~/.config/nancy/daemon.pid
# Windows: %APPDATA%/nancy/daemon.pid
```

### Notification Issues
```bash
# Test notifications
nancy test notification

# Check available notification methods
nancy test notification  # Shows available methods at the end

# Linux: Install notification dependencies
sudo apt install libnotify-bin  # notify-send
# or
sudo apt install dunst          # dunstify

# macOS: Install optional terminal-notifier
brew install terminal-notifier

# Windows: PowerShell should be available by default
```

### Common Issues
- **No notifications**: Run `nancy test notification` to verify system
- **Daemon won't start**: Check if another instance is running with `nancy daemon status`
- **Permission errors**: Ensure Nancy has permission to create files in config directory

## 🐛 Bug Reports

Found a bug? Please open an issue with:
- Your operating system
- Nancy version (`nancy version`)
- Steps to reproduce
- Expected vs actual behavior
- Output of `nancy test notification` (for notification issues)
- Output of `nancy daemon status` (for daemon issues)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- CLI powered by [Cobra](https://github.com/spf13/cobra)
- Cross-platform notifications (notify-send, osascript, PowerShell)

---

**Made with ❤️ for terminal lovers who need gentle reminders to stay productive.**
