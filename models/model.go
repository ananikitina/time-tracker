package models

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey"`
	PassportNumber string `json:"passport_number" gorm:"column:passport_number;unique;not null"`
	Surname        string `json:"surname" gorm:"column:surname"`
	Name           string `json:"name" gorm:"column:name"`
	Patronymic     string `json:"patronymic" gorm:"column:patronymic"`
	Address        string `json:"address" gorm:"column:address"`
}

type Task struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	UserID    uint       `gorm:"not null"`
	TaskName  string     `gorm:"not null"`
	StartTime time.Time  `gorm:"not null"`
	EndTime   *time.Time `gorm:"default:null"`
}
