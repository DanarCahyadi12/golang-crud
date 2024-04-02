package entity

import "time"

type User struct {
	Id        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Product   []Product `gorm:"foreignKey:user_id;references:id"`
}
