package entity

type Skin struct {
	ID          uint `gorm:"primaryKey"`
	TextureHash string

	UsersUsing []User
}
