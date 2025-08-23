package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginLeft(2)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	statusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("235"))

	cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))

	completedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Strikethrough(true)
)

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return "Thanks for using Nagging Nancy! 👋\n"
	}

	if m.showHelp {
		return m.helpView()
	}

	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("📝 Nagging Nancy"))
	s.WriteString(fmt.Sprintf(" - %s\n\n", time.Now().Format("Monday, January 2, 2006")))

	if len(m.reminders) == 0 {
		s.WriteString("🎉 All caught up! No active reminders.\n\n")
		s.WriteString("Press 'q' to quit, '?' for help\n")
		return s.String()
	}

	// List reminders
	for i, reminder := range m.reminders {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		status := "●"
		if reminder.Completed {
			status = "✓"
		}

		line := fmt.Sprintf("%s %s %s %s - %s",
			cursor,
			status,
			reminder.Priority.Icon(),
			reminder.Title,
			reminder.FormattedDueTime(),
		)

		if reminder.Completed {
			// Apply strikethrough to entire line, then color the cursor separately
			styledLine := completedStyle.Render(line)
			// Replace the plain cursor with styled cursor after strikethrough
			if m.cursor == i {
				styledLine = strings.Replace(styledLine, ">", cursorStyle.Render(">"), 1)
			}
			line = styledLine
		} else {
			// Apply cursor styling for non-completed items
			if m.cursor == i {
				line = strings.Replace(line, ">", cursorStyle.Render(">"), 1)
			}
			
			if reminder.IsOverdue() {
				line = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(line + " ⚠️ OVERDUE")
			} else if reminder.IsDueSoon() {
				line = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(line + " ⏰ DUE SOON")
			}
		}

		s.WriteString(line)
		s.WriteString("\n")
	}

	// Status bar
	s.WriteString("\n")
	s.WriteString(m.statusBarView())

	return s.String()
}

func (m Model) helpView() string {
	help := `📝 Nagging Nancy - Help

Navigation:
  ↑/k      Move up
  ↓/j      Move down
  
Actions:
  space    Toggle reminder completion
  d        Delete selected reminder
  r        Refresh list
  f        Toggle show completed
  
Other:
  ?/h      Show/hide help
  q        Quit

Press any key to return...`

	return help
}

func (m Model) statusBarView() string {
	total, active, completed, overdue := m.store.Count()

	status := fmt.Sprintf("Total: %d | Active: %d | Completed: %d | Overdue: %d",
		total, active, completed, overdue)

	controls := "space=toggle d=delete f=filter ?=help q=quit"

	// Pad to full width
	padding := m.width - len(status) - len(controls)
	if padding < 0 {
		padding = 0
	}

	statusBar := status + strings.Repeat(" ", padding) + controls
	return statusBarStyle.Render(statusBar)
}
