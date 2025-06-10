package entity

import "time"

type LauncherVersion struct {
	ID          uint `gorm:"primaryKey"`
	Version     string
	ReleaseDate time.Time
	DownloadUrl string
	SHA256      string
}
