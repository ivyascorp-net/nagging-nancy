package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/ivyascorp-net/nagging-nancy/internal/app"
	"github.com/ivyascorp-net/nagging-nancy/internal/models"
	"github.com/ivyascorp-net/nagging-nancy/internal/tui/components"
)

// Model represents the application state for the TUI
type Model struct {
	store        *models.Store
	config       *app.Config
	width        int
	height       int
	reminders    []*models.Reminder
	cursor       int
	showHelp     bool
	filter       *models.FilterOptions
	quitting     bool
	editing      bool
	editForm     *components.EditForm
}

// NewModel creates a new TUI model
func NewModel(store *models.Store, config *app.Config) Model {
	filter := &models.FilterOptions{
		ShowCompleted: false,
	}

	model := Model{
		store:     store,
		config:    config,
		reminders: store.GetAll(filter),
		cursor:    0,
		showHelp:  false,
		filter:    filter,
		quitting:  false,
	}

	return model
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// refreshReminders loads reminders from store
func (m *Model) refreshReminders() {
	m.reminders = m.store.GetAll(m.filter)
	if m.cursor >= len(m.reminders) && len(m.reminders) > 0 {
		m.cursor = len(m.reminders) - 1
	}
	if len(m.reminders) == 0 {
		m.cursor = 0
	}
}

// getCurrentReminder returns the currently selected reminder
func (m Model) getCurrentReminder() *models.Reminder {
	if len(m.reminders) == 0 || m.cursor < 0 || m.cursor >= len(m.reminders) {
		return nil
	}
	return m.reminders[m.cursor]
}