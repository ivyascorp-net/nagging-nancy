package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Store handles data persistence for reminders
type Store struct {
	filePath  string
	reminders map[string]*Reminder
	mutex     sync.RWMutex
}

// FilterOptions defines options for filtering reminders
type FilterOptions struct {
	ShowCompleted bool
	Priority      *Priority
	DueToday      bool
	Overdue       bool
	Tags          []string
	Limit         int
}

// NewStore creates a new store instance
func NewStore(dataDir string) (*Store, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	filePath := filepath.Join(dataDir, "reminders.json")
	store := &Store{
		filePath:  filePath,
		reminders: make(map[string]*Reminder),
	}

	// Load existing data
	if err := store.Load(); err != nil {
		return nil, fmt.Errorf("failed to load reminders: %w", err)
	}

	return store, nil
}

// Load reads reminders from file
func (s *Store) Load() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// File doesn't exist yet, that's ok for a new installation
		return nil
	}

	// Read file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return nil
	}

	// Parse JSON
	var reminders []*Reminder
	if err := json.Unmarshal(data, &reminders); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert slice to map for efficient lookups
	s.reminders = make(map[string]*Reminder)
	for _, reminder := range reminders {
		if reminder != nil {
			s.reminders[reminder.ID] = reminder
		}
	}

	return nil
}

