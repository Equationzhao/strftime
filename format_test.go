package strftime

import (
	"fmt"
	"testing"
	"time"
)

func TestStrftime_DefaultLocale(t *testing.T) {
	// Fixed time: 2025-02-25 15:30:45, formatted using the default English locale
	loc, _ := time.LoadLocation("Local")
	testTime := time.Date(2025, time.February, 25, 15, 30, 45, 0, loc)
	formatted := Strftime("%Y-%m-%d %H:%M:%S", testTime)
	expected := "2025-02-25 15:30:45"
	if formatted != expected {
		t.Errorf("Strftime failed to format with default locale, got [%s], expected [%s]", formatted, expected)
	}
}

func TestStrftime_CustomChineseLocale(t *testing.T) {
	// Fixed time: 2025-02-25 15:30:45, formatted using a custom Chinese locale
	loc, _ := time.LoadLocation("Local")
	testTime := time.Date(2025, time.February, 25, 15, 30, 45, 0, loc)
	chineseLocale := &Locale{
		WeekdaysFull:   []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"},
		WeekdaysAbbrev: []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"},
		MonthsFull:     []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"},
		MonthsAbbrev:   []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		AM:             "上午",
		PM:             "下午",
	}
	formatted := StrftimeL("%A, %B %d, %Y %H:%M:%S", testTime, chineseLocale)
	expected := "星期二, 二月 25, 2025 15:30:45"
	if formatted != expected {
		t.Errorf("StrftimeL failed to format with Chinese locale, got [%s], expected [%s]", formatted, expected)
	}
}

func TestStrftime_NilLocale(t *testing.T) {
	loc, _ := time.LoadLocation("Local")
	testTime := time.Date(2025, time.February, 25, 15, 30, 45, 0, loc)
	formatted := StrftimeL("%A, %B %d, %Y %H:%M:%S", testTime, nil)
	expected := "Tuesday, February 25, 2025 15:30:45"
	if formatted != expected {
		t.Errorf("StrftimeL failed to format with Chinese locale, got [%s], expected [%s]", formatted, expected)
	}
}

func TestStrftime_LiteralPercent(t *testing.T) {
	// Test that %% in the format string produces a literal percent sign
	testTime := time.Date(2025, 2, 25, 15, 30, 45, 0, time.Local)
	formatted := Strftime("%% %Y", testTime)
	expected := "% 2025"
	if formatted != expected {
		t.Errorf("Literal %% formatting failed, got [%s], expected [%s]", formatted, expected)
	}
}

func TestStrftime_AllSpecifiers(t *testing.T) {
	// Test all supported conversion specifiers (default English locale)
	// Fixed time: 2025-02-03 09:05:07 (UTC)
	loc := time.FixedZone("UTC", 0)
	testTime := time.Date(2025, time.February, 3, 9, 5, 7, 0, loc)
	// %c uses the built-in format: Mon Jan 2 15:04:05 2006
	layoutC := testTime.Format("Mon Jan 2 15:04:05 2006")
	format := "%A %a %B %b %C %c %D %d %e %H %I %M %S %p %y %Y %m %j %F %Z %z %%"
	formatted := Strftime(format, testTime)
	expected := fmt.Sprintf("Monday Mon February Feb 20 %s 02/03/25 03  3 09 09 05 07 AM 25 2025 02 034 2025-02-03 UTC +0000 %%", layoutC)
	if formatted != expected {
		t.Errorf("Strftime_AllSpecifiers failed,\n got: [%s]\n expected: [%s]", formatted, expected)
	}
}

func TestStrftime_PosixExtension(t *testing.T) {
	// Test that the POSIX extension prefixes %E and %O in the format string are correctly skipped
	// The format can be written as either "%E%Y" or "%EY"
	loc := time.FixedZone("UTC", 0)
	testTime := time.Date(2025, 9, 10, 14, 55, 33, 0, loc)

	formatted1 := Strftime("%E%%Y", testTime)
	expected1 := "%2025"
	if formatted1 != expected1 {
		t.Errorf("Expected %%E%%Y to be [%s], got [%s]", expected1, formatted1)
	}
	// Using "%O%H"
	formatted2 := Strftime("%O%%H", testTime)
	expected2 := "%14"
	if formatted2 != expected2 {
		t.Errorf("Expected %%O%%H to be [%s], got [%s]", expected2, formatted2)
	}
}

func TestStrftime_UnknownSpecifier(t *testing.T) {
	testTime := time.Date(2025, 2, 25, 15, 30, 45, 0, time.Local)
	formatted := Strftime("%Q", testTime)
	expected := "Q"
	if formatted != expected {
		t.Errorf("Unknown specifier failed, got [%s], expected [%s]", formatted, expected)
	}
}

func TestStrftime_FormatModifiers(t *testing.T) {
	// Test format modifiers (-, _, 0)
	testTime := time.Date(2025, time.February, 3, 9, 5, 7, 0, time.UTC)

	// Test the - modifier (no padding)
	formatted := Strftime("%-d/%-m/%-Y %-H:%-M:%-S", testTime)
	expected := "3/2/2025 9:5:7"
	if formatted != expected {
		t.Errorf("Format modifier '-' failed, got [%s], expected [%s]", formatted, expected)
	}

	// Test the _ modifier (space padding)
	formatted = Strftime("%_d/%_m/%_Y %_H:%_M:%_S", testTime)
	expected = " 3/ 2/2025  9: 5: 7"
	if formatted != expected {
		t.Errorf("Format modifier '_' failed, got [%s], expected [%s]", formatted, expected)
	}

	// Test the 0 modifier (zero padding, default behavior)
	formatted = Strftime("%0d/%0m/%0Y %0H:%0M:%0S", testTime)
	expected = "03/02/2025 09:05:07"
	if formatted != expected {
		t.Errorf("Format modifier '0' failed, got [%s], expected [%s]", formatted, expected)
	}
}

