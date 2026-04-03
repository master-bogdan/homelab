package healthdto

import "time"

type LivenessStatus struct {
	Status string        `json:"status"`
	Uptime time.Duration `json:"uptime" swaggertype:"string"`
}

type ReadinessStatus struct {
	Status string        `json:"status"`
	Uptime time.Duration `json:"uptime" swaggertype:"string"`
	DB     string        `json:"db"`
	Redis  string        `json:"redis"`
}
