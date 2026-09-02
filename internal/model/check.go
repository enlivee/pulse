package model

import "time"

type Check struct {
	MonitorID  string        `json:"monitor_id"`
	Time       time.Time     `json:"time"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency"`
	Success    bool          `json:"success"`
}
