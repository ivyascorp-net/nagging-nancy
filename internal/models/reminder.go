package models

import (
	"time"

	"github.com/google/uuid"
)

// Priority represents reminder priority levels
type Priority int

const (
	Low Priority = iota
	Medium
	High
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	default:
		return "medium"
	}
}

// ParsePriority converts a string to Priority
func ParsePriority(s string) Priority {
	switch s {
	case "low":
		return Low
	case "high":
		return High
	default:
		return Medium
	}
}

// Color returns the color associated with the priority
func (p Priority) Color() string {
	switch p {
	case Low:
		return "#10B981" // Green
	case Medium:
		return "#F59E0B" // Yellow/Amber
	case High:
		return "#EF4444" // Red
	default:
		return "#6B7280" // Gray
	}
}

// Icon returns the emoji/symbol for the priority
func (p Priority) Icon() string {
	switch p {
	case Low:
		return "🟢"
	case Medium:
		return "🟡"
	case High:
		return "🔴"
	default:
		return "⚪"
	}
}

// Reminder represents a single reminder
type Reminder struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	DueTime     time.Time      `json:"due_time"`
	Priority    Priority       `json:"priority"`
	Completed   bool           `json:"completed"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Tags        []string       `json:"tags,omitempty"`
	Recurring   *RecurringRule `json:"recurring,omitempty"`
}

// RecurringRule defines how often a reminder repeats
type RecurringRule struct {
	Frequency string     `json:"frequency"` // daily, weekly, monthly
	Interval  int        `json:"interval"`  // every N days/weeks/months
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// NewReminder creates a new reminder with generated ID and timestamps
func NewReminder(title string, dueTime time.Time, priority Priority) *Reminder {
	now := time.Now()
	return &Reminder{
		ID:        uuid.New().String(),
		Title:     title,
		DueTime:   dueTime,
		Priority:  priority,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      make([]string, 0),
	}
}

// IsOverdue checks if the reminder is past due
func (r *Reminder) IsOverdue() bool {
	if r.Completed {
		return false
	}
	return time.Now().After(r.DueTime)
}

// IsDueToday checks if the reminder is due today
func (r *Reminder) IsDueToday() bool {
	if r.Completed {
		return false
	}
	today := time.Now()
	due := r.DueTime
	return today.Year() == due.Year() &&
		today.YearDay() == due.YearDay()
}

// IsDueSoon checks if the reminder is due within the next hour
func (r *Reminder) IsDueSoon() bool {
	if r.Completed {
		return false
	}
	return time.Until(r.DueTime) <= time.Hour && time.Until(r.DueTime) > 0
}

// TimeUntilDue returns the duration until the reminder is due
func (r *Reminder) TimeUntilDue() time.Duration {
	if r.Completed {
		return 0
	}
	return time.Until(r.DueTime)
}

// Complete marks the reminder as completed
func (r *Reminder) Complete() {
	if !r.Completed {
		now := time.Now()
		r.Completed = true
		r.CompletedAt = &now
		r.UpdatedAt = now
	}
}

// Uncomplete marks the reminder as not completed
func (r *Reminder) Uncomplete() {
	if r.Completed {
		r.Completed = false
		r.CompletedAt = nil
		r.UpdatedAt = time.Now()
	}
}

// Toggle toggles the completion status
func (r *Reminder) Toggle() {
	if r.Completed {
		r.Uncomplete()
	} else {
		r.Complete()
	}
}

// Update updates the reminder's title and due time
func (r *Reminder) Update(title string, dueTime time.Time, priority Priority) {
	r.Title = title
	r.DueTime = dueTime
	r.Priority = priority
	r.UpdatedAt = time.Now()
}

// SetDescription sets the reminder's description
func (r *Reminder) SetDescription(description string) {
	r.Description = description
	r.UpdatedAt = time.Now()
}

// AddTag adds a tag to the reminder
func (r *Reminder) AddTag(tag string) {
	// Check if tag already exists
	for _, t := range r.Tags {
		if t == tag {
			return
		}
	}
	r.Tags = append(r.Tags, tag)
	r.UpdatedAt = time.Now()
}

// RemoveTag removes a tag from the reminder
func (r *Reminder) RemoveTag(tag string) {
	for i, t := range r.Tags {
		if t == tag {
			r.Tags = append(r.Tags[:i], r.Tags[i+1:]...)
			r.UpdatedAt = time.Now()
			return
		}
	}
}

// HasTag checks if the reminder has a specific tag
func (r *Reminder) HasTag(tag string) bool {
	for _, t := range r.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Status returns a human-readable status string
func (r *Reminder) Status() string {
	if r.Completed {
		return "✓ Completed"
	}
	if r.IsOverdue() {
		return "⚠ Overdue"
	}
	if r.IsDueToday() {
		return "📅 Due Today"
	}
	if r.IsDueSoon() {
		return "⏰ Due Soon"
	}
	return "📝 Pending"
}

// FormattedDueTime returns a nicely formatted due time string
func (r *Reminder) FormattedDueTime() string {
	now := time.Now()
	due := r.DueTime

	// Same day
	if now.Year() == due.Year() && now.YearDay() == due.YearDay() {
		return "Today " + due.Format("3:04 PM")
	}

	// Tomorrow
	tomorrow := now.AddDate(0, 0, 1)
	if tomorrow.Year() == due.Year() && tomorrow.YearDay() == due.YearDay() {
		return "Tomorrow " + due.Format("3:04 PM")
	}

	// This week
	if due.Sub(now) < 7*24*time.Hour && due.Sub(now) > 0 {
		return due.Format("Monday 3:04 PM")
	}

	// This year
	if now.Year() == due.Year() {
		return due.Format("Jan 2 3:04 PM")
	}

	// Different year
	return due.Format("Jan 2, 2006 3:04 PM")
}
