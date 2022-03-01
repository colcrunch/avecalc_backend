package models

import "time"

type User struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time
	UserName     string `json:"username" gorm:"unique"`
	PasswordHash string
	Email        string `json:"email"`
	Admin        bool   `json:"admin"`
}
