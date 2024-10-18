package entity

type Cape struct {
	ID          uint `gorm:"primaryKey"`
	TextureHash string

	UsersUsing []User
}
