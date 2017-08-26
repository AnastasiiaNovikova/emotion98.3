package user

import (
	"db"
	"log"

	"models/comment"
	"models/picture"

	"github.com/jinzhu/gorm"
)

// User account
type User struct {
	gorm.Model

	Nickname     string `gorm:"index"`
	Email        string
	Pictures     []picture.Picture `gorm:"ForeignKey:UserID"`
	CommentsLeft []comment.Comment `gorm:"ForeignKey:AuthorID"`
}

func init() {
	db.Get().AutoMigrate(&User{})
}

// Add new user to the database
func (usr *User) Add() error {
	return db.Get().Create(usr).Error
}

// GetByID finds user (when you are sure that it exists)
func GetByID(id int) *User {
	var usr User
	db.Get().First(&usr, id)
	return &usr
}

// Get user by nickname
func Get(nickname string) *User {
	usr := User{}
	db.Get().Where("nickname = ?", nickname).First(&usr)
	if usr.Nickname == nickname {
		log.Printf("User found, id = %d\n", usr.ID)
		return &usr
	}
	log.Println("User not found")
	return nil
}

// GetPictures of a user
func GetPictures(nickname string) []picture.Picture {
	usr := Get(nickname)
	if usr != nil {
		var pictures []picture.Picture
		db.Get().Model(&usr).Related(&pictures)
		return pictures
	}
	return nil
}
