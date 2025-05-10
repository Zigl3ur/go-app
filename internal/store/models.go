package store

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:255;not null" json:"username"`
	Email     string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
}

type Session struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `json:"token"`
	UserId    uint      `json:"userId"`
	User      User      `gorm:"foreignKey:UserId;references:Id" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"createdAt"`
	ExpireAt  time.Time `json:"expireAt"`
}
