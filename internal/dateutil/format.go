package dateutil

import "time"

func GetLocation(tz string) (*time.Location, error) {
	if tz == "" {
		return time.Local, nil
	}
	return time.LoadLocation(tz)
}

func GetLayout(ms bool) string {
	if ms {
		return "2006-01-02 15:04:05.000"
	}
	return "2006-01-02 15:04:05"
}
