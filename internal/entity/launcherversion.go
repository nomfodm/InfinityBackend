package entity

type LauncherVersion struct {
	ID                  uint `gorm:"primaryKey"`
	CurrentVersion      string
	CurrentBinarySHA256 string
}
