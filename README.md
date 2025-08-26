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

### Homebrew (Recommended)
```bash
# Add the tap and install (macOS/Linux)
brew tap ivyascorp-net/nagging-nancy
brew install nagging-nancy

# Or install from local formula
brew install ./nagging-nancy.rb
```

After installation with Homebrew, Nancy will automatically detect and use the best notification method for your platform. To install enhanced notification support:

**macOS:**
```bash
# Optional: Install terminal-notifier for better notifications
brew install terminal-notifier
```

**Linux:**
```bash
# Install notification dependencies (auto-detected by package manager)
sudo apt install libnotify-bin    # Ubuntu/Debian
sudo dnf install libnotify         # Fedora
sudo pacman -S libnotify          # Arch Linux
```

### Download Binary
```bash
# Download latest release for your platform
curl -L https://github.com/ivyascorp-net/nagging-nancy/releases/latest/download/nancy-linux-amd64 -o nancy
chmod +x nancy
sudo mv nancy /usr/local/bin/
```

### Build from Source
```bash
git clone https://github.com/ivyascorp-net/nagging-nancy.git
cd nagging-nancy
make build
sudo cp build/nancy /usr/local/bin/
```

### Other Package Managers
```bash
# Go install
go install github.com/ivyascorp-net/nagging-nancy/cmd/nancy@latest
```

### Post-Installation Setup
After installing Nancy, set up notification dependencies:
```bash
# Automatic setup script (detects your platform)
make install-notifications

# Or manually run the script
./scripts/install-notifications.sh

# Test your notification setup
nancy test notification
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

# Delete reminders
nancy delete 2               # Delete reminder with ID 2

# Background daemon
nancy daemon start           # Start background notifications
nancy daemon start --foreground  # Run in foreground for debugging
nancy daemon stop            # Stop background notifications
nancy daemon status          # Check daemon status
nancy daemon restart         # Restart daemon

# Test notifications
nancy test notification      # Send test notification

# Setup notifications for your platform  
make install-notifications   # Auto-install notification dependencies
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

# With tags
nancy add "Review code" --tags work,coding --priority medium

# Basic natural language support
nancy add "Doctor appointment tomorrow at 2pm"
nancy add "Team meeting today at 3:30pm"
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
nancy list --tags work,urgent --all
```

### Configuration
Configuration is managed through the config file located at:
- **Linux/macOS**: `~/.config/nancy/config.yaml`
- **Windows**: `%APPDATA%/nancy/config.yaml`

Nancy automatically creates a default configuration file on first run. Edit the file directly to customize settings.

## 🔧 Configuration Files

Nancy stores its files in:
- **Config**: `~/.config/nancy/config.yaml` (Linux/macOS) or `%APPDATA%/nancy/config.yaml` (Windows)
- **Data**: `~/.local/share/nancy/` (Linux) or `%LOCALAPPDATA%/nancy/` (Windows)
- **macOS**: `~/Library/Application Support/nancy/`

Example configuration:
```yaml
# Default settings for new reminders
default:
  priority: medium          # low, medium, high
  advance_minutes: 10       # Default notification advance time

# Notification settings
notifications:
  enabled: true             # Enable desktop notifications
  sound: true               # Play notification sound
  advance_minutes: 15       # How many minutes before due time to notify
  quiet_hours: true         # Respect working hours for notifications

# Appearance settings
appearance:
  theme: auto               # light, dark, auto
  show_completed: false     # Show completed tasks in main list
  compact_mode: false       # Use compact display mode
  show_icons: true          # Show priority and status icons

# Working hours (for quiet notifications)
workhours:
  enabled: true             # Enable working hours
  start: "09:00"            # Work start time (24-hour format)
  end: "17:00"              # Work end time (24-hour format)
  quiet_outside: true       # Quiet notifications outside work hours
  timezone: "Local"         # Timezone (Local or specific timezone)

# Background daemon settings
daemon:
  check_interval: 5         # Check for due reminders every N minutes
  auto_start: false         # Auto-start daemon on system boot
  log_level: "info"         # Logging level: debug, info, warn, error
```

Your reminders and configuration are stored locally:
- **Configuration**: Stored in OS-appropriate config directories 
- **Data**: Stored in OS-appropriate data directories (separate from config)
- **Privacy**: Nancy never sends your data anywhere - everything stays on your machine

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
nancy daemon start --interval 2m0s

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

# Check data directory for daemon files
# Linux: ~/.local/share/nancy/
# macOS: ~/Library/Application Support/nancy/
# Windows: %LOCALAPPDATA%/nancy/
```

### Notification Issues
```bash
# Test notifications
nancy test notification

# Check available notification methods
nancy test notification  # Shows available methods at the end

# Auto-install notification dependencies for your platform
make install-notifications
# or
./scripts/install-notifications.sh

# Manual installation:
# Linux: Install notification dependencies
sudo apt install libnotify-bin  # Ubuntu/Debian (notify-send)
sudo dnf install libnotify       # Fedora (notify-send)  
sudo pacman -S libnotify        # Arch Linux (notify-send)
# or
sudo apt install dunst          # dunstify alternative

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
- Your operating system and version
- Nancy version (`nancy version`)  
- Steps to reproduce the issue
- Expected vs actual behavior
- For notification issues: Output of `nancy test notification`
- For daemon issues: Output of `nancy daemon status`
- Relevant log files if available

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss)  
- CLI powered by [Cobra](https://github.com/spf13/cobra)
- Configuration with [Viper](https://github.com/spf13/viper)
- Cross-platform notifications: notify-send, terminal-notifier, osascript, PowerShell

---

**Made with ❤️ for terminal lovers who need gentle reminders to stay productive.**
