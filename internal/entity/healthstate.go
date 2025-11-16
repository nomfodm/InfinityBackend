package entity

const (
	ServerStatusWorking = iota
	ServerStatusMaintenance
	ServerStatusOff
)

var ServerStatuses = map[int]string{
	ServerStatusWorking:     "working",
	ServerStatusMaintenance: "maintenance",
	ServerStatusOff:         "off",
}

type HealthState struct {
	ID     uint `gorm:"primaryKey"`
	Status int  `json:"status"`
}
