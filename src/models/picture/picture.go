package picture

import (
	"database/sql"
	"db"

	"models/comment"

	"github.com/jinzhu/gorm"
)

type Picture struct {
	gorm.Model

	ImageData  []byte
	UserID     sql.NullInt64 `gorm:"not null;index"`
	Type       string
	LikesCount int
	Comments   []comment.Comment `gorm:"ForeignKey:PictureID"`
}

func (p *Picture) Save() error {
	// Ну такое... при сохранении комментария каждый раз создается новый???
	return db.Get().Create(p).Error
}

func init() {
	db.Get().AutoMigrate(&Picture{})
}
