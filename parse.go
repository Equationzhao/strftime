package strftime

import (
	"fmt"
	"strconv"
	"time"
)

// parseResult stores the fields obtained during parsing
type parseResult struct {
	year    int
	month   int
	day     int
	hour    int
	minute  int
	second  int
	hour12  bool // Whether to use 12-hour format (%I)
	ampmSet bool // Whether %p (AM/PM marker) appeared
	isPM    bool // Whether it is PM when using 12-hour format
}

// parseFixedInt reads a fixed-length numeric string from s[pos:] and returns the corresponding integer and new position
func parseFixedInt(s string, pos, length int) (int, int, error) {
	if pos+length > len(s) {
		return 0, pos, fmt.Errorf("expected %d digits at position %d, but reached end of string", length, pos)
	}
	substr := s[pos : pos+length]
	val, err := strconv.Atoi(substr)
	if err != nil {
		return 0, pos, fmt.Errorf("failed to parse integer '%s' at position %d: %v", substr, pos, err)
	}
	return val, pos + length, nil
}

// parseIntVariable reads at least minDigits and at most maxDigits digits from s[pos:]
func parseIntVariable(s string, pos, minDigits, maxDigits int) (int, int, error) {
	start := pos
	count := 0
	for pos < len(s) && count < maxDigits && s[pos] >= '0' && s[pos] <= '9' {
		pos++
		count++
	}
	if count < minDigits {
		return 0, start, fmt.Errorf("expected at least %d digits at position %d", minDigits, start)
	}
	val, err := strconv.Atoi(s[start:pos])
	if err != nil {
		return 0, start, err
	}
	return val, pos, nil
}

