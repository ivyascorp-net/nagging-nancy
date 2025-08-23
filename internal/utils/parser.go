package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ivyascorp-net/nagging-nancy/internal/models"
)

// ParsedReminder represents the result of parsing reminder text
type ParsedReminder struct {
	Title    string
	DueTime  time.Time
	Priority models.Priority
	Tags     []string
	HasTime  bool
}

// TimePattern represents a regex pattern for parsing time expressions
type TimePattern struct {
	Pattern *regexp.Regexp
	Handler func(matches []string, baseTime time.Time) (time.Time, error)
}

// Common time patterns for natural language parsing
var timePatterns = []TimePattern{
	// "today at 3pm", "today at 15:30"
	{
		regexp.MustCompile(`(?i)today\s+at\s+(\d{1,2}):?(\d{0,2})\s*(am|pm)?`),
		parseTimeToday,
	},
	// "tomorrow at 2:30pm"
	{
		regexp.MustCompile(`(?i)tomorrow\s+at\s+(\d{1,2}):?(\d{0,2})\s*(am|pm)?`),
		parseTimeTomorrow,
	},
	// "in 30 minutes", "in 2 hours"
	{
		regexp.MustCompile(`(?i)in\s+(\d+)\s+(minute|minutes|hour|hours|min|hr|hrs)s?`),
		parseTimeRelative,
	},
	// "at 3pm", "at 15:30"
	{
		regexp.MustCompile(`(?i)at\s+(\d{1,2}):?(\d{0,2})\s*(am|pm)?`),
		parseTimeToday,
	},
	// "3pm", "15:30"
	{
		regexp.MustCompile(`(?i)^(\d{1,2}):?(\d{0,2})\s*(am|pm)?$`),
		parseTimeToday,
	},
	// "monday at 2pm", "friday 3:30pm"
	{
		regexp.MustCompile(`(?i)(monday|tuesday|wednesday|thursday|friday|saturday|sunday)\s+(?:at\s+)?(\d{1,2}):?(\d{0,2})\s*(am|pm)?`),
		parseTimeWeekday,
	},
}

// Priority patterns for detecting priority in text
var priorityPatterns = []struct {
	pattern  *regexp.Regexp
	priority models.Priority
}{
	{regexp.MustCompile(`(?i)\b(urgent|important|high|critical|asap)\b`), models.High},
	{regexp.MustCompile(`(?i)\b(low|minor|sometime|eventually)\b`), models.Low},
}

// ParseReminder parses a reminder string and extracts structured information
func ParseReminder(text string, defaultPriority models.Priority) (*ParsedReminder, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("reminder text cannot be empty")
	}

	result := &ParsedReminder{
		Title:    text,
		DueTime:  time.Now().Add(time.Hour), // Default to 1 hour from now
		Priority: defaultPriority,
		Tags:     make([]string, 0),
		HasTime:  false,
	}

	// Extract time information
	if dueTime, cleanText, hasTime := extractTime(text); hasTime {
		result.DueTime = dueTime
		result.Title = strings.TrimSpace(cleanText)
		result.HasTime = true
	}

	// Extract priority information
	if priority, cleanText := extractPriority(result.Title); priority != defaultPriority {
		result.Priority = priority
		result.Title = strings.TrimSpace(cleanText)
	}

	// Extract tags (#hashtag format)
	if tags, cleanText := extractTags(result.Title); len(tags) > 0 {
		result.Tags = tags
		result.Title = strings.TrimSpace(cleanText)
	}

	// Clean up the title
	result.Title = strings.TrimSpace(result.Title)
	if result.Title == "" {
		return nil, fmt.Errorf("reminder title cannot be empty after parsing")
	}

	return result, nil
}

// extractTime tries to extract time information from text
func extractTime(text string) (time.Time, string, bool) {
	baseTime := time.Now()

	for _, pattern := range timePatterns {
		if matches := pattern.Pattern.FindStringSubmatch(text); matches != nil {
			if parsedTime, err := pattern.Handler(matches, baseTime); err == nil {
				// Remove the matched time expression from text
				cleanText := pattern.Pattern.ReplaceAllString(text, "")
				cleanText = strings.TrimSpace(cleanText)
				return parsedTime, cleanText, true
			}
		}
	}

	return baseTime.Add(time.Hour), text, false
}

// parseTimeToday parses time expressions for today
func parseTimeToday(matches []string, baseTime time.Time) (time.Time, error) {
	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return baseTime, err
	}

	var minute int
	if len(matches) > 2 && matches[2] != "" {
		minute, err = strconv.Atoi(matches[2])
		if err != nil {
			return baseTime, err
		}
	}

	// Handle AM/PM
	if len(matches) > 3 && matches[3] != "" {
		ampm := strings.ToLower(matches[3])
		if ampm == "pm" && hour < 12 {
			hour += 12
		} else if ampm == "am" && hour == 12 {
			hour = 0
		}
	}

	// Validate hour and minute
	if hour > 23 || minute > 59 {
		return baseTime, fmt.Errorf("invalid time: %d:%d", hour, minute)
	}

	now := baseTime
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	// If the time has already passed today, schedule for tomorrow
	if targetTime.Before(now) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}

	return targetTime, nil
}

