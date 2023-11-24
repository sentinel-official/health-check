package types

import (
	"strings"
)

const (
	StatusUnknown Status = 0x00 + iota
	StatusActive
	StatusInactivePending
	StatusInactive
)

type (
	Status byte
)

func (z Status) Is(v Status) bool {
	return z == v
}

func (z Status) IsActive() bool {
	return z.Is(StatusActive)
}

func NewStatusFromString(v string) Status {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	switch v {
	case "active":
		return StatusActive
	case "inactive":
		return StatusInactive
	case "inactive_pending":
		return StatusInactivePending
	default:
		return StatusUnknown
	}
}

func (z Status) String() string {
	switch z {
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	case StatusInactivePending:
		return "inactive_pending"
	default:
		return "unknown"
	}
}
