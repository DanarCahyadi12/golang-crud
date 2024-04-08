package entity

type Product struct {
	Id     string `gorm:"column:id;primaryKey"`
	Name   string `gorm:"column:name"`
	Price  int    `gorm:"column:price"`
	Stock  int    `gorm:"column:stock"`
	UserId string `gorm:"column:user_id;"`
	User   User   `gorm:"foreignKey:user_id;references:id"`
}
