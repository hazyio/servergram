package cron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const maxLineLen = 1000 // Debian cron rejects longer lines

var aliases = map[string]bool{
	"@reboot": true, "@yearly": true, "@annually": true, "@monthly": true,
	"@weekly": true, "@daily": true, "@midnight": true, "@hourly": true,
}

var monthNames = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}
var dowNames = map[string]int{
	"sun": 0, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6,
}

type fieldSpec struct {
	name     string
	min, max int
	names    map[string]int
}

var fieldSpecs = [5]fieldSpec{
	{"minute", 0, 59, nil},
	{"hour", 0, 23, nil},
	{"day-of-month", 1, 31, nil},
	{"month", 1, 12, monthNames},
	{"day-of-week", 0, 7, dowNames}, // 0 and 7 are both Sunday
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false // rejects "+5", " 5", unicode digits
		}
	}
	return true
}

func (f fieldSpec) value(s string) (int, error) {
	n, ok := 0, false
	if isDigits(s) {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("%s: %q out of range", f.name, s)
		}
		n, ok = v, true
	} else if v, found := f.names[strings.ToLower(s)]; found {
		n, ok = v, true
	}
	if !ok {
		return 0, fmt.Errorf("%s: invalid value %q", f.name, s)
	}
	if n < f.min || n > f.max {
		return 0, fmt.Errorf("%s: %d out of range %d-%d", f.name, n, f.min, f.max)
	}
	return n, nil
}

func validateField(f fieldSpec, s string) error {
	for _, part := range strings.Split(s, ",") { // "" part => error below
		base, stepStr, hasStep := strings.Cut(part, "/")
		if hasStep {
			step, err := strconv.Atoi(stepStr)
			if !isDigits(stepStr) || err != nil || step < 1 || step > f.max {
				return fmt.Errorf("%s: invalid step %q", f.name, stepStr)
			}
		}
		if base == "*" {
			continue
		}
		loStr, hiStr, isRange := strings.Cut(base, "-")
		lo, err := f.value(loStr)
		if err != nil {
			return err
		}
		if !isRange {
			if hasStep {
				return fmt.Errorf("%s: step needs * or a range, got %q", f.name, part)
			}
			continue
		}
		hi, err := f.value(hiStr)
		if err != nil {
			return err
		}
		if lo > hi {
			return fmt.Errorf("%s: range %q is backwards", f.name, base)
		}
	}
	return nil
}

// ValidateSchedule checks a 5-field spec or an @alias.
func ValidateSchedule(spec string) error {
	for _, r := range spec { // blocks \n, \r, NUL, unicode spaces
		if (r < 0x20 && r != '\t') || r > 0x7e {
			return errors.New("schedule contains invalid characters")
		}
	}
	spec = strings.TrimSpace(spec)
	if strings.HasPrefix(spec, "@") {
		if !aliases[spec] {
			return fmt.Errorf("unknown alias %q", spec)
		}
		return nil
	}
	parts := strings.Fields(spec)
	if len(parts) != 5 {
		return fmt.Errorf("expected 5 fields, got %d", len(parts))
	}
	for i, p := range parts {
		if err := validateField(fieldSpecs[i], p); err != nil {
			return err
		}
	}
	return nil
}

// ValidateCommand rejects anything that could break the one-line format.
func ValidateCommand(cmd string) error {
	if strings.TrimSpace(cmd) == "" {
		return errors.New("command is empty")
	}
	for _, r := range cmd {
		if r == '\t' {
			continue
		}
		if unicode.IsControl(r) { // \n, \r, NUL, ESC...
			return errors.New("command contains control characters or newlines")
		}
	}
	return nil
}

// nextField returns the first whitespace-delimited field and the rest.
func nextField(s string) (field, rest string) {
	s = strings.TrimLeft(s, " \t")
	i := strings.IndexAny(s, " \t")
	if i < 0 {
		return s, ""
	}
	return s[:i], strings.TrimLeft(s[i:], " \t")
}

// ParseLine splits and validates one crontab job line.
func ParseLine(line string) (schedule, command string, err error) {
	if len(line) > maxLineLen {
		return "", "", fmt.Errorf("line longer than %d characters", maxLineLen)
	}
	if strings.ContainsAny(line, "\n\r\x00") {
		return "", "", errors.New("line contains a newline or NUL")
	}
	rest := strings.TrimLeft(line, " \t")
	if strings.HasPrefix(rest, "#") || rest == "" || strings.Contains(strings.SplitN(rest, " ", 2)[0], "=") {
		return "", "", errors.New("not a job line (comment, blank, or env assignment)")
	}
	if strings.HasPrefix(rest, "@") {
		schedule, command = nextField(rest)
	} else {
		var parts [5]string
		for i := range parts {
			parts[i], rest = nextField(rest)
		}
		schedule, command = strings.Join(parts[:], " "), rest
	}
	if err = ValidateSchedule(schedule); err != nil {
		return "", "", err
	}
	if err = ValidateCommand(command); err != nil {
		return "", "", err
	}
	return schedule, command, nil
}

// BuildLine is what your New Job flow should use: the line is built from a
// validated schedule and an integer id, so no free text reaches the crontab.
func BuildLine(schedule, command string) (string, error) {
	if err := ValidateSchedule(schedule); err != nil {
		return "", err
	}
	if err := ValidateCommand(command); err != nil { // rejects \n, \r, NUL, control chars
		return "", err
	}
	schedule = strings.Join(strings.Fields(schedule), " ")
	command = strings.TrimSpace(command)

	// cron turns every unescaped % into a newline, so store them as \%
	command = strings.ReplaceAll(command, "%", `\%`)

	line := schedule + " " + shQuote(command)
	if len(line) > maxLineLen {
		return "", fmt.Errorf("line longer than %d characters", maxLineLen)
	}
	return line, nil
}

