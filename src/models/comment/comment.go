package comment

import (
	"database/sql"
	"db"

	"github.com/jinzhu/gorm"
)

// Comment under a picture
type Comment struct {
	gorm.Model

	AuthorID  sql.NullInt64 `gorm:"not null;index"`
	PictureID sql.NullInt64 `gorm:"not null;index"`
	Text      string
}

// Leave a comment
func (c *Comment) Leave() error {
	return db.Get().Create(c).Error
}

// Leave a comment from scratch
func Leave(authorID int, pictureID int, text string) error {
	cmt := Comment{
		AuthorID:  db.Int64FK(int64(authorID)),
		PictureID: db.Int64FK(int64(pictureID)),
		Text:      text,
	}
	return db.Get().Create(&cmt).Error
}

func init() {
	db.Get().AutoMigrate(&Comment{})
}
