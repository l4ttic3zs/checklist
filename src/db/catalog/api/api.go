package api

type ItemType struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"unique;not null" json:"name"`
	ImagePath string `json:"image_path"`
}
