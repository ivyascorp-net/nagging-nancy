package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ivyascorp-net/nagging-nancy/internal/models"
	"github.com/ivyascorp-net/nagging-nancy/internal/utils"
)

type EditForm struct {
	reminder    *models.Reminder
	titleInput  textinput.Model
	timeInput   textinput.Model
	dateInput   textinput.Model
	focused     int
	done        bool
	cancelled   bool
	width       int
	height      int
	errorMsg    string
}

const (
	titleField = 0
	timeField  = 1
	dateField  = 2
	numFields  = 3
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func NewEditForm(reminder *models.Reminder) *EditForm {
	ti := textinput.New()
	ti.Placeholder = "Title"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50
	ti.SetValue(reminder.Title)

	timeInput := textinput.New()
	timeInput.Placeholder = "Time (e.g., 3pm, 14:30)"
	timeInput.CharLimit = 20
	timeInput.Width = 30
	timeInput.SetValue(reminder.DueTime.Format("3:04 PM"))

	dateInput := textinput.New()
	dateInput.Placeholder = "Date (e.g., tomorrow, 2024-03-20)"
	dateInput.CharLimit = 30
	dateInput.Width = 30
	dateInput.SetValue(reminder.DueTime.Format("2006-01-02"))

	return &EditForm{
		reminder:   reminder,
		titleInput: ti,
		timeInput:  timeInput,
		dateInput:  dateInput,
		focused:    titleField,
	}
}

func (f *EditForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *EditForm) Update(msg tea.Msg) (*EditForm, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			f.cancelled = true
			return f, nil

		case "enter":
			if f.focused == numFields-1 {
				// Last field, submit form
				return f.submit()
			} else {
				// Move to next field
				f.nextField()
			}

		case "shift+tab":
			f.prevField()

		case "tab":
			f.nextField()
		}
	}

	// Update the focused input
	switch f.focused {
	case titleField:
		f.titleInput, cmd = f.titleInput.Update(msg)
	case timeField:
		f.timeInput, cmd = f.timeInput.Update(msg)
	case dateField:
		f.dateInput, cmd = f.dateInput.Update(msg)
	}

	cmds = append(cmds, cmd)
	return f, tea.Batch(cmds...)
}

func (f *EditForm) View() string {
	var s strings.Builder

	s.WriteString(focusedStyle.Render("✏️  Edit Reminder\n\n"))

	// Title field
	titleLabel := "Title:"
	if f.focused == titleField {
		titleLabel = focusedStyle.Render("> " + titleLabel)
	} else {
		titleLabel = blurredStyle.Render("  " + titleLabel)
	}
	s.WriteString(titleLabel + "\n")
	s.WriteString(f.titleInput.View() + "\n\n")

	// Time field
	timeLabel := "Time:"
	if f.focused == timeField {
		timeLabel = focusedStyle.Render("> " + timeLabel)
	} else {
		timeLabel = blurredStyle.Render("  " + timeLabel)
	}
	s.WriteString(timeLabel + "\n")
	s.WriteString(f.timeInput.View() + "\n\n")

	// Date field
	dateLabel := "Date:"
	if f.focused == dateField {
		dateLabel = focusedStyle.Render("> " + dateLabel)
	} else {
		dateLabel = blurredStyle.Render("  " + dateLabel)
	}
	s.WriteString(dateLabel + "\n")
	s.WriteString(f.dateInput.View() + "\n\n")

	// Error message
	if f.errorMsg != "" {
		s.WriteString(errorStyle.Render("Error: " + f.errorMsg + "\n\n"))
	}

	// Help text
	help := helpStyle.Render("tab: next field • shift+tab: prev field • enter: save • esc: cancel")
	s.WriteString(help)

	return s.String()
}

func (f *EditForm) nextField() {
	f.focused = (f.focused + 1) % numFields
	f.updateFieldFocus()
}

func (f *EditForm) prevField() {
	f.focused = (f.focused - 1 + numFields) % numFields
	f.updateFieldFocus()
}

func (f *EditForm) updateFieldFocus() {
	f.titleInput.Blur()
	f.timeInput.Blur()
	f.dateInput.Blur()

	switch f.focused {
	case titleField:
		f.titleInput.Focus()
	case timeField:
		f.timeInput.Focus()
	case dateField:
		f.dateInput.Focus()
	}
}

func (f *EditForm) submit() (*EditForm, tea.Cmd) {
	f.errorMsg = ""

	// Get values
	title := strings.TrimSpace(f.titleInput.Value())
	timeStr := strings.TrimSpace(f.timeInput.Value())
	dateStr := strings.TrimSpace(f.dateInput.Value())

	// Validate title
	if title == "" {
		f.errorMsg = "Title cannot be empty"
		return f, nil
	}

	// Parse time
	var newTime time.Time
	var err error

	// Try parsing the time string
	if timeStr != "" {
		parsedTime, err := utils.ParseTimeString(timeStr)
		if err != nil {
			f.errorMsg = fmt.Sprintf("Invalid time format: %s", err.Error())
			return f, nil
		}
		newTime = parsedTime
	} else {
		newTime = f.reminder.DueTime
	}

	// Parse date
	var newDate time.Time
	if dateStr != "" {
		// Try different date formats
		dateFormats := []string{
			"2006-01-02",  // 2024-03-20
			"01/02/2006",  // 03/20/2024
			"01-02-2006",  // 03-20-2024
			"Jan 2, 2006", // Mar 20, 2024
			"Jan 2 2006",  // Mar 20 2024
			"2 Jan 2006",  // 20 Mar 2006
		}

		// Handle relative dates
		switch strings.ToLower(dateStr) {
		case "today":
			newDate = time.Now()
		case "tomorrow":
			newDate = time.Now().AddDate(0, 0, 1)
		default:
			// Try parsing as explicit date
			for _, format := range dateFormats {
				if newDate, err = time.Parse(format, dateStr); err == nil {
					break
				}
			}
			if err != nil {
				f.errorMsg = fmt.Sprintf("Invalid date format: %s", dateStr)
				return f, nil
			}
		}
	} else {
		newDate = f.reminder.DueTime
	}

	// Combine date and time
	finalTime := time.Date(
		newDate.Year(), newDate.Month(), newDate.Day(),
		newTime.Hour(), newTime.Minute(), 0, 0,
		time.Local,
	)

	// Update the reminder
	f.reminder.Title = title
	f.reminder.DueTime = finalTime
	f.reminder.UpdatedAt = time.Now()

	f.done = true
	return f, nil
}

func (f *EditForm) Done() bool {
	return f.done
}

func (f *EditForm) Cancelled() bool {
	return f.cancelled
}

func (f *EditForm) GetReminder() *models.Reminder {
	return f.reminder
}