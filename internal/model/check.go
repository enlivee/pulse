package model

import "time"

type Check struct {
	MonitorID  string
	Time       time.Time
	StatusCode int
	Latency    time.Duration
	Success    bool
}