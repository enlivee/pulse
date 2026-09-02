package model

import "time"

type Monitor struct {
	ID       string        `json:"id"`
	URL      string        `json:"url"`
	Interval time.Duration `json:"interval"`
}