// Save writes reminders to file
func (s *Store) Save() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Convert map to slice for JSON serialization
	reminders := make([]*Reminder, 0, len(s.reminders))
	for _, reminder := range s.reminders {
		if reminder != nil {
			reminders = append(reminders, reminder)
		}
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(reminders, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file with proper permissions
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Add adds a new reminder to the store
func (s *Store) Add(reminder *Reminder) error {
	if reminder == nil {
		return fmt.Errorf("reminder cannot be nil")
	}

	s.mutex.Lock()
	s.reminders[reminder.ID] = reminder
	s.mutex.Unlock()

	return s.Save()
}

// Get retrieves a reminder by ID
func (s *Store) Get(id string) (*Reminder, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	reminder, exists := s.reminders[id]
	if !exists {
		return nil, fmt.Errorf("reminder with ID %s not found", id)
	}

	// Return a copy to prevent external modification
	reminderCopy := *reminder
	return &reminderCopy, nil
}

// Update updates an existing reminder
func (s *Store) Update(reminder *Reminder) error {
	if reminder == nil {
		return fmt.Errorf("reminder cannot be nil")
	}

	s.mutex.Lock()
	_, exists := s.reminders[reminder.ID]
	if !exists {
		s.mutex.Unlock()
		return fmt.Errorf("reminder with ID %s not found", reminder.ID)
	}

	reminder.UpdatedAt = time.Now()
	s.reminders[reminder.ID] = reminder
	s.mutex.Unlock()

	return s.Save()
}

// Delete removes a reminder from the store
func (s *Store) Delete(id string) error {
	s.mutex.Lock()
	_, exists := s.reminders[id]
	if !exists {
		s.mutex.Unlock()
		return fmt.Errorf("reminder with ID %s not found", id)
	}

	delete(s.reminders, id)
	s.mutex.Unlock()

	return s.Save()
}

// GetAll returns all reminders with optional filtering
func (s *Store) GetAll(filter *FilterOptions) []*Reminder {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	reminders := make([]*Reminder, 0, len(s.reminders))

	for _, reminder := range s.reminders {
		if reminder == nil {
			continue
		}

		// Apply filters
		if filter != nil {
			if !filter.ShowCompleted && reminder.Completed {
				continue
			}

			if filter.Priority != nil && reminder.Priority != *filter.Priority {
				continue
			}

			if filter.DueToday && !reminder.IsDueToday() {
				continue
			}

			if filter.Overdue && !reminder.IsOverdue() {
				continue
			}

			// Check tags filter
			if len(filter.Tags) > 0 {
				hasTag := false
				for _, filterTag := range filter.Tags {
					if reminder.HasTag(filterTag) {
						hasTag = true
						break
					}
				}
				if !hasTag {
					continue
				}
			}
		}

		// Create a copy to prevent external modification
		reminderCopy := *reminder
		reminders = append(reminders, &reminderCopy)
	}

	// Sort by due time (ascending)
	sort.Slice(reminders, func(i, j int) bool {
		// Completed items go to the bottom
		if reminders[i].Completed && !reminders[j].Completed {
			return false
		}
		if !reminders[i].Completed && reminders[j].Completed {
			return true
		}

		// Sort by due time
		return reminders[i].DueTime.Before(reminders[j].DueTime)
	})

	// Apply limit if specified
	if filter != nil && filter.Limit > 0 && len(reminders) > filter.Limit {
		reminders = reminders[:filter.Limit]
	}

	return reminders
}

// GetByPriority returns reminders filtered by priority
func (s *Store) GetByPriority(priority Priority) []*Reminder {
	filter := &FilterOptions{
		Priority: &priority,
	}
	return s.GetAll(filter)
}

// GetDueToday returns reminders due today
func (s *Store) GetDueToday() []*Reminder {
	filter := &FilterOptions{
		DueToday: true,
	}
	return s.GetAll(filter)
}

// GetOverdue returns overdue reminders
func (s *Store) GetOverdue() []*Reminder {
	filter := &FilterOptions{
		Overdue: true,
	}
	return s.GetAll(filter)
}

// GetCompleted returns completed reminders
func (s *Store) GetCompleted() []*Reminder {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	reminders := make([]*Reminder, 0)

	for _, reminder := range s.reminders {
		if reminder != nil && reminder.Completed {
			reminderCopy := *reminder
			reminders = append(reminders, &reminderCopy)
		}
	}

	// Sort by completion time (most recent first)
	sort.Slice(reminders, func(i, j int) bool {
		if reminders[i].CompletedAt == nil && reminders[j].CompletedAt == nil {
			return reminders[i].CreatedAt.After(reminders[j].CreatedAt)
		}
		if reminders[i].CompletedAt == nil {
			return false
		}
		if reminders[j].CompletedAt == nil {
			return true
		}
		return reminders[i].CompletedAt.After(*reminders[j].CompletedAt)
	})

	return reminders
}

// GetActive returns non-completed reminders
func (s *Store) GetActive() []*Reminder {
	filter := &FilterOptions{
		ShowCompleted: false,
	}
	return s.GetAll(filter)
}

// Count returns counts of different reminder categories
func (s *Store) Count() (total, active, completed, overdue int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, reminder := range s.reminders {
		if reminder == nil {
			continue
		}

		total++
		if reminder.Completed {
			completed++
		} else {
			active++
			if reminder.IsOverdue() {
				overdue++
			}
		}
	}

	return
}

// GetTags returns all unique tags used in reminders
func (s *Store) GetTags() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	tagSet := make(map[string]bool)
	for _, reminder := range s.reminders {
		if reminder == nil {
			continue
		}

		for _, tag := range reminder.Tags {
			tagSet[tag] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	sort.Strings(tags)
	return tags
}

// CompleteReminder marks a reminder as completed by ID
func (s *Store) CompleteReminder(id string) error {
	s.mutex.Lock()
	reminder, exists := s.reminders[id]
	if !exists {
		s.mutex.Unlock()
		return fmt.Errorf("reminder with ID %s not found", id)
	}

	reminder.Complete()
	s.mutex.Unlock()

	return s.Save()
}

// ToggleReminder toggles the completion status of a reminder by ID
func (s *Store) ToggleReminder(id string) error {
	s.mutex.Lock()
	reminder, exists := s.reminders[id]
	if !exists {
		s.mutex.Unlock()
		return fmt.Errorf("reminder with ID %s not found", id)
	}

	reminder.Toggle()
	s.mutex.Unlock()

	return s.Save()
}

// Cleanup removes old completed reminders (older than 30 days)
func (s *Store) Cleanup() error {
	s.mutex.Lock()
	cutoff := time.Now().AddDate(0, 0, -30) // 30 days ago
	deleted := 0

	for id, reminder := range s.reminders {
		if reminder != nil && reminder.Completed {
			completedAt := reminder.CompletedAt
			if completedAt != nil && completedAt.Before(cutoff) {
				delete(s.reminders, id)
				deleted++
			}
		}
	}
	s.mutex.Unlock()

	if deleted > 0 {
		return s.Save()
	}

	return nil
}

// Export exports all reminders to a JSON string
func (s *Store) Export() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	reminders := make([]*Reminder, 0, len(s.reminders))
	for _, reminder := range s.reminders {
		if reminder != nil {
			reminders = append(reminders, reminder)
		}
	}

	return json.MarshalIndent(reminders, "", "  ")
}

// Import imports reminders from JSON data (merges with existing)
func (s *Store) Import(data []byte) error {
	var importedReminders []*Reminder
	if err := json.Unmarshal(data, &importedReminders); err != nil {
		return fmt.Errorf("failed to parse import data: %w", err)
	}

	s.mutex.Lock()
	imported := 0
	for _, reminder := range importedReminders {
		if reminder != nil {
			// Check if reminder with same ID already exists
			if _, exists := s.reminders[reminder.ID]; !exists {
				s.reminders[reminder.ID] = reminder
				imported++
			}
		}
	}
	s.mutex.Unlock()

	if imported > 0 {
		return s.Save()
	}

	return nil
}
