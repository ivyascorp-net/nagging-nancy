# 📝 Nagging Nancy

> Your friendly, persistent terminal-based reminders app

Nancy is a fast, lightweight terminal application that helps you manage reminders and tasks without leaving your command line. Built with Go and Bubble Tea for a smooth, responsive experience.

## ✨ Features

- 🚀 **Lightning fast** - Single binary, instant startup
- 🎨 **Beautiful TUI** - Clean, interactive terminal interface
- ⚡ **Quick CLI commands** - Add reminders without opening the interface
- 🔔 **Smart notifications** - Cross-platform desktop notifications
- 📅 **Natural language** - "Call mom tomorrow at 2pm"
- 🎯 **Priority system** - High, medium, low priority tasks
- 💾 **Local storage** - Your data stays on your machine
- 🌙 **Background daemon** - Persistent reminder monitoring
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
nancy daemon stop            # Stop background notifications
nancy daemon status          # Check daemon status
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

## 🔄 Background Daemon

Start the background daemon to receive notifications even when Nancy isn't running:

```bash
# Start daemon
nancy daemon start

# Check status
nancy daemon status

# View logs
nancy daemon logs

# Stop daemon
nancy daemon stop
```

The daemon will:
- Send desktop notifications for due reminders
- Respect your working hours settings
- Run quietly in the background
- Auto-start on system boot (optional)

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

## 🐛 Bug Reports

Found a bug? Please open an issue with:
- Your operating system
- Nancy version (`nancy version`)
- Steps to reproduce
- Expected vs actual behavior

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- CLI powered by [Cobra](https://github.com/spf13/cobra)
- Notifications via [Beeep](https://github.com/gen2brain/beeep)

---

**Made with ❤️ for terminal lovers who need gentle reminders to stay productive.**
