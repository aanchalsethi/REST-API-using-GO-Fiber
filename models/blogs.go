package models

import "gorm.io/gorm"

type Blog struct {
	ID     int     `gorm:"primary key;autoIncrement" json:"id"`
	Author *string `json:"author"`
	Title  *string `json:"title"`
	Body   *string `json:"body"`
}

func MigrateBlogs(db *gorm.DB) error {
	err := db.AutoMigrate(&Blog{})
	return err
}
