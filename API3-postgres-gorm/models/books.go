package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

// The MigrateBooks function takes a *gorm.DB argument, which represents the database connection,
// and returns an error. The function uses the AutoMigrate method provided by gorm to automatically
// create or update the Books table in the database schema based on the Books struct definition.
// The err variable captures any error that may occur during the migration and returns it.
func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
