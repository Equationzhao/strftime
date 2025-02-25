package strftime

import (
	"strings"
	"testing"
	"time"
)

func TestParse_DefaultLocale_24Hour(t *testing.T) {
	// 24-hour format parsing: 2025-02-25 15:30:45
	loc, _ := time.LoadLocation("Local")
	expectedTime := time.Date(2025, time.February, 25, 15, 30, 45, 0, loc)
	input := "2025-02-25 15:30:45"
	format := "%Y-%m-%d %H:%M:%S"
	parsedTime, err := Parse(format, input)
	if err != nil {
		t.Fatalf("Parse default locale 24-hour format parsing error: %v", err)
	}
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Parse 24-hour format parsing failed, expected %v, got %v", expectedTime, parsedTime)
	}
}

func TestParse_DefaultLocale_12Hour(t *testing.T) {
	// 12-hour format parsing: "2025-02-25 03:30:45 PM" corresponds to 15:30:45
	loc, _ := time.LoadLocation("Local")
	expectedTime := time.Date(2025, time.February, 25, 15, 30, 45, 0, loc)
	input := "2025-02-25 03:30:45 PM"
	format := "%Y-%m-%d %I:%M:%S %p"
	parsedTime, err := Parse(format, input)
	if err != nil {
		t.Fatalf("Parse default locale 12-hour format parsing error: %v", err)
	}
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Parse 12-hour format parsing failed, expected %v, got %v", expectedTime, parsedTime)
	}
}

func TestParse_CustomChineseLocale(t *testing.T) {
	// Chinese locale parsing example: weekdays and months use full Chinese names
	chineseLocale := &Locale{
		WeekdaysFull:   []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"},
		WeekdaysAbbrev: []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"},
		MonthsFull:     []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"},
		MonthsAbbrev:   []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		AM:             "上午",
		PM:             "下午",
	}
	loc, _ := time.LoadLocation("Local")
	expectedTime := time.Date(2025, time.February, 25, 15, 30, 45, 0, loc)
	input := "星期二, 二月 25, 2025 03:30:45 下午"
	format := "%A, %B %d, %Y %I:%M:%S %p"
	parsedTime, err := ParseL(format, input, chineseLocale)
	if err != nil {
		t.Fatalf("ParseL Chinese locale parsing error: %v", err)
	}
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("ParseL Chinese locale parsing failed, expected %v, got %v", expectedTime, parsedTime)
	}
}

func TestParse_LiteralPercent(t *testing.T) {
	// Test parsing of literal %%
	input := "% 2025"
	format := "%% %Y"
	parsedTime, err := Parse(format, input)
	if err != nil {
		t.Fatalf("Parse literal %% parsing error: %v", err)
	}
	if parsedTime.Year() != 2025 {
		t.Errorf("Parse literal %% parsing failed, expected year 2025, got %d", parsedTime.Year())
	}
}

func TestParse_UnsupportedSpecifier(t *testing.T) {
	// Using unsupported conversion specifier (e.g., %C) should return an error
	input := "2025-02-25"
	format := "%C-%m-%d"
	_, err := Parse(format, input)
	if err == nil {
		t.Errorf("Expected error for unsupported conversion specifier, but got none")
	}
}

func TestParse_IncompleteSpecifier(t *testing.T) {
	// Format string ending with a single '%', expected to return an error
	input := "2025"
	format := "%Y%"
	_, err := Parse(format, input)
	if err == nil {
		t.Errorf("Expected error for incomplete format specifier, but got none")
	}
}

func TestParse_LiteralMismatch(t *testing.T) {
	// Literal mismatch between format and input string should return an error
	input := "2025/02/25"
	format := "%Y-%m-%d"
	_, err := Parse(format, input)
	if err == nil {
		t.Errorf("Expected error for literal mismatch, but got none")
	}
}

