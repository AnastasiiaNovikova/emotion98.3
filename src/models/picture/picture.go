package picture

import (
	"database/sql"
	"db"

	"models/comment"

	"github.com/jinzhu/gorm"
)

// Picture stores URL of user picture User uploaded
type Picture struct {
	gorm.Model
	URL string

	UserID   sql.NullInt64     `gorm:"not null;index"`
	Comments []comment.Comment `gorm:"ForeignKey:PictureID"`
}

// Save picture in a database
func (p *Picture) Save() error {
	return db.Get().Create(p).Error
}

// Get a picture by ID
func Get(id int) *Picture {
	pict := Picture{}
	db.Get().First(&pict, id)
	if int(pict.ID) == id {
		return &pict
	}
	return nil
}

// GetComments of a picture
func GetComments(pictureID int, maxCount int) {

}

func init() {
	db.Get().AutoMigrate(&Picture{})
}
