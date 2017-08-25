package comment

import (
	"database/sql"
	"db"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model

	AuthorID  sql.NullInt64 `gorm:"not null;index"`
	PictureID sql.NullInt64 `gorm:"not null;index"`
	Text      string
}

func (c *Comment) Save() error {
	// Ну такое... при сохранении комментария каждый раз создается новый???
	return db.Get().Create(c).Error
}

func init() {
	db.Get().AutoMigrate(&Comment{})
}
