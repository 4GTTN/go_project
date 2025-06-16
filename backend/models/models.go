package models

import (
	"time"
)

type Map struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Name        string     `json:"name"`
	Price       float64    `json:"price"`
	IsDeleted   bool       `json:"is_deleted"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	Image       string     `json:"image"`
}

type Player struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
	Username  string     `gorm:"unique" json:"username"`
	Email     string     `gorm:"unique" json:"email"`
	Role      string     `json:"role"`
	Token     string     `json:"token"`
	Password  string     `json:"password"`
}

type Game struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	CreatorID   uint       `json:"creator_id"`
	Status      string     `json:"status"`
	FormedAt    *time.Time `json:"formed_at"`
	CompletedAt *time.Time `json:"completed_at"`
	ModeratorID *uint      `json:"moderator_id"`
	Maps        []Map      `gorm:"many2many:game_maps;" json:"maps"`
}

type GameMap struct {
	GameID uint `gorm:"primaryKey" json:"game_id"`
	MapID  uint `gorm:"primaryKey" json:"map_id"`
}
