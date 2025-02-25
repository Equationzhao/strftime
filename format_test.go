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

	formatted1 := Strftime("%E%Y", testTime)
	expected1 := "2025"
	if formatted1 != expected1 {
		t.Errorf("Expected %%E%%Y to be [%s], got [%s]", expected1, formatted1)
	}
	// Using "%O%H"
	formatted2 := Strftime("%O%H", testTime)
	expected2 := "14"
	if formatted2 != expected2 {
		t.Errorf("Expected %%O%%H to be [%s], got [%s]", expected2, formatted2)
	}
}

func TestStrftime_UnknownSpecifier(t *testing.T) {
	// When encountering an unsupported conversion specifier, output '%' and the character as-is
	testTime := time.Date(2025, 2, 25, 15, 30, 45, 0, time.Local)
	formatted := Strftime("%X", testTime)
	expected := "%X"
	if formatted != expected {
		t.Errorf("Unknown specifier failed, got [%s], expected [%s]", formatted, expected)
	}
}
