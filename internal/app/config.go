package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	DataDir       string             `mapstructure:"data_dir"`
	Default       DefaultConfig      `mapstructure:"default"`
	Notifications NotificationConfig `mapstructure:"notifications"`
	Appearance    AppearanceConfig   `mapstructure:"appearance"`
	WorkHours     WorkHoursConfig    `mapstructure:"workhours"`
	Daemon        DaemonConfig       `mapstructure:"daemon"`
}

// DefaultConfig holds default settings for new reminders
type DefaultConfig struct {
	Priority       string `mapstructure:"priority"`
	AdvanceMinutes int    `mapstructure:"advance_minutes"`
}

// NotificationConfig holds notification settings
type NotificationConfig struct {
	Enabled        bool `mapstructure:"enabled"`
	Sound          bool `mapstructure:"sound"`
	AdvanceMinutes int  `mapstructure:"advance_minutes"`
	QuietHours     bool `mapstructure:"quiet_hours"`
}

// AppearanceConfig holds UI appearance settings
type AppearanceConfig struct {
	Theme         string `mapstructure:"theme"` // "light", "dark", "auto"
	ShowCompleted bool   `mapstructure:"show_completed"`
	CompactMode   bool   `mapstructure:"compact_mode"`
	ShowIcons     bool   `mapstructure:"show_icons"`
}

// WorkHoursConfig defines working hours for quiet notifications
type WorkHoursConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Start        string `mapstructure:"start"` // "09:00"
	End          string `mapstructure:"end"`   // "17:00"
	QuietOutside bool   `mapstructure:"quiet_outside"`
	Timezone     string `mapstructure:"timezone"`
}

// DaemonConfig holds daemon-specific settings
type DaemonConfig struct {
	CheckInterval int    `mapstructure:"check_interval"` // minutes
	AutoStart     bool   `mapstructure:"auto_start"`
	LogLevel      string `mapstructure:"log_level"`
}

// getConfigDir returns the appropriate config directory for the OS
func getConfigDir() string {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	default: // linux and other unix-like systems
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("HOME"), ".config")
		}
	}

	return filepath.Join(configDir, "nancy")
}

