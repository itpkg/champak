package web

import "time"

//Model basic db model
type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"updated_at"`
	UpdatedAt time.Time `json:"created_at"`
}
