package domain

import "fmt"

type Frequency string

const (
	Hourly Frequency = "hourly"
	Daily  Frequency = "daily"
)

func ParseFrequency(s string) (Frequency, error) {
	switch s {
	case string(Hourly):
		return Hourly, nil
	case string(Daily):
		return Daily, nil
	default:
		return "", fmt.Errorf("invalid frequency: %q", s)
	}
}