// getDataDir returns the appropriate data directory for the OS
func getDataDir() string {
	var dataDir string

	switch runtime.GOOS {
	case "windows":
		dataDir = os.Getenv("LOCALAPPDATA")
		if dataDir == "" {
			dataDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
	case "darwin":
		dataDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	default: // linux and other unix-like systems
		dataDir = os.Getenv("XDG_DATA_HOME")
		if dataDir == "" {
			dataDir = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
	}

	return filepath.Join(dataDir, "nancy")
}

// DefaultConfig returns a config with sensible defaults
func NewDefaultConfig() *Config {
	return &Config{
		DataDir: getDataDir(),
		Default: DefaultConfig{
			Priority:       "medium",
			AdvanceMinutes: 10,
		},
		Notifications: NotificationConfig{
			Enabled:        true,
			Sound:          true,
			AdvanceMinutes: 15,
			QuietHours:     true,
		},
		Appearance: AppearanceConfig{
			Theme:         "auto",
			ShowCompleted: false,
			CompactMode:   false,
			ShowIcons:     true,
		},
		WorkHours: WorkHoursConfig{
			Enabled:      true,
			Start:        "09:00",
			End:          "17:00",
			QuietOutside: true,
			Timezone:     "Local",
		},
		Daemon: DaemonConfig{
			CheckInterval: 5, // check every 5 minutes
			AutoStart:     false,
			LogLevel:      "info",
		},
	}
}

// LoadConfig loads configuration from file or creates default if not found
func LoadConfig() (*Config, error) {
	configDir := getConfigDir()

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Setup viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Set default values
	config := NewDefaultConfig()
	setViperDefaults(config)

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default
			if err := saveDefaultConfig(configDir); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal into config struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// setViperDefaults sets default values in viper
func setViperDefaults(config *Config) {
	viper.SetDefault("data_dir", config.DataDir)
	viper.SetDefault("default.priority", config.Default.Priority)
	viper.SetDefault("default.advance_minutes", config.Default.AdvanceMinutes)
	viper.SetDefault("notifications.enabled", config.Notifications.Enabled)
	viper.SetDefault("notifications.sound", config.Notifications.Sound)
	viper.SetDefault("notifications.advance_minutes", config.Notifications.AdvanceMinutes)
	viper.SetDefault("notifications.quiet_hours", config.Notifications.QuietHours)
	viper.SetDefault("appearance.theme", config.Appearance.Theme)
	viper.SetDefault("appearance.show_completed", config.Appearance.ShowCompleted)
	viper.SetDefault("appearance.compact_mode", config.Appearance.CompactMode)
	viper.SetDefault("appearance.show_icons", config.Appearance.ShowIcons)
	viper.SetDefault("workhours.enabled", config.WorkHours.Enabled)
	viper.SetDefault("workhours.start", config.WorkHours.Start)
	viper.SetDefault("workhours.end", config.WorkHours.End)
	viper.SetDefault("workhours.quiet_outside", config.WorkHours.QuietOutside)
	viper.SetDefault("workhours.timezone", config.WorkHours.Timezone)
	viper.SetDefault("daemon.check_interval", config.Daemon.CheckInterval)
	viper.SetDefault("daemon.auto_start", config.Daemon.AutoStart)
	viper.SetDefault("daemon.log_level", config.Daemon.LogLevel)
}

// saveDefaultConfig creates a default config file
func saveDefaultConfig(configDir string) error {
	configPath := filepath.Join(configDir, "config.yaml")

	// Don't overwrite existing config
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}

	configContent := `# Nagging Nancy Configuration

# Data storage directory (leave empty for auto-detection)
data_dir: ""

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
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Save saves the current configuration to file
func (c *Config) Save() error {
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set values in viper
	viper.Set("data_dir", c.DataDir)
	viper.Set("default.priority", c.Default.Priority)
	viper.Set("default.advance_minutes", c.Default.AdvanceMinutes)
	viper.Set("notifications.enabled", c.Notifications.Enabled)
	viper.Set("notifications.sound", c.Notifications.Sound)
	viper.Set("notifications.advance_minutes", c.Notifications.AdvanceMinutes)
	viper.Set("notifications.quiet_hours", c.Notifications.QuietHours)
	viper.Set("appearance.theme", c.Appearance.Theme)
	viper.Set("appearance.show_completed", c.Appearance.ShowCompleted)
	viper.Set("appearance.compact_mode", c.Appearance.CompactMode)
	viper.Set("appearance.show_icons", c.Appearance.ShowIcons)
	viper.Set("workhours.enabled", c.WorkHours.Enabled)
	viper.Set("workhours.start", c.WorkHours.Start)
	viper.Set("workhours.end", c.WorkHours.End)
	viper.Set("workhours.quiet_outside", c.WorkHours.QuietOutside)
	viper.Set("workhours.timezone", c.WorkHours.Timezone)
	viper.Set("daemon.check_interval", c.Daemon.CheckInterval)
	viper.Set("daemon.auto_start", c.Daemon.AutoStart)
	viper.Set("daemon.log_level", c.Daemon.LogLevel)

	// Write to file
	configPath := filepath.Join(configDir, "config.yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	// Validate priority
	if c.Default.Priority != "low" && c.Default.Priority != "medium" && c.Default.Priority != "high" {
		return fmt.Errorf("invalid default priority: %s", c.Default.Priority)
	}

	// Validate advance minutes (reasonable range)
	if c.Default.AdvanceMinutes < 0 || c.Default.AdvanceMinutes > 1440 {
		return fmt.Errorf("invalid default advance minutes: %d", c.Default.AdvanceMinutes)
	}

	if c.Notifications.AdvanceMinutes < 0 || c.Notifications.AdvanceMinutes > 1440 {
		return fmt.Errorf("invalid notification advance minutes: %d", c.Notifications.AdvanceMinutes)
	}

	// Validate theme
	if c.Appearance.Theme != "light" && c.Appearance.Theme != "dark" && c.Appearance.Theme != "auto" {
		return fmt.Errorf("invalid theme: %s", c.Appearance.Theme)
	}

	// Validate working hours
	if c.WorkHours.Enabled {
		if err := c.validateTimeFormat(c.WorkHours.Start); err != nil {
			return fmt.Errorf("invalid work start time: %w", err)
		}
		if err := c.validateTimeFormat(c.WorkHours.End); err != nil {
			return fmt.Errorf("invalid work end time: %w", err)
		}
	}

	// Validate daemon settings
	if c.Daemon.CheckInterval < 1 || c.Daemon.CheckInterval > 60 {
		return fmt.Errorf("invalid daemon check interval: %d (must be 1-60 minutes)", c.Daemon.CheckInterval)
	}

	logLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !logLevels[c.Daemon.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.Daemon.LogLevel)
	}

	return nil
}

// validateTimeFormat validates time format (HH:MM)
func (c *Config) validateTimeFormat(timeStr string) error {
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		return fmt.Errorf("invalid time format '%s', expected HH:MM", timeStr)
	}
	return nil
}

// IsWorkingHours checks if the given time falls within working hours
func (c *Config) IsWorkingHours(t time.Time) bool {
	if !c.WorkHours.Enabled {
		return true // If working hours not enabled, always return true
	}

	// Parse work hours
	start, err := time.Parse("15:04", c.WorkHours.Start)
	if err != nil {
		return true // If invalid format, assume working hours
	}

	end, err := time.Parse("15:04", c.WorkHours.End)
	if err != nil {
		return true // If invalid format, assume working hours
	}

	// Get current time in same format
	currentTime, err := time.Parse("15:04", t.Format("15:04"))
	if err != nil {
		return true // If can't parse, assume working hours
	}

	// Handle case where end time is before start time (overnight)
	if end.Before(start) {
		return currentTime.After(start) || currentTime.Before(end)
	}

	return currentTime.After(start) && currentTime.Before(end)
}

// ShouldNotify determines if notifications should be sent at the given time
func (c *Config) ShouldNotify(t time.Time) bool {
	if !c.Notifications.Enabled {
		return false
	}

	if !c.Notifications.QuietHours {
		return true
	}

	if c.WorkHours.QuietOutside {
		return c.IsWorkingHours(t)
	}

	return true
}

// GetConfigDir returns the configuration directory path
func (c *Config) GetConfigDir() string {
	return getConfigDir()
}

// GetDataDir returns the data directory path
func (c *Config) GetDataDir() string {
	if c.DataDir != "" {
		return c.DataDir
	}
	return getDataDir()
}

// Set sets a configuration value by key
func (c *Config) Set(key, value string) error {
	switch key {
	case "default.priority":
		if value != "low" && value != "medium" && value != "high" {
			return fmt.Errorf("invalid priority: %s", value)
		}
		c.Default.Priority = value
	case "appearance.theme":
		if value != "light" && value != "dark" && value != "auto" {
			return fmt.Errorf("invalid theme: %s", value)
		}
		c.Appearance.Theme = value
	case "workhours.start":
		if err := c.validateTimeFormat(value); err != nil {
			return err
		}
		c.WorkHours.Start = value
	case "workhours.end":
		if err := c.validateTimeFormat(value); err != nil {
			return err
		}
		c.WorkHours.End = value
	case "notifications.enabled":
		c.Notifications.Enabled = value == "true"
	case "notifications.sound":
		c.Notifications.Sound = value == "true"
	case "appearance.show_completed":
		c.Appearance.ShowCompleted = value == "true"
	case "appearance.compact_mode":
		c.Appearance.CompactMode = value == "true"
	case "appearance.show_icons":
		c.Appearance.ShowIcons = value == "true"
	case "workhours.enabled":
		c.WorkHours.Enabled = value == "true"
	case "workhours.quiet_outside":
		c.WorkHours.QuietOutside = value == "true"
	case "daemon.auto_start":
		c.Daemon.AutoStart = value == "true"
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return c.Save()
}

// Get gets a configuration value by key
func (c *Config) Get(key string) (string, error) {
	switch key {
	case "default.priority":
		return c.Default.Priority, nil
	case "appearance.theme":
		return c.Appearance.Theme, nil
	case "workhours.start":
		return c.WorkHours.Start, nil
	case "workhours.end":
		return c.WorkHours.End, nil
	case "notifications.enabled":
		if c.Notifications.Enabled {
			return "true", nil
		}
		return "false", nil
	case "notifications.sound":
		if c.Notifications.Sound {
			return "true", nil
		}
		return "false", nil
	case "appearance.show_completed":
		if c.Appearance.ShowCompleted {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("unknown configuration key: %s", key)
	}
}
