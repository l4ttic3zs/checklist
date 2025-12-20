package api

type Item struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	ImagePath string `json:"image_path"`
}

func (Item) TableName() string {
	return "items"
}