func TestStrftime_UnusedSpecifiers(t *testing.T) {
	// Test less common format specifiers
	testTime := time.Date(2025, time.February, 3, 15, 5, 7, 0, time.UTC)

	tests := []struct {
		format   string
		expected string
	}{
		{"%n", "\n"},
		{"%t", "\t"},
		{"%k", "15"},                          // Hour in 24h format (space-padded)
		{"%l", " 3"},                          // Hour in 12h format (space-padded)
		{"%r", "03:05:07 PM"},                 // 12-hour time format
		{"%R", "15:05"},                       // 24-hour time format, without seconds
		{"%s", fmt.Sprint(testTime.Unix())},   // Unix timestamp
		{"%u", "1"},                           // ISO 8601 weekday (1-7, Monday is 1)
		{"%V", "06"},                          // ISO 8601 week number
		{"%w", "1"},                           // Weekday (0-6, Sunday is 0)
		{"%+", "Mon Feb 3 15:05:07 UTC 2025"}, // Format similar to date(1)
		{"%v", " 3-Feb-2025"},                 // %e-%b-%Y format
	}

	for _, tt := range tests {
		formatted := Strftime(tt.format, testTime)
		if formatted != tt.expected {
			t.Errorf("For format [%s]: got [%s], expected [%s]", tt.format, formatted, tt.expected)
		}
	}
}

func TestStrftime_EdgeCases(t *testing.T) {
	testTime := time.Date(2025, time.February, 3, 9, 5, 7, 0, time.UTC)

	formatted := Strftime("", testTime)
	if formatted != "" {
		t.Errorf("Empty format string should return empty result, got [%s]", formatted)
	}

	formatted = Strftime("Date: %", testTime)
	expected := "Date: "
	if formatted != expected {
		t.Errorf("Format ending with %% should ignore it, got [%s], expected [%s]", formatted, expected)
	}

	formatted = Strftime("%%% %%%%", testTime)
	expected = "% %%"
	if formatted != expected {
		t.Errorf("Multiple %% should be processed correctly, got [%s], expected [%s]", formatted, expected)
	}

	formatted = Strftime("%%%% %%%%", testTime)
	expected = "%% %%"
	if formatted != expected {
		t.Errorf("Multiple %% should be processed correctly, got [%s], expected [%s]", formatted, expected)
	}

	formatted = Strftime("%% %%", testTime)
	expected = "% %"
	if formatted != expected {
		t.Errorf("Multiple %% should be processed correctly, got [%s], expected [%s]", formatted, expected)
	}
}

func TestStrftime_PosixExtensionsWithActualFormats(t *testing.T) {
	testTime := time.Date(2025, time.February, 3, 9, 5, 7, 0, time.UTC)

	tests := []struct {
		format   string
		expected string
	}{
		{"%EY", "2025"},
		{"%OY", "2025"},
		{"%Ed", "03"},
		{"%OH", "09"},
		{"%Em", "02"},
	}

	for _, tt := range tests {
		formatted := Strftime(tt.format, testTime)
		if formatted != tt.expected {
			t.Errorf("POSIX extension format [%s]: got [%s], expected [%s]", tt.format, formatted, tt.expected)
		}
	}
}

func TestStrftime_LocaleAMPM(t *testing.T) {
	testTimeMorning := time.Date(2025, time.February, 3, 9, 5, 7, 0, time.UTC)
	testTimeAfternoon := time.Date(2025, time.February, 3, 15, 5, 7, 0, time.UTC)

	customLocale := &Locale{
		WeekdaysFull:   DefaultLocale.WeekdaysFull,
		WeekdaysAbbrev: DefaultLocale.WeekdaysAbbrev,
		MonthsFull:     DefaultLocale.MonthsFull,
		MonthsAbbrev:   DefaultLocale.MonthsAbbrev,
		AM:             "上午",
		PM:             "下午",
	}

	formatted := StrftimeL("%p", testTimeMorning, customLocale)
	expected := "上午"
	if formatted != expected {
		t.Errorf("AM format with custom locale failed, got [%s], expected [%s]", formatted, expected)
	}

	formatted = StrftimeL("%p", testTimeAfternoon, customLocale)
	expected = "下午"
	if formatted != expected {
		t.Errorf("PM format with custom locale failed, got [%s], expected [%s]", formatted, expected)
	}

	formatted = StrftimeL("%r", testTimeAfternoon, customLocale)
	expected = "03:05:07 下午"
	if formatted != expected {
		t.Errorf("12-hour time format with custom locale failed, got [%s], expected [%s]", formatted, expected)
	}
}

func TestStrftime_CompleteFormatCoverage(t *testing.T) {
	winter := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC) // January 1st is the 1st day of the year, January 1, 2025 is Wednesday
	summer := time.Date(2025, time.July, 4, 12, 34, 56, 0, time.UTC) // July 4th is the 185th day of the year, it is Friday

	tests := []struct {
		time     time.Time
		format   string
		expected string
	}{
		{winter, "%j", "001"},
		{summer, "%j", "185"},
		{winter, "%C", "20"},
		{summer, "%g", "25"},
		{winter, "%G", "2025"},
	}

	for _, tt := range tests {
		formatted := Strftime(tt.format, tt.time)
		if formatted != tt.expected {
			t.Errorf("Format [%s] with time [%v]: got [%s], expected [%s]", tt.format, tt.time, formatted, tt.expected)
		}
	}
}
