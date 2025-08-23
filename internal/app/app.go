package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ivyascorp-net/nagging-nancy/internal/models"
)

// App represents the main application instance
type App struct {
	config *Config
	store  *models.Store
}

// New creates a new application instance
func New() (*App, error) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize data store
	store, err := models.NewStore(config.GetDataDir())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize store: %w", err)
	}

	app := &App{
		config: config,
		store:  store,
	}

	return app, nil
}

// GetConfig returns the application configuration
func (a *App) GetConfig() *Config {
	return a.config
}

// GetStore returns the data store
func (a *App) GetStore() *models.Store {
	return a.store
}

// AddReminder adds a new reminder to the store
func (a *App) AddReminder(title string, dueTime interface{}, priority models.Priority) (*models.Reminder, error) {
	// TODO: Parse dueTime (could be time.Time, string, etc.)
	// For now, assume it's already a time.Time
	reminder := models.NewReminder(title, dueTime.(time.Time), priority)

	if err := a.store.Add(reminder); err != nil {
		return nil, fmt.Errorf("failed to add reminder: %w", err)
	}

	return reminder, nil
}

// GetReminders returns reminders with optional filtering
func (a *App) GetReminders(filter *models.FilterOptions) []*models.Reminder {
	return a.store.GetAll(filter)
}

// CompleteReminder marks a reminder as completed
func (a *App) CompleteReminder(id string) error {
	return a.store.CompleteReminder(id)
}

// DeleteReminder removes a reminder
func (a *App) DeleteReminder(id string) error {
	return a.store.Delete(id)
}

// GetStats returns application statistics
func (a *App) GetStats() (total, active, completed, overdue int) {
	return a.store.Count()
}

// Cleanup performs maintenance tasks
func (a *App) Cleanup() error {
	return a.store.Cleanup()
}

// Version information
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// GetVersion returns version information
func GetVersion() string {
	return fmt.Sprintf("nancy %s (built %s, commit %s)", Version, BuildTime, GitCommit)
}

// InitApp initializes the application directories and files
func InitApp() error {
	config := NewDefaultConfig()

	// Create config directory
	configDir := config.GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create data directory
	dataDir := config.GetDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Save default config if it doesn't exist
	configPath := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := config.Save(); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
	}

	return nil
}
