package models

type Contract struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	Ref        string  `json:"ref" gorm:"unique"`
	Link       string  `json:"link"`
	Size       int     `json:"size"`
	Collateral float32 `json:"collateral"`
	Price      float32 `json:"price"`
	Status     int8    `json:"status" gorm:"default:0"`
}
