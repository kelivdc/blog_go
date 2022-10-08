package database

import (
	"blog/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func InitDatabase() {
	dsn := os.Getenv("DATABASE")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("Failed to open database")
	}
	fmt.Println("Success open database.")
	db.AutoMigrate(
		&models.Category{},
		&models.User{},
	)
	fmt.Println("Success database migration.")

	Database = DbInstance{Db: db}
}
