package api

type ItemType struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"unique;not null" json:"name"`
	ImagePath string `json:"image_path"`
}

type Item struct {
	ID         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemTypeID uint     `gorm:"unique;not null" json:"-"`
	ItemType   ItemType `gorm:"foreignKey:ItemTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"item_type"`
	Count      int      `json:"count"`
}
