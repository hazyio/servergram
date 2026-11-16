package cron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

	line := schedule + " " + command
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
