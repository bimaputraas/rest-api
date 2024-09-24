package pkgdatabase

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGorm(uri string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(uri), &gorm.Config{})
}
