package model

import "time"

type Monitor struct {
	ID       string
	URL      string
	Interval time.Duration
}