// BothDayFieldsRestricted reports the OR quirk: when neither day-of-month
// nor day-of-week starts with "*", cron runs on days matching EITHER.
func BothDayFieldsRestricted(schedule string) bool {
	f := strings.Fields(schedule)
	return len(f) == 5 && !strings.HasPrefix(f[2], "*") && !strings.HasPrefix(f[4], "*")
}

var aliasSchedules = map[string]string{
	"@yearly":   "0 0 1 1 *",
	"@annually": "0 0 1 1 *",
	"@monthly":  "0 0 1 * *",
	"@weekly":   "0 0 * * 0",
	"@daily":    "0 0 * * *",
	"@midnight": "0 0 * * *",
	"@hourly":   "0 * * * *",
}

const maxNextRunYears = 10

type cronSchedule struct {
	minutes  [60]bool
	hours    [24]bool
	days     [32]bool // 1-31
	months   [13]bool // 1-12
	weekdays [7]bool  // 0-6, Sunday = 0

	domRestricted bool
	dowRestricted bool
}

// parseCronField marks every value selected by expr in set. It reports whether
// the field is restricted (i.e. does not start with "*").
func parseCronField(f fieldSpec, expr string, set []bool) (bool, error) {
	restricted := !strings.HasPrefix(expr, "*")

	for _, part := range strings.Split(expr, ",") {
		base, stepStr, hasStep := strings.Cut(part, "/")
		step := 1
		if hasStep {
			n, err := strconv.Atoi(stepStr)
			if !isDigits(stepStr) || err != nil || n < 1 || n > f.max {
				return false, fmt.Errorf("%s: invalid step %q", f.name, stepStr)
			}
			step = n
		}

		lo, hi := f.min, f.max
		if base != "*" {
			loStr, hiStr, isRange := strings.Cut(base, "-")
			var err error
			lo, err = f.value(loStr)
			if err != nil {
				return false, err
			}
			hi = lo
			if isRange {
				hi, err = f.value(hiStr)
				if err != nil {
					return false, err
				}
			} else if hasStep {
				return false, fmt.Errorf("%s: step needs * or a range, got %q", f.name, part)
			}
			if lo > hi {
				return false, fmt.Errorf("%s: range %q is backwards", f.name, base)
			}
		}

		for v := lo; v <= hi; v += step {
			idx := v
			if f.name == "day-of-week" && idx == 7 {
				idx = 0 // 0 and 7 are both Sunday
			}
			set[idx] = true
		}
	}

	return restricted, nil
}

func parseCronSchedule(spec string) (*cronSchedule, error) {
	if err := ValidateSchedule(spec); err != nil {
		return nil, err
	}

	spec = strings.TrimSpace(spec)
	if strings.HasPrefix(spec, "@") {
		expanded, ok := aliasSchedules[spec]
		if !ok {
			if spec == "@reboot" {
				return nil, errors.New("@reboot has no scheduled run time")
			}
			return nil, fmt.Errorf("unknown alias %q", spec)
		}
		spec = expanded
	}

	parts := strings.Fields(spec)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	s := &cronSchedule{}
	var err error

	if _, err = parseCronField(fieldSpecs[0], parts[0], s.minutes[:]); err != nil {
		return nil, err
	}
	if _, err = parseCronField(fieldSpecs[1], parts[1], s.hours[:]); err != nil {
		return nil, err
	}
	if s.domRestricted, err = parseCronField(fieldSpecs[2], parts[2], s.days[:]); err != nil {
		return nil, err
	}
	if _, err = parseCronField(fieldSpecs[3], parts[3], s.months[:]); err != nil {
		return nil, err
	}
	if s.dowRestricted, err = parseCronField(fieldSpecs[4], parts[4], s.weekdays[:]); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *cronSchedule) matchesDay(t time.Time) bool {
	dom := s.days[t.Day()]
	dow := s.weekdays[int(t.Weekday())]

	// When both day fields are restricted cron matches either one.
	if s.domRestricted && s.dowRestricted {
		return dom || dow
	}
	return dom && dow
}

// next returns the first matching minute strictly after from.
func (s *cronSchedule) next(from time.Time) (time.Time, error) {
	loc := from.Location()
	t := from.Truncate(time.Minute).Add(time.Minute)
	limit := from.AddDate(maxNextRunYears, 0, 0)

	for t.Before(limit) {
		if !s.months[int(t.Month())] {
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, loc)
			continue
		}
		if !s.matchesDay(t) {
			t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, loc)
			continue
		}
		if !s.hours[t.Hour()] {
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, loc)
			continue
		}
		if !s.minutes[t.Minute()] {
			t = t.Add(time.Minute)
			continue
		}
		return t, nil
	}

	return time.Time{}, fmt.Errorf("no next run found within %d years", maxNextRunYears)
}

// NextRun returns the first time the schedule fires strictly after from.
// It supports 5-field cron expressions and the standard @aliases, except
// @reboot which has no scheduled run time.
func NextRun(schedule string, from time.Time) (time.Time, error) {
	s, err := parseCronSchedule(schedule)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid schedule %q: %w", schedule, err)
	}
	return s.next(from)
}

// NextRuns returns the next n fire times, useful for previews in the UI.
func NextRuns(schedule string, from time.Time, n int) ([]time.Time, error) {
	if n <= 0 {
		return []time.Time{}, nil
	}

	s, err := parseCronSchedule(schedule)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule %q: %w", schedule, err)
	}

	out := make([]time.Time, 0, n)
	t := from
	for range n {
		next, err := s.next(t)
		if err != nil {
			return nil, err
		}
		out = append(out, next)
		t = next
	}
	return out, nil
}