// parseTimeTomorrow parses time expressions for tomorrow
func parseTimeTomorrow(matches []string, baseTime time.Time) (time.Time, error) {
	targetTime, err := parseTimeToday(matches, baseTime)
	if err != nil {
		return baseTime, err
	}

	// Add one day to make it tomorrow
	return targetTime.AddDate(0, 0, 1), nil
}

// parseTimeRelative parses relative time expressions like "in 30 minutes"
func parseTimeRelative(matches []string, baseTime time.Time) (time.Time, error) {
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return baseTime, err
	}

	unit := strings.ToLower(matches[2])
	var duration time.Duration

	switch unit {
	case "minute", "minutes", "min":
		duration = time.Duration(amount) * time.Minute
	case "hour", "hours", "hr", "hrs":
		duration = time.Duration(amount) * time.Hour
	default:
		return baseTime, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return baseTime.Add(duration), nil
}

// parseTimeWeekday parses weekday time expressions
func parseTimeWeekday(matches []string, baseTime time.Time) (time.Time, error) {
	weekdayStr := strings.ToLower(matches[1])

	// Map weekday names to time.Weekday
	weekdays := map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}

	targetWeekday, exists := weekdays[weekdayStr]
	if !exists {
		return baseTime, fmt.Errorf("invalid weekday: %s", weekdayStr)
	}

	// Parse the time part
	hour, err := strconv.Atoi(matches[2])
	if err != nil {
		return baseTime, err
	}

	var minute int
	if len(matches) > 3 && matches[3] != "" {
		minute, err = strconv.Atoi(matches[3])
		if err != nil {
			return baseTime, err
		}
	}

	// Handle AM/PM
	if len(matches) > 4 && matches[4] != "" {
		ampm := strings.ToLower(matches[4])
		if ampm == "pm" && hour < 12 {
			hour += 12
		} else if ampm == "am" && hour == 12 {
			hour = 0
		}
	}

	// Calculate target date
	now := baseTime
	currentWeekday := now.Weekday()

	daysUntilTarget := int(targetWeekday - currentWeekday)
	if daysUntilTarget <= 0 {
		daysUntilTarget += 7 // Next week
	}

	targetDate := now.AddDate(0, 0, daysUntilTarget)
	targetTime := time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		hour, minute, 0, 0, now.Location(),
	)

	return targetTime, nil
}

// extractPriority extracts priority keywords from text
func extractPriority(text string) (models.Priority, string) {
	for _, pattern := range priorityPatterns {
		if pattern.pattern.MatchString(text) {
			cleanText := pattern.pattern.ReplaceAllString(text, "")
			cleanText = strings.TrimSpace(cleanText)
			return pattern.priority, cleanText
		}
	}
	return models.Medium, text
}

// extractTags extracts hashtags from text
func extractTags(text string) ([]string, string) {
	tagPattern := regexp.MustCompile(`#(\w+)`)
	matches := tagPattern.FindAllStringSubmatch(text, -1)

	if len(matches) == 0 {
		return nil, text
	}

	tags := make([]string, 0, len(matches))
	for _, match := range matches {
		tags = append(tags, match[1])
	}

	// Remove hashtags from text
	cleanText := tagPattern.ReplaceAllString(text, "")
	cleanText = strings.TrimSpace(cleanText)

	return tags, cleanText
}

// ParseTimeString parses various time string formats
func ParseTimeString(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	now := time.Now()

	// Try different time formats
	formats := []string{
		"15:04",             // 14:30
		"3:04 PM",           // 2:30 PM
		"3:04PM",            // 2:30PM
		"3PM",               // 2PM
		"15",                // 14 (hour only)
		"2006-01-02 15:04",  // Full datetime
		"Jan 2 15:04",       // Month day time
		"Jan 2, 2006 15:04", // Month day year time
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			// If only time is provided, use today's date
			if format == "15:04" || format == "3:04 PM" || format == "3:04PM" || format == "3PM" || format == "15" {
				return time.Date(now.Year(), now.Month(), now.Day(),
					t.Hour(), t.Minute(), 0, 0, now.Location()), nil
			}
			return t, nil
		}
	}

	return now, fmt.Errorf("unable to parse time: %s", timeStr)
}

// FormatDuration returns a human-readable duration string
func FormatDuration(d time.Duration) string {
	if d < 0 {
		return "overdue"
	}

	if d < time.Minute {
		return "now"
	}

	if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	}

	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if hours == 1 && minutes == 0 {
			return "1 hour"
		}
		if minutes == 0 {
			return fmt.Sprintf("%d hours", hours)
		}
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}

	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

// ParsePriorityString converts a string to Priority
func ParsePriorityString(priorityStr string) models.Priority {
	return models.ParsePriority(strings.ToLower(strings.TrimSpace(priorityStr)))
}

// ValidateReminderInput validates reminder input
func ValidateReminderInput(title string, dueTime time.Time) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("reminder title cannot be empty")
	}

	// Don't allow reminders too far in the past (more than 1 hour)
	if time.Since(dueTime) > time.Hour {
		return fmt.Errorf("due time cannot be more than 1 hour in the past")
	}

	// Don't allow reminders too far in the future (more than 10 years)
	if time.Until(dueTime) > 10*365*24*time.Hour {
		return fmt.Errorf("due time cannot be more than 10 years in the future")
	}

	return nil
}
