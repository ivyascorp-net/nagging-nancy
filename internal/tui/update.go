package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivyascorp-net/nagging-nancy/internal/tui/components"
)

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle edit form updates when in edit mode
	if m.editing && m.editForm != nil {
		var cmd tea.Cmd
		m.editForm, cmd = m.editForm.Update(msg)
		
		if m.editForm.Done() {
			// Save the edited reminder
			reminder := m.editForm.GetReminder()
			if err := m.store.Update(reminder); err == nil {
				m.refreshReminders()
			}
			m.editing = false
			m.editForm = nil
		} else if m.editForm.Cancelled() {
			// Cancel editing
			m.editing = false
			m.editForm = nil
		}
		
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// If showing help, any key press should hide help
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "?", "h":
			m.showHelp = true
			return m, nil

		case "j", "down":
			if len(m.reminders) > 0 {
				m.cursor++
				if m.cursor >= len(m.reminders) {
					m.cursor = 0
				}
			}
			return m, nil

		case "k", "up":
			if len(m.reminders) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.reminders) - 1
				}
			}
			return m, nil

		case " ":
			// Toggle completion
			if current := m.getCurrentReminder(); current != nil {
				m.store.ToggleReminder(current.ID)
				m.refreshReminders()
			}
			return m, nil

		case "d":
			// Delete current reminder
			if current := m.getCurrentReminder(); current != nil {
				m.store.Delete(current.ID)
				m.refreshReminders()
			}
			return m, nil

		case "e":
			if current := m.getCurrentReminder(); current != nil {
				reminder, err := m.store.Get(current.ID)
				if err != nil {
					return m, nil
				}
				m.editing = true
				m.editForm = components.NewEditForm(reminder)
				return m, m.editForm.Init()
			}
			return m, nil

		case "r":
			// Refresh reminders
			m.refreshReminders()
			return m, nil

		case "f":
			// Toggle show completed filter
			m.filter.ShowCompleted = !m.filter.ShowCompleted
			m.refreshReminders()
			return m, nil
		}
	}

	return m, nil
}
