package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func NewGorm(viper *viper.Viper) *gorm.DB {
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	name := viper.GetString("database.name")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	idl := viper.GetInt("database.pooling.idle")
	maxConn := viper.GetInt("database.pooling.max")
	lifetime := viper.GetInt("database.pooling.lifetime")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", user, password, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	connection, err := db.DB()
	if err != nil {
		panic(err)
	}
	connection.SetMaxIdleConns(idl)
	connection.SetMaxOpenConns(maxConn)
	connection.SetConnMaxLifetime(time.Second * time.Duration(lifetime))
	return db
}
