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

type ServerStatus struct {
	ID     uint `gorm:"primaryKey"`
	Status int  `json:"status"`
}
