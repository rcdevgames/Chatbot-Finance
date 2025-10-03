package util

import (
	"fmt"
	"math"
	"strings"
)

func FormatCurrency(amount float64) string {
	// Convert to string and format with thousand separators
	str := fmt.Sprintf("%.0f", amount)

	// Add thousand separators
	var result string
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += "."
		}
		result += string(digit)
	}

	return "Rp " + result
}

func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}

func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}

func CapitalizeFirst(text string) string {
	if len(text) == 0 {
		return text
	}
	return strings.ToUpper(text[:1]) + strings.ToLower(text[1:])
}

func ParseAmount(input string) (float64, error) {
	// Remove all non-digit and non-decimal characters
	cleaned := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == ',' {
			return r
		}
		return -1
	}, input)

	// Replace comma with dot for decimal
	cleaned = strings.ReplaceAll(cleaned, ",", ".")

	// Handle common Indonesian abbreviations
	lowerInput := strings.ToLower(input)
	multiplier := 1.0

	if strings.Contains(lowerInput, "jt") || strings.Contains(lowerInput, "juta") {
		multiplier = 1000000
	} else if strings.Contains(lowerInput, "rb") || strings.Contains(lowerInput, "ribu") {
		multiplier = 1000
	}

	// Parse to float
	var amount float64
	if cleaned != "" {
		_, err := fmt.Sscanf(cleaned, "%f", &amount)
		if err != nil {
			return 0, err
		}
	}

	return amount * multiplier, nil
}

func FormatNumberWithSuffix(num float64) string {
	if num >= 1000000 {
		return fmt.Sprintf("%.1fjt", num/1000000)
	} else if num >= 1000 {
		return fmt.Sprintf("%.0frb", num/1000)
	}
	return fmt.Sprintf("%.0f", num)
}

func EscapeMarkdown(text string) string {
	// Escape special markdown characters
	chars := []string{"_", "*", "`", "["}
	for _, char := range chars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

func GenerateProgressBar(current, max float64, length int) string {
	if max == 0 {
		return strings.Repeat("░", length)
	}

	percentage := current / max
	if percentage > 1 {
		percentage = 1
	}

	filled := int(math.Round(float64(length) * percentage))
	empty := length - filled

	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

func FormatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	if hours > 0 {
		return fmt.Sprintf("%dj %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}