func TestParse_ExtraTrailing(t *testing.T) {
	// Input string has trailing non-whitespace characters, should return an error
	input := "2025-02-25 15:30:45 EXTRA"
	format := "%Y-%m-%d %H:%M:%S"
	_, err := Parse(format, input)
	if err == nil {
		t.Errorf("Expected error for trailing characters, but got none")
	}
}

func TestParse_12HourMissingAMPM(t *testing.T) {
	// Using 12-hour format but missing AM/PM marker, should return an error
	input := "2025-02-25 03:30:45"
	format := "%Y-%m-%d %I:%M:%S"
	_, err := Parse(format, input)
	if err == nil {
		t.Errorf("Expected error for missing AM/PM marker, but got none")
	}
}

func TestParse_InvalidNumber(t *testing.T) {
	// Input string has invalid number format, should return an error
	input := "20a5-02-25"
	format := "%Y-%m-%d"
	_, err := Parse(format, input)
	if err == nil {
		t.Errorf("Expected error for invalid number format, but got none")
	}
}

func TestParse_PosixExtension(t *testing.T) {
	// Test parsing with POSIX extension prefix %E being skipped and correctly handling subsequent specifier
	input := "2025-03-15"
	format := "%E%Y-%m-%d"
	parsedTime, err := Parse(format, input)
	if err != nil {
		t.Fatalf("Parse POSIX extension parsing error: %v", err)
	}
	if parsedTime.Year() != 2025 || int(parsedTime.Month()) != 3 || parsedTime.Day() != 15 {
		t.Errorf("Parse POSIX extension parsing failed, got %v", parsedTime)
	}
}

func TestParse_IncompletePosixExtension(t *testing.T) {
	// When encountering %E at the end of the format string without a subsequent specifier, should return an error
	input := "2025"
	format := "%E"
	_, err := Parse(format, input)
	if err == nil {
		t.Errorf("Expected error for incomplete POSIX extension specifier, but got none")
	}
}

func TestParse_TrailingWhitespace(t *testing.T) {
	// Input string allows trailing whitespace characters
	loc, _ := time.LoadLocation("Local")
	expectedTime := time.Date(2025, time.February, 25, 15, 30, 45, 0, loc)
	input := "2025-02-25 15:30:45    " // Trailing spaces
	format := "%Y-%m-%d %H:%M:%S"
	parsedTime, err := Parse(format, input)
	if err != nil {
		t.Fatalf("Parse trailing whitespace parsing error: %v", err)
	}
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Parse trailing whitespace parsing failed, expected %v, got %v", expectedTime, parsedTime)
	}
}

// -------------------------
// Tests for ParseL error and edge cases
// -------------------------

func TestParseL_IncompleteFormat(t *testing.T) {
	_, err := ParseL("%", "anything", nil)
	if err == nil || !strings.Contains(err.Error(), "incomplete format specifier at end") {
		t.Errorf("Expected 'incomplete format specifier at end' error, got %v", err)
	}
}

func TestParseL_IncompletePosixExtension(t *testing.T) {
	_, err := ParseL("%E", "2025", nil)
	if err == nil || !strings.Contains(err.Error(), "incomplete format specifier after posix extension") {
		t.Errorf("Expected 'incomplete format specifier after posix extension' error, got %v", err)
	}
}

func TestParseL_AMPMError(t *testing.T) {
	// Using %p, but input does not contain a valid AM/PM marker (default locale AM=="AM", PM=="PM")
	_, err := ParseL("%Y %p", "2025 NK", nil)
	if err == nil || !strings.Contains(err.Error(), "expected AM/PM marker") {
		t.Errorf("Expected AM/PM marker error, got %v", err)
	}
}

func TestParseL_D_ErrorMissingSlash(t *testing.T) {
	// %D should be "%m/%d/%y", if the separator is incorrect, should return an error
	_, err := ParseL("%D", "12-31/99", nil)
	if err == nil || !strings.Contains(err.Error(), "expected '/' after month in %D") {
		t.Errorf("Expected %%D missing '/' error, got %v", err)
	}
}

