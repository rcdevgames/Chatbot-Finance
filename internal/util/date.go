package util

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func ParseRelativeDate(input string) (string, error) {
	now := time.Now()
	input = strings.ToLower(strings.TrimSpace(input))

	// Today patterns
	todayPatterns := []string{
		`^hari ini$`,
		`^today$`,
		`^skrng$`,
		`^sekarang$`,
	}

	for _, pattern := range todayPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			return now.Format("2006-01-02"), nil
		}
	}

	// Yesterday patterns
	yesterdayPatterns := []string{
		`^kemarin$`,
		`^yesterday$`,
	}

	for _, pattern := range yesterdayPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			yesterday := now.AddDate(0, 0, -1)
			return yesterday.Format("2006-01-02"), nil
		}
	}

	// This week patterns
	weekPatterns := []string{
		`^minggu ini$`,
		`^this week$`,
		`^week ini$`,
	}

	for _, pattern := range weekPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			startDate, _ := getWeekRange(now)
			return startDate, nil
		}
	}

	// Last week patterns
	lastWeekPatterns := []string{
		`^minggu lalu$`,
		`^last week$`,
	}

	for _, pattern := range lastWeekPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			lastWeek := now.AddDate(0, 0, -7)
			startDate, _ := getWeekRange(lastWeek)
			return startDate, nil
		}
	}

	// This month patterns
	monthPatterns := []string{
		`^bulan ini$`,
		`^this month$`,
	}

	for _, pattern := range monthPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			startDate, _ := getMonthRange(now)
			return startDate, nil
		}
	}

	// Last month patterns
	lastMonthPatterns := []string{
		`^bulan lalu$`,
		`^last month$`,
	}

	for _, pattern := range lastMonthPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			lastMonth := now.AddDate(0, -1, 0)
			startDate, _ := getMonthRange(lastMonth)
			return startDate, nil
		}
	}

	// Try to parse specific date formats
	dateFormats := []string{
		"2006-01-02",
		"02-01-2006",
		"02/01/2006",
		"2 Jan 2006",
		"2 Januari 2006",
		"Jan 2, 2006",
	}

	for _, format := range dateFormats {
		if parsed, err := time.Parse(format, input); err == nil {
			return parsed.Format("2006-01-02"), nil
		}
	}

	// If no pattern matches, return today as default
	return now.Format("2006-01-02"), nil
}

func getWeekRange(t time.Time) (string, string) {
	// Get Monday of the week
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}

	startOfWeek := t.AddDate(0, 0, -(weekday - 1))
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	return startOfWeek.Format("2006-01-02"), endOfWeek.Format("2006-01-02")
}

func getMonthRange(t time.Time) (string, string) {
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	return startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02")
}

func GetDateRangeFromPeriod(period string) (string, string) {
	now := time.Now()

	switch strings.ToLower(period) {
	case "today":
		today := now.Format("2006-01-02")
		return today, today
	case "this_week":
		return getWeekRange(now)
	case "last_week":
		lastWeek := now.AddDate(0, 0, -7)
		return getWeekRange(lastWeek)
	case "this_month":
		return getMonthRange(now)
	case "last_month":
		lastMonth := now.AddDate(0, -1, 0)
		return getMonthRange(lastMonth)
	default:
		// Default to this month
		return getMonthRange(now)
	}
}

func IsValidDate(dateString string) bool {
	_, err := time.Parse("2006-01-02", dateString)
	return err == nil
}

func FormatDate(dateString string, format string) string {
	if date, err := time.Parse("2006-01-02", dateString); err == nil {
		switch format {
		case "short":
			return date.Format("02 Jan")
		case "long":
			return date.Format("2 January 2006")
		case "relative":
			return getRelativeDateString(date)
		default:
			return dateString
		}
	}
	return dateString
}

func getRelativeDateString(date time.Time) string {
	now := time.Now()

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	givenDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, now.Location())

	daysDiff := int(givenDate.Sub(today).Hours() / 24)

	switch {
	case daysDiff == 0:
		return "Hari ini"
	case daysDiff == -1:
		return "Kemarin"
	case daysDiff == 1:
		return "Besok"
	case daysDiff > -7 && daysDiff < 7:
		if daysDiff > 0 {
			return fmt.Sprintf("%d hari lagi", daysDiff)
		}
		return fmt.Sprintf("%d hari lalu", -daysDiff)
	default:
		return date.Format("2 Jan 2006")
	}
}