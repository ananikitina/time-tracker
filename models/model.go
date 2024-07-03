package models

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey"`
	PassportNumber string `json:"passportNumber" gorm:"unique;not null"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	TaskName  string `gorm:"not null"`
	StartTime time.Time
	EndTime   time.Time
}
