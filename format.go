package strftime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Strftime formats time according to the specified format string using the default locale
func Strftime(format string, t time.Time) string {
	return StrftimeL(format, t, DefaultLocale)
}

// StrftimeL formats time according to the specified format string and locale
func StrftimeL(format string, t time.Time, loc *Locale) string {
	if loc == nil {
		loc = DefaultLocale
	}

	var result strings.Builder
	i := 0
	for i < len(format) {
		if format[i] != '%' {
			result.WriteByte(format[i])
			i++
			continue
		}

		// Handle % directives
		i++
		if i >= len(format) {
			break
		}

		// Handle %% escape sequence to produce a single %
		if format[i] == '%' {
			result.WriteByte('%')
			i++
			continue
		}

		// Handle format modifiers for GNU libc extension
		padChar := byte('0')
		padWidth := 2

		switch format[i] {
		case '-': // No padding
			padWidth = 0
			i++
		case '_': // Space padding
			padChar = ' '
			i++
		case '0': // Zero padding (default)
			i++
		}

		if i >= len(format) {
			break
		}

		// Handle POSIX locale extensions
		if format[i] == 'E' || format[i] == 'O' {
			i++
			if i >= len(format) {
				break
			}
		}

		switch format[i] {
		case 'A': // Full weekday name
			result.WriteString(loc.WeekdaysFull[t.Weekday()])
		case 'a': // Abbreviated weekday name
			result.WriteString(loc.WeekdaysAbbrev[t.Weekday()])
		case 'B': // Full month name
			result.WriteString(loc.MonthsFull[t.Month()-1])
		case 'b', 'h': // Abbreviated month name
			result.WriteString(loc.MonthsAbbrev[t.Month()-1])
		case 'C': // Century
			century := t.Year() / 100
			if padWidth > 0 {
				result.WriteString(formatInt(century, padWidth, padChar))
			} else {
				result.WriteString(strconv.Itoa(century))
			}
		case 'c': // Date and time representation
			result.WriteString(t.Format("Mon Jan 2 15:04:05 2006"))
		case 'D': // %m/%d/%y
			result.WriteString(t.Format("01/02/06"))
		case 'd': // Day of month (01-31)
			if padWidth > 0 {
				result.WriteString(formatInt(t.Day(), padWidth, padChar))
			} else {
				result.WriteString(strconv.Itoa(t.Day()))
			}
		case 'e': // Day of month (space-padded)
			result.WriteString(fmt.Sprintf("%2d", t.Day()))
		case 'F': // ISO 8601 date
			result.WriteString(t.Format("2006-01-02"))
		case 'G': // ISO 8601 year
			year, _ := t.ISOWeek()
			if padWidth > 0 {
				result.WriteString(formatInt(year, 4, padChar))
			} else {
				result.WriteString(strconv.Itoa(year))
			}
		case 'g': // ISO 8601 year (2 digits)
			year, _ := t.ISOWeek()
			if padWidth > 0 {
				result.WriteString(formatInt(year%100, 2, padChar))
			} else {
				result.WriteString(strconv.Itoa(year % 100))
			}
		case 'H': // Hour in 24h format (00-23)
			if padWidth > 0 {
				result.WriteString(formatInt(t.Hour(), padWidth, padChar))
			} else {
				result.WriteString(strconv.Itoa(t.Hour()))
			}
		case 'I': // Hour in 12h format (01-12)
			hour := t.Hour() % 12
			if hour == 0 {
				hour = 12
			}
			if padWidth > 0 {
				result.WriteString(formatInt(hour, padWidth, padChar))
			} else {
				result.WriteString(strconv.Itoa(hour))
			}
		case 'j': // Day of year (001-366)
			yday := t.YearDay()
			if padWidth > 0 {
				result.WriteString(formatInt(yday, 3, padChar))
			} else {
				result.WriteString(strconv.Itoa(yday))
			}
		case 'k': // Hour in 24h format (space-padded)
			result.WriteString(fmt.Sprintf("%2d", t.Hour()))
		case 'l': // Hour in 12h format (space-padded)
			hour := t.Hour() % 12
			if hour == 0 {
				hour = 12
			}
			result.WriteString(fmt.Sprintf("%2d", hour))
		case 'M': // Minute (00-59)
			if padWidth > 0 {
				result.WriteString(formatInt(t.Minute(), padWidth, padChar))
			} else {
				result.WriteString(strconv.Itoa(t.Minute()))
			}
		case 'm': // Month (01-12)
			if padWidth > 0 {
				result.WriteString(formatInt(int(t.Month()), padWidth, padChar))
			} else {
				result.WriteString(strconv.Itoa(int(t.Month())))
			}
		case 'n': // Newline
			result.WriteString("\n")
		case 'p': // AM/PM
			if t.Hour() < 12 {
				result.WriteString(loc.AM)
			} else {
				result.WriteString(loc.PM)
			}
		case 'R': // %H:%M
			result.WriteString(t.Format("15:04"))
		case 'r': // %I:%M:%S %p
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			ampm := loc.AM
			if t.Hour() >= 12 {
				ampm = loc.PM
			}
			result.WriteString(fmt.Sprintf("%02d:%02d:%02d %s", h, t.Minute(), t.Second(), ampm))
		case 'S': // Second (00-59)
			if padWidth > 0 {
				result.WriteString(formatInt(t.Second(), padWidth, padChar))
			} else {
				result.WriteString(strconv.Itoa(t.Second()))
			}
		case 's': // Seconds since Unix epoch
			result.WriteString(strconv.FormatInt(t.Unix(), 10))
		case 'T': // %H:%M:%S
			result.WriteString(t.Format("15:04:05"))
		case 't': // Tab
			result.WriteString("\t")
		case 'U': // Week number (Sunday first day)
			_, week := t.ISOWeek()
			if padWidth > 0 {
				result.WriteString(formatInt(week, 2, padChar))
			} else {
				result.WriteString(strconv.Itoa(week))
			}
		case 'u': // Weekday (1-7, Monday is 1)
			wd := int(t.Weekday())
			if wd == 0 {
				wd = 7 // Sunday should be 7 in this format
			}
			result.WriteString(strconv.Itoa(wd))
		case 'V': // ISO 8601 week number
			_, week := t.ISOWeek()
			if padWidth > 0 {
				result.WriteString(formatInt(week, 2, padChar))
			} else {
				result.WriteString(strconv.Itoa(week))
			}
		case 'v': // %e-%b-%Y
			result.WriteString(fmt.Sprintf("%2d-%s-%04d", t.Day(), loc.MonthsAbbrev[t.Month()-1], t.Year()))
		case 'W': // Week number (Monday first day)
			_, week := t.ISOWeek()
			if padWidth > 0 {
				result.WriteString(formatInt(week, 2, padChar))
			} else {
				result.WriteString(strconv.Itoa(week))
			}
		case 'w': // Weekday (0-6, Sunday is 0)
			result.WriteString(strconv.Itoa(int(t.Weekday())))
		case 'X': // Time representation
			result.WriteString(t.Format("15:04:05"))
		case 'x': // Date representation
			result.WriteString(t.Format("01/02/06"))
		case 'Y': // Year with century
			if padWidth > 0 {
				result.WriteString(formatInt(t.Year(), 4, padChar))
			} else {
				result.WriteString(strconv.Itoa(t.Year()))
			}
		case 'y': // Year without century
			if padWidth > 0 {
				result.WriteString(formatInt(t.Year()%100, 2, padChar))
			} else {
				result.WriteString(strconv.Itoa(t.Year() % 100))
			}
		case 'Z': // Time zone name
			result.WriteString(t.Format("MST"))
		case 'z': // Time zone offset
			result.WriteString(t.Format("-0700"))
		case '+': // Date and time like date(1)
			result.WriteString(t.Format("Mon Jan 2 15:04:05 MST 2006"))
		case '%': // Literal %
			result.WriteByte('%')
		default:
			result.WriteByte(format[i])
		}
		i++
	}

	return result.String()
}

// formatInt formats an integer with specified padding
func formatInt(value, width int, padChar byte) string {
	s := strconv.Itoa(value)
	if len(s) >= width {
		return s
	}

	padding := strings.Repeat(string(padChar), width-len(s))
	return padding + s
}