// ParseL parses the input string s according to the specified format and locale, and returns a time.Time object.
// Supported conversion specifiers include:
//
//	%Y,%y,%m,%d,%e,%H,%I,%M,%S,%p,%D,%F,%B,%b,%h,%A,%a, and %%.
//
// For POSIX extensions (e.g., starting with %E or %O), the extension prefix is skipped, and formats like "%EY" and "%E%Y" are supported.
func ParseL(format, s string, locale *Locale) (time.Time, error) {
	if locale == nil {
		locale = DefaultLocale
	}

	// Use the current time as the default value, parts not parsed will use the corresponding parts of the current time
	base := time.Now()
	result := parseResult{
		year:    base.Year(),
		month:   int(base.Month()),
		day:     base.Day(),
		hour:    base.Hour(),
		minute:  base.Minute(),
		second:  base.Second(),
		hour12:  false,
		ampmSet: false,
		isPM:    false,
	}

	i, j := 0, 0
	// Traverse the format string
	for i < len(format) {
		if format[i] == '%' {
			i++ // Skip '%'
			if i >= len(format) {
				return time.Time{}, fmt.Errorf("incomplete format specifier at end")
			}
			// Check for POSIX extension prefix %E or %O
			if format[i] == 'E' || format[i] == 'O' {
				i++ // Skip extension marker
				// If followed by a '%', skip it as well (support "%E%Y" format)
				if i < len(format) && format[i] == '%' {
					i++
				}
				if i >= len(format) {
					return time.Time{}, fmt.Errorf("incomplete format specifier after posix extension")
				}
			}
			// Get the conversion specifier character and increment the pointer
			spec := format[i]
			i++
			switch spec {
			case 'Y': // 4-digit year
				result.year, j, _ = parseFixedInt(s, j, 4)
			case 'y': // 2-digit year, converted to 1900s or 2000s by convention
				var twoDigit int
				twoDigit, j, _ = parseFixedInt(s, j, 2)
				if twoDigit < 69 {
					result.year = 2000 + twoDigit
				} else {
					result.year = 1900 + twoDigit
				}
			case 'm': // Month (two digits)
				result.month, j, _ = parseFixedInt(s, j, 2)
			case 'd': // Day (two digits)
				result.day, j, _ = parseFixedInt(s, j, 2)
			case 'e': // Day (1-2 digits, leading space may exist)
				if j < len(s) && s[j] == ' ' {
					j++
				}
				result.day, j, _ = parseIntVariable(s, j, 1, 2)
			case 'H': // 24-hour format hour
				result.hour, j, _ = parseFixedInt(s, j, 2)
			case 'I': // 12-hour format hour
				result.hour, j, _ = parseFixedInt(s, j, 2)
				result.hour12 = true
			case 'M': // Minute
				result.minute, j, _ = parseFixedInt(s, j, 2)
			case 'S': // Second
				result.second, j, _ = parseFixedInt(s, j, 2)
			case 'p': // AM/PM marker
				if len(s[j:]) >= len(locale.AM) && s[j:j+len(locale.AM)] == locale.AM {
					result.ampmSet = true
					result.isPM = false
					j += len(locale.AM)
				} else if len(s[j:]) >= len(locale.PM) && s[j:j+len(locale.PM)] == locale.PM {
					result.ampmSet = true
					result.isPM = true
					j += len(locale.PM)
				} else {
					return time.Time{}, fmt.Errorf("expected AM/PM marker at position %d", j)
				}
			case 'D':
				// "%D" equals "%m/%d/%y"
				result.month, j, _ = parseFixedInt(s, j, 2)
				if j >= len(s) || s[j] != '/' {
					return time.Time{}, fmt.Errorf("expected '/' after month in %%D")
				}
				j++
				result.day, j, _ = parseFixedInt(s, j, 2)
				if j >= len(s) || s[j] != '/' {
					return time.Time{}, fmt.Errorf("expected '/' after day in %%D")
				}
				j++
				var twoDigit int
				twoDigit, j, _ = parseFixedInt(s, j, 2)
				if twoDigit < 69 {
					result.year = 2000 + twoDigit
				} else {
					result.year = 1900 + twoDigit
				}
			case 'F': // Equivalent to "%Y-%m-%d"
				result.year, j, _ = parseFixedInt(s, j, 4)
				if j >= len(s) || s[j] != '-' {
					return time.Time{}, fmt.Errorf("expected '-' after year in %%F")
				}
				j++
				result.month, j, _ = parseFixedInt(s, j, 2)
				if j >= len(s) || s[j] != '-' {
					return time.Time{}, fmt.Errorf("expected '-' after month in %%F")
				}
				j++
				result.day, j, _ = parseFixedInt(s, j, 2)
			case 'B': // Full month name (based on locale.MonthsFull)
				found := false
				for iMonth, mName := range locale.MonthsFull {
					if len(s[j:]) >= len(mName) && s[j:j+len(mName)] == mName {
						result.month = iMonth + 1
						j += len(mName)
						found = true
						break
					}
				}
				if !found {
					return time.Time{}, fmt.Errorf("failed to parse full month name at position %d", j)
				}
			case 'b', 'h': // Abbreviated month name (based on locale.MonthsAbbrev)
				found := false
				for iMonth, mName := range locale.MonthsAbbrev {
					if len(s[j:]) >= len(mName) && s[j:j+len(mName)] == mName {
						result.month = iMonth + 1
						j += len(mName)
						found = true
						break
					}
				}
				if !found {
					return time.Time{}, fmt.Errorf("failed to parse abbreviated month name at position %d", j)
				}
			case 'A': // Full weekday name (consumed but does not affect values)
				found := false
				for _, wName := range locale.WeekdaysFull {
					if len(s[j:]) >= len(wName) && s[j:j+len(wName)] == wName {
						j += len(wName)
						found = true
						break
					}
				}
				if !found {
					return time.Time{}, fmt.Errorf("failed to parse full weekday name at position %d", j)
				}
			case 'a': // Abbreviated weekday name (consumed but does not affect values)
				found := false
				for _, wName := range locale.WeekdaysAbbrev {
					if len(s[j:]) >= len(wName) && s[j:j+len(wName)] == wName {
						j += len(wName)
						found = true
						break
					}
				}
				if !found {
					return time.Time{}, fmt.Errorf("failed to parse abbreviated weekday name at position %d", j)
				}
			case '%': // Literal '%'
				if j >= len(s) || s[j] != '%' {
					return time.Time{}, fmt.Errorf("expected literal '%%' at position %d", j)
				}
				j++
			default:
				// For unknown conversion specifiers, output '%' and the character as is
				return time.Time{}, fmt.Errorf("unsupported conversion specifier: %%%c", spec)
			}
		} else {
			// Non-conversion specifier part, requires literal match
			if j >= len(s) || s[j] != format[i] {
				return time.Time{}, fmt.Errorf("literal mismatch at position %d: expected '%c', got '%c'", j, format[i], s[j])
			}
			i++
			j++
		}
	}

	// Skip trailing whitespace characters in the input string
	for j < len(s) && (s[j] == ' ' || s[j] == '\t') {
		j++
	}
	if j != len(s) {
		return time.Time{}, fmt.Errorf("unparsed trailing characters at position %d", j)
	}

	// For 12-hour format, %p must be used
	if result.hour12 && !result.ampmSet {
		return time.Time{}, fmt.Errorf("12-hour format specified but missing AM/PM marker")
	}

	// Adjust based on 12-hour format and AM/PM
	if result.hour12 {
		if result.hour < 1 || result.hour > 12 {
			return time.Time{}, fmt.Errorf("invalid hour %d for 12-hour format", result.hour)
		}
		if result.isPM && result.hour != 12 {
			result.hour += 12
		} else if !result.isPM && result.hour == 12 {
			result.hour = 0
		}
	}

	parsedTime := time.Date(result.year, time.Month(result.month), result.day, result.hour, result.minute, result.second, 0, base.Location())
	return parsedTime, nil
}

// Parse parses the string using the default locale
func Parse(format, s string) (time.Time, error) {
	return ParseL(format, s, DefaultLocale)
}
