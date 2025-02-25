package strftime

// Locale defines the date and time names required for locale settings
type Locale struct {
	WeekdaysFull   []string // Full names (starting from Sunday)
	WeekdaysAbbrev []string // Abbreviated names
	MonthsFull     []string // Full month names (starting from January)
	MonthsAbbrev   []string // Abbreviated month names
	AM             string   // AM identifier
	PM             string   // PM identifier
}

// Default English Locale
var DefaultLocale = &Locale{
	WeekdaysFull:   []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	WeekdaysAbbrev: []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
	MonthsFull: []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	},
	MonthsAbbrev: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	AM:           "AM",
	PM:           "PM",
}
