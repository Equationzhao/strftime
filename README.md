# strftime

A Go implementation of C/Python-style strftime for formatting time using format specifiers, with locale support.

## Installation

```bash
go get -u github.com/equationzhao/strftime
```

## Overview

This package provides strftime-style time formatting and parsing functions for Go. It supports customizable locales and most common format specifiers.

## Usage

### Basic Formatting

```go
package main

import (
	"fmt"
	"time"
	
	"github.com/equationzhao/strftime"
)

func main() {
	now := time.Now()
	
	// Format with the default locale (English)
	formatted := strftime.Strftime("%Y-%m-%d %H:%M:%S", now)
	fmt.Println(formatted) // Output: 2023-04-05 15:30:45 (current time)
	
	// Using additional specifiers
	fullFormat := strftime.Strftime("%A, %B %d, %Y at %I:%M %p", now)
	fmt.Println(fullFormat) // Output: Wednesday, April 05, 2023 at 03:30 PM
}
```

### Custom Locales

```go
package main

import (
	"fmt"
	"time"
	
	"github.com/equationzhao/strftime"
)

func main() {
	now := time.Now()
	
	// Create a custom locale for Chinese formatting
	chineseLocale := &strftime.Locale{
		WeekdaysFull:   []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"},
		WeekdaysAbbrev: []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"},
		MonthsFull:     []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"},
		MonthsAbbrev:   []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		AM:             "上午",
		PM:             "下午",
	}
	
	// Format with the custom locale
	formatted := strftime.StrftimeL("%A, %B %d, %Y at %I:%M %p", now, chineseLocale)
	fmt.Println(formatted) // Output: 星期三, 四月 05, 2023 at 03:30 下午
}
```

### Parsing Time

```go
package main

import (
	"fmt"
	"time"
	
	"github.com/equationzhao/strftime"
)

func main() {
	// Parse time using default locale
	t, err := strftime.Parse("%Y-%m-%d %H:%M:%S", "2023-04-05 15:30:45")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t) // Output: 2023-04-05 15:30:45 +0000 UTC
	
	// Parse with custom locale
	chineseLocale := &strftime.Locale{
		// ... locale definition ...
	}
	t, err = strftime.ParseL("%Y年%m月%d日 %H时%M分", "2023年04月05日 15时30分", chineseLocale)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t)
}
```

## Supported Format Specifiers

| Specifier | Description | Example |
|-----------|-------------|---------|
| %A | Full weekday name | "Sunday", "Monday", ... |
| %a | Abbreviated weekday name | "Sun", "Mon", ... |
| %B | Full month name | "January", "February", ... |
| %b, %h | Abbreviated month name | "Jan", "Feb", ... |
| %C | Century (year/100) | "20" for 2023 |
| %c | Locale's date and time representation | "Mon Jan 2 15:04:05 2006" |
| %D | Same as %m/%d/%y | "04/05/23" |
| %d | Day of month (01-31) | "01", "02", ... |
| %e | Day of month (space-padded) | " 1", " 2", ... |
| %F | ISO 8601 date format (%Y-%m-%d) | "2023-04-05" |
| %H | Hour in 24-hour format (00-23) | "00", "01", ... |
| %I | Hour in 12-hour format (01-12) | "01", "02", ... |
| %j | Day of year (001-366) | "001", "002", ... |
| %M | Minute (00-59) | "00", "01", ... |
| %m | Month (01-12) | "01", "02", ... |
| %p | AM or PM | "AM", "PM" |
| %S | Second (00-59) | "00", "01", ... |
| %Y | Year with century | "2023" |
| %y | Year without century | "23" |
| %Z | Time zone name | "UTC", "EST", ... |
| %z | Time zone offset | "+0000", "-0700", ... |
| %% | A literal percent sign | "%" |

Note: POSIX extensions (like `%E` and `%O` prefixes) are supported by being skipped.

