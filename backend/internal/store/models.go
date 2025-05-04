package store

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

type Session struct {
	Id        uint `gorm:"primaryKey"`
	Token     string
	UserId    uint
	User      User      `gorm:"foreignKey:UserId;references:Id"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	ExpiresAt time.Time
}
