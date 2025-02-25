package strftime

import (
	"fmt"
	"strings"
	"time"
)

// Strftime formats time using the default locale
func Strftime(format string, t time.Time) string {
	return StrftimeL(format, t, DefaultLocale)
}

// StrftimeL returns a formatted string based on the specified format, time t, and locale
// Supported conversion specifiers include: %A, %a, %B, %b/%h, %C, %c, %D, %d, %e, %H, %I, %M, %S, %p, %y, %Y, %m, %j, %F, %Z, %z, %%
// For POSIX extensions (such as those starting with %E or %O), this implementation will skip the extension prefix.
// Format strings can be written as either "%EY" or "%E%Y", both equivalent to "%Y".
func StrftimeL(format string, t time.Time, locale *Locale) string {
	if locale == nil {
		locale = DefaultLocale
	}
	runes := []rune(format)
	var b strings.Builder
	i := 0
	for i < len(runes) {
		if runes[i] == '%' {
			i++
			if i >= len(runes) {
				b.WriteRune('%')
				break
			}

			// Check if it's a POSIX extension prefix (%E or %O)
			if runes[i] == 'E' || runes[i] == 'O' {
				// Skip 'E' or 'O'
				i++
				if i >= len(runes) {
					b.WriteRune('%')
					b.WriteRune(runes[i-1])
					break
				}
				// If we encounter an additional '%' after skipping (e.g., "%E%Y"), skip it too, making it equivalent to "%Y"
				if runes[i] == '%' {
					i++
					if i >= len(runes) {
						b.WriteRune('%')
						break
					}
				}
			}
			spec := runes[i]
			switch spec {
			case 'A': // Full weekday name
				b.WriteString(locale.WeekdaysFull[int(t.Weekday())])
			case 'a': // Abbreviated weekday name
				b.WriteString(locale.WeekdaysAbbrev[int(t.Weekday())])
			case 'B': // Full month name
				monthIndex := int(t.Month()) - 1
				if monthIndex >= 0 && monthIndex < len(locale.MonthsFull) {
					b.WriteString(locale.MonthsFull[monthIndex])
				}
			case 'b', 'h': // Abbreviated month name
				monthIndex := int(t.Month()) - 1
				if monthIndex >= 0 && monthIndex < len(locale.MonthsAbbrev) {
					b.WriteString(locale.MonthsAbbrev[monthIndex])
				}
			case 'C': // Century (year/100, two digits)
				century := t.Year() / 100
				b.WriteString(fmt.Sprintf("%02d", century))
			case 'c': // Local date and time
				b.WriteString(t.Format("Mon Jan 2 15:04:05 2006"))
			case 'D': // Equivalent to "%m/%d/%y"
				b.WriteString(t.Format("01/02/06"))
			case 'd': // Day of month (two digits)
				b.WriteString(fmt.Sprintf("%02d", t.Day()))
			case 'e': // Day of month (width 2, space padded)
				b.WriteString(fmt.Sprintf("%2d", t.Day()))
			case 'H': // Hour in 24-hour format (two digits)
				b.WriteString(fmt.Sprintf("%02d", t.Hour()))
			case 'I': // Hour in 12-hour format (two digits)
				hour := t.Hour() % 12
				if hour == 0 {
					hour = 12
				}
				b.WriteString(fmt.Sprintf("%02d", hour))
			case 'M': // Minutes (two digits)
				b.WriteString(fmt.Sprintf("%02d", t.Minute()))
			case 'S': // Seconds (two digits)
				b.WriteString(fmt.Sprintf("%02d", t.Second()))
			case 'p': // AM/PM indicator
				if t.Hour() < 12 {
					b.WriteString(locale.AM)
				} else {
					b.WriteString(locale.PM)
				}
			case 'y': // Year without century (two digits)
				b.WriteString(fmt.Sprintf("%02d", t.Year()%100))
			case 'Y': // Year with century (four digits)
				b.WriteString(fmt.Sprintf("%04d", t.Year()))
			case 'm': // Month (two digits)
				b.WriteString(fmt.Sprintf("%02d", int(t.Month())))
			case 'j': // Day of year (three digits)
				b.WriteString(fmt.Sprintf("%03d", t.YearDay()))
			case 'F': // ISO 8601 date format (equivalent to "%Y-%m-%d")
				b.WriteString(t.Format("2006-01-02"))
			case 'Z': // Time zone abbreviation
				b.WriteString(t.Format("MST"))
			case 'z': // Numeric time zone offset, e.g., -0700
				b.WriteString(t.Format("-0700"))
			case '%': // Output a percent sign
				b.WriteRune('%')
			default:
				// Unrecognized conversion specifier, output as-is including the leading '%'
				b.WriteRune('%')
				b.WriteRune(spec)
			}
		} else {
			b.WriteRune(runes[i])
		}
		i++
	}
	return b.String()
}
