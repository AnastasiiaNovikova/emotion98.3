package picture

import (
	"database/sql"
	"db"
	"log"

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
		log.Printf("Picture found, id = %d\n", id)
		return &pict
	}
	log.Println("Picture not found")
	return nil
}

// GetComments of a picture
func GetComments(pictureID int) []comment.Comment {
	pict := Get(pictureID)
	if pict != nil {
		var comments []comment.Comment
		db.Get().Model(&pict).Related(&comments)
		return comments
	}
	return nil
}

// GetComments of a picture
func (p *Picture) GetComments() []comment.Comment {
	var comments []comment.Comment
	db.Get().Model(p).Related(&comments)
	return comments
}

func init() {
	db.Get().AutoMigrate(&Picture{})
}