func TestParseL_F_ErrorMissingHyphen(t *testing.T) {
	// %F should be "%Y-%m-%d", if missing '-', should return an error
	_, err := ParseL("%F", "2025/02-25", nil)
	if err == nil || !strings.Contains(err.Error(), "expected '-' after year in %F") {
		t.Errorf("Expected %%F missing '-' error, got %v", err)
	}
}

func TestParseL_FullMonthNameError(t *testing.T) {
	_, err := ParseL("%B", "NotAMonth", nil)
	if err == nil || !strings.Contains(err.Error(), "failed to parse full month name") {
		t.Errorf("Expected full month name parsing error, got %v", err)
	}
}

func TestParseL_AbbrevMonthNameError(t *testing.T) {
	_, err := ParseL("%b", "XYZ", nil)
	if err == nil || !strings.Contains(err.Error(), "failed to parse abbreviated month name") {
		t.Errorf("Expected abbreviated month name parsing error, got %v", err)
	}
}

func TestParseL_FullWeekdayError(t *testing.T) {
	_, err := ParseL("%A", "NotAWeekday", nil)
	if err == nil || !strings.Contains(err.Error(), "failed to parse full weekday name") {
		t.Errorf("Expected full weekday name parsing error, got %v", err)
	}
}

func TestParseL_AbbrevWeekdayError(t *testing.T) {
	_, err := ParseL("%a", "NotAWeekday", nil)
	if err == nil || !strings.Contains(err.Error(), "failed to parse abbreviated weekday name") {
		t.Errorf("Expected abbreviated weekday name parsing error, got %v", err)
	}
}

func TestParseL_LiteralPercentError(t *testing.T) {
	// When format requires literal '%' but input does not match, should return an error
	_, err := ParseL("%%", "not%", nil)
	if err == nil || !strings.Contains(err.Error(), "expected literal '%'") {
		t.Errorf("Expected literal %% error, got %v", err)
	}
}

func TestParseL_ExtraTrailing(t *testing.T) {
	_, err := ParseL("%Y", "2025X", nil)
	if err == nil || !strings.Contains(err.Error(), "unparsed trailing characters") {
		t.Errorf("Expected trailing characters error, got %v", err)
	}
}

func TestParseL_12HourMissingAMPM(t *testing.T) {
	_, err := ParseL("%Y-%m-%d %I:%M:%S", "2025-02-25 03:30:45", nil)
	if err == nil || !strings.Contains(err.Error(), "12-hour format specified but missing AM/PM marker") {
		t.Errorf("Expected 12-hour missing AM/PM error, got %v", err)
	}
}

func TestParseL_Invalid12Hour(t *testing.T) {
	// For 12-hour format, if the hour is not in the range [1,12], should return an error. Here "00" is invalid.
	_, err := ParseL("%Y-%m-%d %I:%M:%S %p", "2025-02-25 00:30:45 AM", nil)
	if err == nil || !strings.Contains(err.Error(), "invalid hour") {
		t.Errorf("Expected invalid hour error, got %v", err)
	}
}

func TestParseL_Valid12HourAM(t *testing.T) {
	// For 12-hour format, 12 AM should be converted to 0 hour. For example, "12:30:45 AM" -> 00:30:45
	tme, err := ParseL("%Y-%m-%d %I:%M:%S %p", "2025-02-25 12:30:45 AM", nil)
	if err != nil {
		t.Fatalf("ParseL parsing error: %v", err)
	}
	if tme.Hour() != 0 {
		t.Errorf("Expected 12:30:45 AM to be 00:30:45, got %02d:30:45", tme.Hour())
	}
}

func TestParseL_Valid12HourPM(t *testing.T) {
	// For 12-hour format, 12 PM remains 12 hour
	tme, err := ParseL("%Y-%m-%d %I:%M:%S %p", "2025-02-25 12:30:45 PM", nil)
	if err != nil {
		t.Fatalf("ParseL parsing error: %v", err)
	}
	if tme.Hour() != 12 {
		t.Errorf("Expected 12:30:45 PM to be 12:30:45, got %02d:30:45", tme.Hour())
	}
}

