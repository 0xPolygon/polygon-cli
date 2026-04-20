package render

import (
	"fmt"
	"strconv"
	"time"
)

// AnnotateUnixSeconds formats a unix-second timestamp like cast does:
// the integer, a human UTC string, and a coarse "ago" relative time.
// Returns the input unchanged if it cannot be parsed as a positive
// integer.
//
//	1776640801  (2026-04-19 23:20:01 UTC, 2h 4m ago)
func AnnotateUnixSeconds(raw string) string {
	return annotateAt(raw, time.Now())
}

// annotateAt is the deterministic variant used by tests.
func annotateAt(raw string, now time.Time) string {
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || n <= 0 {
		return raw
	}
	t := time.Unix(n, 0).UTC()
	delta := now.Sub(t)
	return fmt.Sprintf("%s  (%s, %s)", raw, t.Format("2006-01-02 15:04:05 UTC"), humanAgo(delta))
}

// humanAgo returns a coarse human duration. Positive deltas are "ago",
// negative deltas are "from now".
func humanAgo(d time.Duration) string {
	suffix := "ago"
	if d < 0 {
		d = -d
		suffix = "from now"
	}
	switch {
	case d < time.Minute:
		secs := int(d.Seconds())
		return fmt.Sprintf("%ds %s", secs, suffix)
	case d < time.Hour:
		return fmt.Sprintf("%dm %s", int(d.Minutes()), suffix)
	case d < 24*time.Hour:
		h := int(d.Hours())
		m := int(d.Minutes()) - h*60
		return fmt.Sprintf("%dh %dm %s", h, m, suffix)
	default:
		days := int(d.Hours() / 24)
		h := int(d.Hours()) - days*24
		return fmt.Sprintf("%dd %dh %s", days, h, suffix)
	}
}