// -------------------------
// Tests for internal helper functions
// -------------------------

func TestParseHelpers(t *testing.T) {
	// Test correct branch of parseFixedInt
	input := "98765"
	val, pos, err := parseFixedInt(input, 0, 3)
	if err != nil || val != 987 || pos != 3 {
		t.Errorf("parseFixedInt valid failed: got (%d, %d, %v), expected (987, 3, nil)", val, pos, err)
	}

	// Test parseFixedInt returns error when characters are insufficient
	_, _, err = parseFixedInt(input, 4, 3)
	if err == nil {
		t.Error("parseFixedInt expected error for insufficient characters, but got none")
	}

	// Test correct branch of parseIntVariable
	tests := []struct {
		s         string
		pos       int
		minDigits int
		maxDigits int
		expVal    int
		expPos    int
		expErr    bool
	}{
		{"12345", 0, 2, 4, 1234, 4, false},
		{"12abc", 0, 2, 4, 12, 2, false}, // Stop at non-digit character
		{"a123", 0, 1, 3, 0, 0, true},    // First character is non-digit
	}
	for idx, tt := range tests {
		v, newPos, err := parseIntVariable(tt.s, tt.pos, tt.minDigits, tt.maxDigits)
		if tt.expErr {
			if err == nil {
				t.Errorf("Test parseIntVariable #%d: expected error but got none", idx)
			}
		} else {
			if err != nil {
				t.Errorf("Test parseIntVariable #%d: unexpected error: %v", idx, err)
			}
			if v != tt.expVal || newPos != tt.expPos {
				t.Errorf("Test parseIntVariable #%d: expected (%d, %d) but got (%d, %d)", idx, tt.expVal, tt.expPos, v, newPos)
			}
		}
	}
}

func TestParse_FullMonthName_Default(t *testing.T) {
	// Default locale, %B corresponds to full month name (e.g., "February")
	input := "2025-February-25"
	format := "%Y-%B-%d"
	tme, err := Parse(format, input)
	if err != nil {
		t.Fatalf("TestParse_FullMonthName_Default error: %v", err)
	}
	if tme.Month() != time.February {
		t.Errorf("Expected month to be February, got %v", tme.Month())
	}
}

func TestParse_AbbrevMonthName_Default(t *testing.T) {
	// Default locale, %b corresponds to abbreviated month name (e.g., "Feb")
	input := "2025-Feb-25"
	format := "%Y-%b-%d"
	tme, err := Parse(format, input)
	if err != nil {
		t.Fatalf("TestParse_AbbrevMonthName_Default error: %v", err)
	}
	if tme.Month() != time.February {
		t.Errorf("Expected month to be February, got %v", tme.Month())
	}
}

func TestParse_FullWeekday_Default(t *testing.T) {
	// Default locale, %A corresponds to full weekday name (e.g., "Friday")
	// Note: ParseL consumes the weekday name but does not change the time data, so only verify no error here.
	input := "Friday2025-02-25"
	format := "%A%Y-%m-%d"
	tme, err := Parse(format, input)
	if err != nil {
		t.Fatalf("TestParse_FullWeekday_Default error: %v", err)
	}
	if tme.Year() != 2025 {
		t.Errorf("Expected year to be 2025, got %d", tme.Year())
	}
}

func TestParse_AbbrevWeekday_Default(t *testing.T) {
	// Default locale, %a corresponds to abbreviated weekday name (e.g., "Fri")
	input := "Fri2025-02-25"
	format := "%a%Y-%m-%d"
	tme, err := Parse(format, input)
	if err != nil {
		t.Fatalf("TestParse_AbbrevWeekday_Default error: %v", err)
	}
	if tme.Year() != 2025 {
		t.Errorf("Expected year to be 2025, got %d", tme.Year())
	}
}
